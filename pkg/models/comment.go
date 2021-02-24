package models

import "time"

type CommentFirst struct {
	Id           int64         `xorm:"pk autoincr BIGINT" json:"commentId"`
	Vid          string        `xorm:"not null comment('全局vid') VARCHAR(50)" json:"-"`
	Uid          int64         `xorm:"not null comment('全局uid') BIGINT(20)" json:"uid"`
	MentionUsers []MentionUser `xorm:"TEXT"`
	Message      string        `xorm:"not null default '' VARCHAR(100)" json:"comment"`
	CreateTime   time.Time     `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime   time.Time     `xorm:"updated TIMESTAMP"`
	SubComCount  int32         `xorm:"not null default 0 INT"`
	Score        float64       `xorm:"not null default 0 DOUBLE"'`
}
type MentionUser struct {
	UserName string `json:"userName"`
	Uid      int64  `json:"uid"`
}
type CommentReply struct {
	Id             int64         `xorm:"pk autoincr BIGINT" json:"replyId"`
	CommentId      int64         `xorm:"not null default 0 BIGINT" json:"commentId"`
	Uid            int64         `xorm:"not null comment('全局uid') BIGINT(20)" json:"uid"`
	MentionUsers   []MentionUser `xorm:"TEXT"`
	CommentContent string        `xorm:"not null default '' VARCHAR(100)" json:"comment"`
	Score          float64       `xorm:"not null default 0 DOUBLE"'`
	CreateTime     time.Time     `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime     time.Time     `xorm:"updated TIMESTAMP"`
}
