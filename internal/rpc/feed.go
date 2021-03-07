package rpc

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	lauthRpc "github.com/tizi-local/commonapis/api/authority"
	"github.com/tizi-local/commonapis/api/schedule"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/llib/utils/convert"
	"github.com/tizi-local/lvodQuery/config"
	"github.com/tizi-local/lvodQuery/internal/cache"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/internal/feed"
	"github.com/tizi-local/lvodQuery/pkg/models"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

type VodQueryService struct {
	lvodQuery.UnimplementedVodQueryServiceServer
	*log.Logger
	feedGenerator *feed.FeedGenerator
	config        *config.RpcConfig
}

func NewVodQueryService(logger *log.Logger, c *config.RpcConfig) *VodQueryService {
	return &VodQueryService{
		Logger:        logger,
		feedGenerator: feed.NewFeedGenerator(),
		config:        c,
	}
}

func (a *VodQueryService) FeedQuery(ctx context.Context, page *lvodQuery.FeedQueryReq) (*lvodQuery.FeedQueryResp, error) {
	a.Debugf("vSession: %s", page.GetVSession())
	if page.GetVSession() != "" {
		// session input
		// 从redis里拿，拿不到报错
		feedCacheKey := cache.Key(cache.VodKeyFeedSession, page.GetVSession())
		if cache.Exist(ctx, feedCacheKey) == 1 {
			start := page.Page * feed.FeedPageSize
			stop := (page.Page + 1) * feed.FeedPageSize - 1
			cacheLen, err := cache.LLen(ctx, feedCacheKey)
			if err != nil {
				a.Error("invalid request for cache", err.Error())
				return nil, fmt.Errorf("query Feed failed:%s", err.Error())
			}
			if start >= cacheLen {
				a.Error("exceed feed scan ", start, stop, cacheLen)
				return nil, feed.ErrorLastOfSession
			}
			keys, err := cache.LRange(ctx, feedCacheKey, start, stop)
			if err != nil {
				a.Error("request cache failed ", err.Error())
				return nil, err
			}
			a.Debug("feed keys ", keys)
			vDatas, err := a.mgetCachedVideoInfo(ctx, keys)
			if err != nil {
				a.Error("request cache failed ", err.Error())
				return nil, err
			}
			responseVideos, err := a.ConvertVideoModel2Pb(ctx, vDatas)
			if err != nil{
				a.Error("convert model to pb video error:", err)
				return nil, err
			}
			return &lvodQuery.FeedQueryResp{
				Session: page.GetVSession(),
				Total:   int64(len(responseVideos)),
				Page:    page.Page,
				Videos:  responseVideos,
			}, nil
		} else {
			a.Error("not existed feed session")
			return nil, fmt.Errorf("not existed feed session, need legal session")
		}
	} else {
		// no session input, create a new one`
		videoInfos := make([]*models.VideoInfo, 0)
		var videoCount int64
		// get videos info from db
		videoCount, err := db.GetDb().
			Where("video_info.success = ?", 1).
			Count(new(models.VideoInfo))
		if err != nil {
			a.Error("query video_info count failed", err.Error())
			return nil, fmt.Errorf("query feed failed")
		}
		// random scan
		offset := rand.Int63n(videoCount)
		err = db.GetDb().Table("video_info").
			Where("video_info.id >= ? AND video_info.success = ?", offset, 1).
			Limit(feed.FeedCountLimit).
			Find(&videoInfos)
		if err != nil {
			a.Error("query video_info data failed", err.Error())
			return nil, fmt.Errorf("query feed failed")
		}
		if len(videoInfos) != 0 {
			// generate new feed and cache
			newFeed := a.feedGenerator.NewFeed(videoInfos)
			session := newFeed.Session
			responseVideos := make([]*lvodQuery.Videos, 0)
			feedCacheKey := cache.Key(cache.VodKeyFeedSession, session)
			for _, v := range videoInfos {
				// get poi from db
				//err = db.GetDb().Table("poi").Where("poi.vid = ?", v.Vid).
				//	Find(&(v.Poi))
				//if err != nil {
				//	a.Errorf("get data from db failed,err:%v\n", err)
				//	return nil, err
				//}
				// get user info from lauth
				// set videoInfo cache
				err = a.cacheVideoInfo(ctx, v)
				if err != nil {
					continue
				}
				_, err = cache.RPush(ctx, feedCacheKey, cache.Key(cache.VodKeyVideo, v.Vid))
				if err != nil {
					a.Errorf("Insert redis failed,err:%v\n", err)
				}
			}
			_, err = cache.Expire(ctx, feedCacheKey, time.Hour*24)
			if err != nil {
				a.Errorf("expire redis key failed,err:%v\n", err)
				return nil, err
			}

			// return the first FeedPageSize videos in this session
			if len(videoInfos) > feed.FeedPageSize {
				videoInfos = videoInfos[:feed.FeedPageSize]
			}
			responseVideos, err = a.ConvertVideoModel2Pb(ctx, videoInfos)
			if err != nil{
				a.Error("convert model to pb video error:", err)
				return nil, err
			}
			return &lvodQuery.FeedQueryResp{
				Session: session,
				Page:    0,
				Total:   int64(len(videoInfos)),
				Videos:  responseVideos,
			}, nil
		} else {
			return nil, fmt.Errorf("failed get video feeds")
		}
	}
}

