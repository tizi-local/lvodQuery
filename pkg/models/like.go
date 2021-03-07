package models

import "time"

type LikeList struct {
	Id    int64     `xorm:pk wq`
	Uid   string    `xorm:"not null comment('全局uid') VARCHAR(32)"`
	Vid   string    `xorm:"not null comment('全局vid') VARCHAR(50)"`
	State bool      `xorm:"not null default 0 Bool"` //0为👍，1为取消
	Time  time.Time `xorm:"updated TIMESTAMP"`
}
