package rpc

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	lauthRpc "github.com/tizi-local/commonapis/api/authority"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/llib/log"
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
	session := page.GetVSession()
	a.Debugf("vSession: %s", session)
	if session != "" {
		// session input
		// 从redis里拿，拿不到报错
		if cache.Exist(ctx, session) == 1 {

			sessionCacheKey := cache.Key(cache.VodKeyFeedSession, session)
			start := page.Page * feed.FeedPageSize
			stop := (page.Page + 1) * feed.FeedPageSize
			cacheLen, err := cache.LLen(ctx, sessionCacheKey)
			if err != nil {
				a.Error("invalid request for cache", err.Error())
				return nil, fmt.Errorf("query Feed failed:%s", err.Error())
			}
			if stop > cacheLen {
				a.Error("exceed feed scan ", start, stop, cacheLen)
				return nil, fmt.Errorf("execeed feed scan")
			}
			keys, err := cache.LRange(ctx, sessionCacheKey, start, stop)
			if err != nil {
				a.Error("request cache failed ", err.Error())
				return nil, err
			}
			vDatas, err := cache.MGet(ctx, keys...)
			if err != nil {
				a.Error("request cache failed ", err.Error())
				return nil, err
			}
			responseVideos := make([]*lvodQuery.Videos, len(vDatas))
			for _, vData := range vDatas {
				video := lvodQuery.Videos{}
				err := jsoniter.Unmarshal(vData.([]byte), &video)
				if err != nil {
					a.Error("unmarshal error", err.Error())
					continue
				}
				responseVideos = append(responseVideos, &video)
			}
			return &lvodQuery.FeedQueryResp{
				Session: session,
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
		videoInfos := make([]models.VideoInfo, 0)
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
			Limit(100).
			Find(&videoInfos)
		if err != nil {
			a.Error("query video_info data failed", err.Error())
			return nil, fmt.Errorf("query feed failed")
		}
		if len(videoInfos) != 0 {
			// generate new session and cache
			session := a.feedGenerator.GenerateSession()
			responseVideos := make([]*lvodQuery.Videos, 0)
			cc, err := grpc.Dial(a.config.Auth, grpc.WithInsecure())
			if err != nil {
				a.Error("dial sms rpc error:", err.Error())
				return nil, err
			}
			authClient := lauthRpc.NewAuthServiceClient(cc)
			for _, v := range videoInfos {
				// get poi from db
				// TODO use join
				err = db.GetDb().Table("poi").Where("poi.vid = ?", v.Vid).
					Find(&(v.Poi))
				if err != nil {
					a.Errorf("get data from db failed,err:%v\n", err)
					return nil, err
				}
				// get user info from lauth
				user, err := authClient.GetUserInfo(ctx, &lauthRpc.UserRequest{
					Uid: v.AuthorUid,
				})
				if err != nil {
					return nil, err
				}
				// set videoInfo cache
				err = a.cacheVideoInfo(ctx, &v)
				if err != nil {
					continue
				}
				_, err = cache.RPush(ctx, session, v.Vid)
				if err != nil {
					a.Errorf("Insert redis failed,err:%v\n", err)
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
				})
			}
			_, err = cache.Expire(ctx, session, time.Hour*24)
			if err != nil {
				a.Errorf("expire redis key failed,err:%v\n", err)
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
	l, err := jsoniter.Marshal(v)
	if err != nil {
		a.Errorf("JSON marshal failed,err:%v\n", err)
		return err
	}
	_, err = cache.SetExpire(ctx, cache.Key(cache.VodKeyVideo, v.Vid), string(l), 24*time.Hour)
	return err
}
