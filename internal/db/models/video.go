package models

import "time"

type VideoInfo struct {
	Id            int64  `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Vid           string `xorm:"not null comment('全局vid') VARCHAR(50)" json:"vid"`
	Url           string `xorm:"not null TEXT" json:"videoUrl"`
	AssetId       string `xorm:"not null default '' VARCHAR(50)" json:"-"`
	VideoTitle    string `xorm:"not null default '' VARCHAR(50)" json:"videoTitle"`
	VideoDesc     string `xorm:"not null default '' VARCHAR(255)" json:"videoDesc"`
	Success       bool   `xorm:"not null default 0 TINYINT(1)" json:"-"`
	CommentCount  int64  `xorm:"not null default 0 INT"`
	FavoriteCount int64  `xorm:"not null default 0 INT"`
	LikeCount     int64  `xorm:"not null default 0 INT"`
	CreateTime    time.Time `xorm:"created TIMESTAMP" json:"timestamp"`
	UpdateTime    time.Time `xorm:"updated TIMESTAMP"`
	Poi           []Poi     `json:"poi"`
	AuthorUid     string `json:"author_uid"`
}

type Poi struct {
	Vid       string  `xorm:"not null comment('全局vid') VARCHAR(50)" json:"-"`
	Name      string  `xorm:"not null default '' VARCHAR(50)" json:"name"`
	Loid      string  `xorm:"not null default '' VARCHAR(50)" json:"poiId"`
	Longitude float64 `xorm:"not null default 0 DOUBLE" json:"longitude"`
	Latitude  float64 `xorm:"not null default 0 DOUBLE" json:"latitude"`
	City      string  `xorm:"not null default '' VARCHAR(50)"`
	Province  string  `xorm:"not null default '' VARCHAR(50)"`
	Country   string  `xorm:"not null default '' VARCHAR(50)"`
}
