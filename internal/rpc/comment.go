package rpc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/lvodQuery/internal/cache"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/internal/db/models"
	"strconv"
	"time"
)

func (a *VodQueryService) CommentQueryFirst(ctx context.Context, page *lvodQuery.CommentQueryReq) (*lvodQuery.CommentQueryResp, error) {
	if cache.Exist(ctx,page.GetVid()) == 0|| page.GetPage() > 5{
		var start int = 0
		var limit = 100
		if page.GetPage() > 5{
			start = int(page.Page-1)*20
			limit = 20
		}
		CommentIndex := make([]models.CommentFirst, 0)
		db.GetDb().Table("comment_first").Where("vid = ?", page.GetVid()).Limit(limit, start).Find(&CommentIndex)
		for _, v := range CommentIndex {
			info, err := jsoniter.Marshal(v)
			if err != nil {
				a.Errorf("Marshal failed,err:%v\n", err)
				continue
			}
			item := &redis.Z{
				Score:  v.Score,
				Member: info,
			}
			cache.ZAdd(ctx, page.Vid, item)
		}
	}
	_, err := cache.Expire(ctx, page.GetVid(), 24*time.Hour)
	if err != nil {
		a.Errorf("Expire key:%v failed,err:%v\n",page.GetVid(),err)
		return nil, err
	}
	resps := make([]*lvodQuery.Comments,0)
	items := make([][]byte,0)
	err = cache.ZRange(ctx,page.GetVid(),items,int64((page.GetPage()-1)*20),int64(page.GetPage()*20),)
	if err != nil{
		a.Errorf("Cache have some trouble in ZRange,err:%v\n",err)
	}
	for _,v := range items{
		resp := &lvodQuery.Comments{}
		err = jsoniter.Unmarshal(v,resp)
		if err != nil {
			a.Errorf("Comment Item unmarshal failed,err:%v\n",err)
			continue
		}
		resps = append(resps, resp)
	}
	if len(resps) == 0 {
		return nil, err
	}
	return &lvodQuery.CommentQueryResp{
		Total:    int64(len(resps)),
		Page:     page.GetPage(),
		Comments: resps,
	},nil
}
func (a *VodQueryService) CommentQuerySecond(ctx context.Context, page *lvodQuery.CommentQueryReq) (*lvodQuery.CommentQueryResp, error) {
	keyId := string(page.GetCommentId())
	if cache.Exist(ctx,keyId) == 0{
		CommentIndex := make([]models.CommentReply, 0)
		db.GetDb().Table("comment_reply").Where("comment_id = ?", page.GetCommentId()).Find(&CommentIndex)
		for _, v := range CommentIndex {
			info, err := jsoniter.Marshal(v)
			if err != nil {
				a.Errorf("Marshal failed,err:%v\n", err)
				continue
			}
			item := &redis.Z{
				Score:  v.Score,
				Member: info,
			}
			cache.ZAdd(ctx, keyId, item)
		}
		count,err := cache.ZNum(ctx,keyId)
		if count == 0 || err!=nil{
			return nil, err
		}
	}
	resps := make([]*lvodQuery.Comments,0)
	items := make([][]byte,0)
	err := cache.ZRange(ctx,keyId,items,int64((page.GetPage()-1)*20),int64(page.GetPage()*20),)
	if err != nil{
		a.Errorf("Cache have some trouble in ZRange,err:%v\n",err)
	}
	for _,v := range items{
		resp := &lvodQuery.Comments{}
		err = jsoniter.Unmarshal(v,resp)
		if err != nil {
			a.Errorf("Comment Item unmarshal failed,err:%v\n",err)
			continue
		}
		resps = append(resps, resp)
	}
	if len(resps) == 0 {
		return nil, err
	}
	return &lvodQuery.CommentQueryResp{
		Total:    int64(len(resps)),
		Page:     page.GetPage(),
		Comments: resps,
	},nil
}
func (a *VodQueryService) CommentCreateFirst(ctx context.Context, req *lvodQuery.CommentCreateReq)(*lvodQuery.Error,error) {
	uid, err := strconv.Atoi(req.GetUid())
	if err != nil {
		a.Errorf("convert uid %s to int failed,err:%v\n",req.GetUid(),err)
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("convert uid to int failed"),
			Details: nil,
		},err
	}
	dbSession := db.GetDb().NewSession()
	dbSession.Table("video_info").Where("vid = ?",req.Vid).Incr("comment_count")
	writeToDB := models.CommentFirst{
		Vid:     req.Vid,
		Uid:     int64(uid),
		Message: req.Comment,
		Score:   float64(req.CommentId),
	}
	for _,v := range req.MentionUser{
		uid, _ := strconv.Atoi(v.Uid)
		mentionUser := models.MentionUser{
			UserName: v.UserName,
			Uid:      int64(uid),
		}
		writeToDB.MentionUsers = append(writeToDB.MentionUsers, mentionUser)
	}
	rows, err := dbSession.InsertOne(writeToDB)
	if err != nil || rows == 0{
		a.Errorf("Comment write to db failed,err:%v\n",err)
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("Insert to db failed"),
			Details: nil,
		},err
	}
	return &lvodQuery.Error{
		ErrCode: lvodQuery.ErrCode_Success,
		ErrMsg:  "",
	}, nil
}
func (a *VodQueryService) CommentCreateSecond(ctx context.Context, req *lvodQuery.CommentCreateReq)(*lvodQuery.Error,error){
	uid, err := strconv.Atoi(req.GetUid())
	if err != nil {
		a.Errorf("convert uid %s to int failed,err:%v\n",req.GetUid(),err)
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("convert uid to int failed"),
			Details: nil,
		},err
	}
	dbSession := db.GetDb().NewSession()
	dbSession.Table("video_info").Where("vid = ?",req.Vid).Incr("comment_count")
	writeToDB := models.CommentReply{
		CommentId:      req.CommentId,
		Uid:            int64(uid),
		CommentContent: req.Comment,
		Score:          float64(req.CommentId),
	}
	for _,v := range req.MentionUser{
		uid, _ := strconv.Atoi(v.Uid)
		mentionUser := models.MentionUser{
			UserName: v.UserName,
			Uid:      int64(uid),
		}
		writeToDB.MentionUsers = append(writeToDB.MentionUsers, mentionUser)
	}
	rows, err := dbSession.InsertOne(writeToDB)
	sql := "update `comment_first` set sub_com_count=sub_com_count+1 where id = ?"
	_, err = dbSession.Exec(sql, int(writeToDB.CommentId))
	if err != nil || rows == 0{
		a.Errorf("Comment write to db failed,err:%v\n",err)
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("Insert to db failed"),
		},err
	}
	return &lvodQuery.Error{
		ErrCode: lvodQuery.ErrCode_Success,
		ErrMsg:  "",
	}, nil
}
