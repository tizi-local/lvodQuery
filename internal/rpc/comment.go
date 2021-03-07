package rpc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	lvodQuery "github.com/tizi-local/commonapis/api/vodQuery"
	"github.com/tizi-local/lvodQuery/internal/cache"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/pkg/models"
	"strconv"
	"time"
)

// 一级评论查询
func (a *VodQueryService) CommentQueryFirst(ctx context.Context, req *lvodQuery.CommentQueryReq) (*lvodQuery.CommentQueryResp, error) {
	metaDataKey := cache.Key(cache.VodCommentKeyFirst, req.Vid)
	//scoreRank := cache.Key(cache.VodCommentScore, req.Vid)
	timelineRank := cache.Key(cache.VodCommentKeyTimeline, req.Vid)
	if cache.Exist(ctx, metaDataKey) == 0 {
		// redis 数据不存在 需要从mysql中构建 缓存
		a.Warnf("reconstruct cache for video: %s", req.Vid)
		a.rebuildVideoCommentCache(ctx, req.Vid)
	}
	// update cache expire time
	_, err := cache.Expire(ctx, req.GetVid(), 24*time.Hour)
	if err != nil {
		a.Errorf("Expire key:%v failed,err:%v", req.GetVid(), err)
		return nil, err
	}
	// 根据时间线拿comment
	// TODO 根据scoreRank
	keys := cache.ZREVRange(ctx, timelineRank, req.GetPage()*20, (req.GetPage()+1)*20-1)
	a.Debug("timeline comment keys:", keys)
	res, err := cache.HMGet(ctx, metaDataKey, keys...)
	if err != nil {
		a.Errorf("get comment data error: %s", err.Error())
		return nil, err
	}
	comments := make([]*lvodQuery.Comments, 0)
	for i := range res {
		comment := models.CommentFirst{}
		err := jsoniter.UnmarshalFromString(res[i].(string), &comment)
		if err != nil {
			a.Error("unmarshal error:", err.Error(), res[i].(string))
			continue
		}
		comments = append(comments, &lvodQuery.Comments{
			CommentId:   GenerateCommentId(comment.Id),
			Comment:     comment.Message,
			Uid:         comment.Uid,
			Timestamp:   comment.CreateTime.Unix(),
			SubComCount: int64(comment.SubComCount),
		})
	}

	return &lvodQuery.CommentQueryResp{
		Total:    int64(len(comments)),
		Page:     req.GetPage(),
		Comments: comments,
	}, nil
}

//TODO 抽象，加缓存
// 二级评论查询
func (a *VodQueryService) CommentQuerySecond(ctx context.Context, req *lvodQuery.CommentQueryReq) (*lvodQuery.CommentQueryResp, error) {
	commentId := req.CommentId
	if commentId == "" {
		a.Errorf("invalid comment id: %s", commentId)
		return nil, fmt.Errorf("invalid commentId: %s", commentId)
	}
	commentReplies := make([]models.CommentReply, 0)
	err := db.GetDb().Table("comment_reply").Where("comment_id = ?", ParseCommentId(commentId)).
		Limit(20, int(req.Page)*20).
		Find(&commentReplies)
	if err != nil {
		a.Errorf("query db err:%v", err)
		return nil, err
	}
	a.Debug("db comment reply:", commentReplies)
	comments := make([]*lvodQuery.Comments, 0)
	for _, comment := range commentReplies {
		comments = append(comments, &lvodQuery.Comments{
			CommentId: GenerateCommentReplaId(comment.Id),
			Comment:   comment.CommentContent,
			Uid:       comment.Uid,
			Timestamp: comment.CreateTime.Unix(),
			ReplyId:   GenerateCommentId(comment.CommentId),
		})
	}

	return &lvodQuery.CommentQueryResp{
		Total:    int64(len(comments)),
		Page:     req.GetPage(),
		Comments: comments,
	}, nil
}

// 一级评论创建
func (a *VodQueryService) CommentCreateFirst(ctx context.Context, req *lvodQuery.CommentCreateReq) (*lvodQuery.Error, error) {
	dbSession := db.GetDb().NewSession()
	if req.Vid == "" {
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("empty vid"),
			Details: nil,
		}, fmt.Errorf("empty vid")
	}
	dbSession.Begin()
	_, err := dbSession.Table("video_info").Where("vid = ?", req.Vid).Incr("comment_count").Update(&models.VideoInfo{
		Vid: req.Vid,
	})
	if err != nil {
		dbSession.Rollback()
		a.Errorf("incr %s count error: %v", req.Vid, err)
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  err.Error(),
			Details: nil,
		}, fmt.Errorf("incr error")
	}

	newComment := &models.CommentFirst{
		Vid:     req.Vid,
		Uid:     req.Uid,
		Message: req.Comment,
	}
	for _, v := range req.MentionUser {
		mentionUser := models.MentionUser{
			UserName: v.UserName,
			Uid:      v.Uid,
		}
		newComment.MentionUsers = append(newComment.MentionUsers, mentionUser)
	}
	a.Debugf("new comment:%+v", newComment)
	rows, err := dbSession.InsertOne(newComment)
	if err != nil || rows == 0 {
		a.Errorf("Comment write to db failed,err:%v\n", err)
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("Insert to db failed"),
			Details: nil,
		}, err
	}
	dbSession.Commit()
	// 插入新评论，更新缓存
	go a.rebuildVideoCommentCache(ctx, newComment.Vid)
	return &lvodQuery.Error{
		ErrCode: lvodQuery.ErrCode_Success,
		ErrMsg:  "",
	}, nil
}

