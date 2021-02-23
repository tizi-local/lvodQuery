package models

import "time"

type FavoriteList struct {
	CollectId    int64     `xorm:pk wq`
	CollectUid   int64     `xorm:"not null comment('全局uid') BIGINT(20)"`
	CollectVid   string    `xorm:"not null comment('全局vid') VARCHAR(50)"`
	CollectState bool      `xorm:"not null default 0 Bool"` //0为收藏，1为取消
	CollectTime  time.Time `xorm:"updated TIMESTAMP"`
}

type LikeList struct {
	CollectId    int64     `xorm:pk wq`
	CollectUid   int64     `xorm:"not null comment('全局uid') BIGINT(20)"`
	CollectVid   string    `xorm:"not null comment('全局vid') VARCHAR(50)"`
	CollectState bool      `xorm:"not null default 0 Bool"` //0为👍，1为取消
	CollectTime  time.Time `xorm:"updated TIMESTAMP"`
}
