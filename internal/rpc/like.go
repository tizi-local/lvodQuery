package rpc

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/internal/db/models"
)

func (a *VodQueryService)Like(ctx context.Context,req *lvodQuery.LikeReq)(*lvodQuery.Error,error){
	dbSession := db.GetDb().NewSession()
	dbSession.Table("video_info").Where("vid = ?",req.Vid).Incr("like_count")
	writeToDB := models.LikeList{
		CollectUid:   req.Uid,
		CollectVid:   req.Vid,
		CollectState: false,
	}
	_, err := dbSession.Table("like_list").InsertOne(writeToDB)
	if err != nil {
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		},err
	}
	return &lvodQuery.Error{
		ErrCode: 0,
		ErrMsg:  "",
	}, nil
}
func (a *VodQueryService)LikeQuery(ctx context.Context,req *lvodQuery.ListQuery)(*lvodQuery.FeedQueryResp,error){
	vids := make([]string,0)
	err := db.GetDb().Table("like_list").Where("uid=?", int(req.GetUid())).Find(&vids)
	if err != nil {
		return nil, err
	}
	videoInfos := make([]*lvodQuery.Videos, 0)
	for _,v := range vids{
		videoInfo := models.VideoInfo{}
		db.GetDb().Table("video_info").Where("vid=?",v).Get(&videoInfo)
		info, err := jsoniter.Marshal(videoInfo)
		if err != nil {
			a.Errorf("Json marshal failed,err:%v+%v",videoInfo,err)
			continue
		}
		video := &lvodQuery.Videos{}
		err = jsoniter.Unmarshal(info, video)
		if err != nil {
			a.Errorf("Json Unmarshal failed,err:%v+%v",videoInfo,err)
			continue
		}
		videoInfos = append(videoInfos, video)
	}
	return &lvodQuery.FeedQueryResp{
		Total:   int64(len(videoInfos)),
		Page:    req.Page,
		Videos:  videoInfos,
	},nil
}

