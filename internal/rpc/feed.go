package rpc

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/lvodQuery/internal/cache"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/internal/db/models"
	lvodQuery "github.com/tizi-local/lvodQuery/proto/vodQuery"
	"time"
)

type VodQueryService struct {
	lvodQuery.UnimplementedVodQueryServiceServer
	*log.Logger
}

func NewVodQueryService(logger *log.Logger) *VodQueryService {
	return &VodQueryService{
		Logger: logger,
	}
}

func (a *VodQueryService) FeedQuery(ctx context.Context, page *lvodQuery.FeedQueryReq) (*lvodQuery.FeedQueryResp, error) {
	session := page.GetVSession()
	sessionByte := []byte(session)
	fmt.Println("\n", session)
	count, err := cache.SNum(ctx, string(sessionByte))
	if err != nil {
		a.Errorf("Get Token failed", err)
		return nil, fmt.Errorf("get session failed,err:%v\n", err)
	} else if count == 0 {
		session, err := NewSession(64)
		if err != nil {
		}
		sessionByte = []byte(session)
		videoInfos := make([]models.VideoInfo, 0)
		//		rand := rand2.Int()
		err = db.GetDb().Table("video_info").Where("video_info.id > ? AND video_info.id <? AND video_info.success = 1", 0, 1+99).
			Find(&videoInfos)
		for _, v := range videoInfos {
			err = db.GetDb().Table("poi").Where("poi.vid = ?", v.Vid).
				Find(&(v.Poi))
			if err != nil {
				a.Errorf("get data from db failed,err:%v\n", err)
				return nil, err
			}
			l, err := jsoniter.Marshal(v)
			if err != nil {
				a.Errorf("JSON marshal failed,err:%v\n", err)
			}
			_, err = cache.SAdd(ctx, session, l)
			if err != nil {
				a.Errorf("Insert redis failed,err:%v\n", err)
			}
		}

		_, err = cache.SExpire(ctx, session, time.Hour*24)
		if count, _ := cache.SNum(ctx, session); count == 0 {
			a.Errorf("this feed is empty")
			return nil, errors.New("this feed is empty")
		}
	}
	responseVideos := make([]*lvodQuery.Videos, 0)
	for i := 0; i < 5; i++ {
		info, err := cache.SPop(ctx, string(sessionByte))
		if err != nil {
			continue
		}
		video := &lvodQuery.Videos{}
		jsoniter.Unmarshal([]byte(info), video)
		responseVideos = append(responseVideos, video)
	}
	return &lvodQuery.FeedQueryResp{
		Session: string(sessionByte),
		Total:   int64(len(responseVideos)),
		Page:    page.Page,
		Videos:  responseVideos,
	}, nil
}
