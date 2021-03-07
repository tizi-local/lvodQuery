package rpc

import (
	"context"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/pkg/models"
)

func (a *VodQueryService) Like(ctx context.Context, req *lvodQuery.LikeReq) (*lvodQuery.Error, error) {
	dbSession := db.GetDb().NewSession()
	dbSession.Begin()
	_, err := dbSession.Table("video_info").Where("vid = ?", req.Vid).Incr("like_count").Update(new(models.VideoInfo))
	if err != nil{
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, err
	}
	newLikeRecord := models.LikeList{
		Uid:   req.Uid,
		Vid:   req.Vid,
		State: false,
	}
	_, err = dbSession.Table("like_list").InsertOne(newLikeRecord)
	if err != nil {
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, err
	}
	dbSession.Commit()
	return &lvodQuery.Error{
		ErrCode: 0,
		ErrMsg:  "",
	}, nil
}

func (a *VodQueryService) LikeQuery(ctx context.Context, req *lvodQuery.ListQuery) (*lvodQuery.FeedQueryResp, error) {
	vids := make([]string, 0)
	err := db.GetDb().Table("like_list").Select("vid").Where("uid=?", req.GetUid()).
		Limit(20, int(20 * req.Page)).
		Find(&vids)
	if err != nil {
		a.Errorf("query like list error: %v", err)
		return nil, err
	}
	videos := make([]*models.VideoInfo, 0)
	err = db.GetDb().Table("video_info").In("vid", vids).Find(&videos)
	if err != nil{
		a.Errorf("query videos error: %v", err)
		return nil, err
	}
	a.Debug("db get videos:", videos)

	respVideos, err := a.ConvertVideoModel2Pb(ctx, videos)
	return &lvodQuery.FeedQueryResp{
		Total:  int64(len(respVideos)),
		Page:   req.Page,
		Videos: respVideos,
	}, nil
}