// 二级评论创建
func (a *VodQueryService) CommentCreateSecond(ctx context.Context, req *lvodQuery.CommentCreateReq) (*lvodQuery.Error, error) {
	commentId := ParseCommentId(req.CommentId)
	if commentId < 0 {
		return nil, fmt.Errorf("invalid param commentId")
	}
	dbSession := db.GetDb().NewSession()
	dbSession.Begin()
	newCommentReply := models.CommentReply{
		CommentId:      commentId,
		Uid:            req.Uid,
		CommentContent: req.Comment,
		Vid:            req.Vid,
	}
	for _, v := range req.MentionUser {
		mentionUser := models.MentionUser{
			UserName: v.UserName,
			Uid:      v.Uid,
		}
		newCommentReply.MentionUsers = append(newCommentReply.MentionUsers, mentionUser)
	}
	rows, err := dbSession.InsertOne(&newCommentReply)
	if err != nil {
		dbSession.Rollback()
		return nil, err
	}
	dbSession.Table("video_info").Where("vid = ?", req.Vid).Incr("comment_count").Update(&models.CommentFirst{
		Id: commentId,
	})
	rows, err = dbSession.Table("comment_first").Where("id = ?", commentId).Incr("sub_com_count").Update(new(models.CommentReply))
	if err != nil || rows == 0 {
		a.Errorf("Comment write to db failed,err:%v\n", err)
		_ = dbSession.Rollback()
		return &lvodQuery.Error{
			ErrCode: lvodQuery.ErrCode_Failed,
			ErrMsg:  fmt.Sprint("Insert to db failed"),
		}, err
	}
	dbSession.Commit()
	return &lvodQuery.Error{
		ErrCode: lvodQuery.ErrCode_Success,
		ErrMsg:  "",
	}, nil
}

//TODO
func DeleteComment() {

}

func LikeComment() {

}

// 构建评论缓存
func (a *VodQueryService) rebuildVideoCommentCache(ctx context.Context, vid string) error {
	metaDataKey := cache.Key(cache.VodCommentKeyFirst, vid)
	scoreRank := cache.Key(cache.VodCommentScore, vid)
	timelineRank := cache.Key(cache.VodCommentKeyTimeline, vid)

	if err := a.deleteVodCommentCache(ctx, vid); err != nil {
		return err
	}

	commentsFirst := make([]models.CommentFirst, 0)
	session := db.GetDb().NewSession()
	session.Begin()
	page := 0

	//for commentsFirst {
	err := session.SQL("select * from comment_first where vid = ? order by id limit ?, ? for update", vid, page, (page+1)*5000).Find(&commentsFirst)
	if err != nil {
		session.Rollback()
		return err
	}
	for _, comment := range commentsFirst {
		a.Debug("comment:", comment.CreateTime.UnixNano(), comment)
		str, err := jsoniter.MarshalToString(comment)
		if err != nil {
			return err
		}
		commentId := GenerateCommentId(comment.Id)
		cache.HSet(ctx, metaDataKey, commentId, str)
		cache.ZAdd(ctx, scoreRank, &redis.Z{Score: comment.Score, Member: commentId})
		cache.ZAdd(ctx, timelineRank, &redis.Z{Score: float64(comment.CreateTime.UnixNano()), Member: commentId})
	}
	cache.Expire(ctx, metaDataKey, time.Minute*60)
	cache.Expire(ctx, scoreRank, time.Minute*60)
	cache.Expire(ctx, timelineRank, time.Minute*60)
	//}
	session.Commit()
	return nil
}

// 删除评论缓存
func (a *VodQueryService) deleteVodCommentCache(ctx context.Context, vid string) error {
	metaDataKey := cache.Key(cache.VodCommentKeyFirst, vid)
	scoreRank := cache.Key(cache.VodCommentScore, vid)
	timelineRank := cache.Key(cache.VodCommentKeyTimeline, vid)

	_, err := cache.Del(ctx,
		metaDataKey,
		scoreRank,
		timelineRank,
	)
	if err != nil {
		return err
	}
	return nil
}

// 一级评论ID转换
func GenerateCommentId(commentId int64) string {
	return strconv.FormatInt(1<<60+1<<61+commentId, 10)
}

func ParseCommentId(commentIdStr string) int64 {
	id, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		return 0
	}
	return id - 1<<60 - 1<<61
}

// 二级评论 GenerateCommentReply ID
func GenerateCommentReplaId(commentReplyId int64) string {
	return strconv.FormatInt(1<<61+1<<62+commentReplyId, 10)
}

func ParseCommentReplyId(comment *lvodQuery.Comments) int64 {
	id, err := strconv.ParseInt(comment.CommentId, 10, 64)
	if err != nil {
		return 0
	}
	return id - 1<<61 - 1<<62
}
