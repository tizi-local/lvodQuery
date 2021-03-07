package models

import "time"

type CommentFirst struct {
	Id           int64         `xorm:"pk autoincr BIGINT" json:"commentId"`
	Vid          string        `xorm:"not null comment('全局vid') VARCHAR(50)" json:"vid"`
	Uid          string        `xorm:"not null comment('全局uid') VARCHAR(32)" json:"uid"`
	MentionUsers []MentionUser `xorm:"TEXT" json:"mention_users"`
	Message      string        `xorm:"not null default '' VARCHAR(100)" json:"comment" `
	CreateTime   time.Time     `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime   time.Time     `xorm:"updated TIMESTAMP"`
	SubComCount  int32         `xorm:"not null default 0 INT" json:"sub_com_count"`
	Score        float64       `xorm:"not null default 0 DOUBLE" json:"score"`
}

type MentionUser struct {
	UserName string `json:"userName"`
	Uid      string `json:"uid"`
}

type CommentReply struct {
	Id             int64         `xorm:"pk autoincr BIGINT" json:"replyId"`
	CommentId      int64         `xorm:"not null default 0 BIGINT comment('一级评论Id')" json:"commentId"`
	Vid            string        `xorm:"not null comment('全局vid') VARCHAR(50)" json:"vid"`
	Uid            string        `xorm:"not null comment('全局uid') BIGINT(20)" json:"uid"`
	MentionUsers   []MentionUser `xorm:"TEXT"`
	CommentContent string        `xorm:"not null default '' VARCHAR(100)" json:"comment"`
	Score          float64       `xorm:"not null default 0 DOUBLE"`
	CreateTime     time.Time     `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime     time.Time     `xorm:"updated TIMESTAMP"`
}