func (a *VodQueryService) cacheVideoInfo(ctx context.Context, v *models.VideoInfo) error {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		a.Errorf("JSON marshal failed,err:%v\n", err)
		return err
	}
	_, err = cache.SetExpire(ctx, cache.Key(cache.VodKeyVideo, v.Vid), string(bytes), 24*time.Hour)
	return err
}

func (a *VodQueryService) mgetCachedVideoInfo(ctx context.Context, vids []string) ([]*models.VideoInfo, error){
	res, err := cache.MGet(ctx,  vids...)
	if err != nil{
		return []*models.VideoInfo{}, err
	}
	videos := make([]*models.VideoInfo, 0)
	for i := range res {
		bs := convert.String2Bytes(res[i].(string))
		v := &models.VideoInfo{}
		err := jsoniter.Unmarshal(bs, v)
		if err != nil{
			continue
		}
		videos = append(videos, v)
	}

	return videos, nil
}

func (a *VodQueryService) ConvertVideoModel2Pb(ctx context.Context, videos []*models.VideoInfo)([]*lvodQuery.Videos, error){
	cc, err := grpc.Dial(a.config.Auth, grpc.WithInsecure())
	if err != nil {
		a.Error("dial sms rpc error:", err.Error())
		return nil, err
	}
	authClient := lauthRpc.NewAuthServiceClient(cc)
	responseVideos := make([]*lvodQuery.Videos, 0)
	for _, v := range videos {
		user, err := authClient.GetUserInfo(ctx, &lauthRpc.UserRequest{
			Uid: v.AuthorUid,
		})
		if err != nil {
			return nil, err
		}
		//TODO query locations
		responseVideos = append(responseVideos, &lvodQuery.Videos{
			Vid:           v.Vid,
			VideoUrl:      v.Url,
			VideoDesc:     v.VideoDesc,
			VideoTitle:    v.VideoTitle,
			LikeCount:     v.LikeCount,
			CommentCount:  v.CommentCount,
			ForwardCount:  0, // TODO forward count
			FavoriteCount: v.FavoriteCount,
			Author: &lvodQuery.Author{
				Uid:  v.AuthorUid,
				Name: user.UserName,
			},
			Locations: []*schedule.Location{},
			CoverUrl:      v.CovertUrl,
			GoodsUrl:      v.GoodsUrl,
		})
	}

	return responseVideos, nil
}