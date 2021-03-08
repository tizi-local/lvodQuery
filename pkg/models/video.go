package models

import "time"

type VideoInfo struct {
	Id            int64      `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Vid           string     `xorm:"not null comment('全局vid') VARCHAR(50)" json:"vid"`
	Url           string     `xorm:"not null TEXT" json:"videoUrl"`
	AssetId       string     `xorm:"not null default '' VARCHAR(50)" json:"-"`
	Bucket        string     `xorm:"not null default '' VARCHAR(50)"`
	Object        string     `xorm:"not null default '' VARCHAR(255)"`
	VideoTitle    string     `xorm:"not null default '' VARCHAR(50)" json:"videoTitle"`
	VideoDesc     string     `xorm:"not null default '' VARCHAR(255)" json:"videoDesc"`
	Success       bool       `xorm:"not null default 0 TINYINT(1)" json:"-"`
	CommentCount  int64      `xorm:"not null default 0 INT"`
	FavoriteCount int64      `xorm:"not null default 0 INT"`
	LikeCount     int64      `xorm:"not null default 0 INT"`
	CreateTime    time.Time  `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime    time.Time  `xorm:"updated TIMESTAMP"`
	AuthorUid     string     `json:"author_uid"`
	GoodsUrl      string     `xorm:"not null TEXT" json:"goods_url"`
	CovertUrl     string     `xorm:"not null TEXT" json:"covert_url"`
	Deleted       bool       `xorm:"not null default 0 TINYINT(1)" json:"-"` // 软删
	Locations     []Location `xorm:"-" json:"locations"`
}
