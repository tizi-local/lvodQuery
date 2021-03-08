package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/tizi-local/commonapis/api/schedule"
	"time"
)

// 多对多
type VideoLocation struct {
	Id        int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Lid       string    `xorm:"index notnull"` //添加逻辑外键
	Vid       string    `xorm:"index notnull"` //添加逻辑外键
	CreatedAt time.Time `xorm:"created TIMESTAMP"`
	UpdatedAt time.Time `xorm:"updated TIMESTAMP"`
	DeletedAt time.Time `xorm:"deleted TIMESTAMP"`
}

type Location struct {
	Id int64 `xorm:"pk autoincr BIGINT(20)" json:"id"`
	// 经_纬
	Lid string `xorm:"not null default '' Varchar(45) unique(lid_index)" json:"lid"`
	// 纬度
	Latitude float64 `xorm:"not null default 0 DOUBLE"`
	// 经度
	Longitude float64 `xorm:"not null default 0 DOUBLE"`
	// 地点类型
	Category int `xorm:"not null default 0 TINYINT(2)"`
	// 地址描述
	Address string `xorm:"not null default '' Varchar(255)"`
	// 国家
	Country string `xorm:"not null default '' Varchar(255)"`
	// 省份
	Province string `xorm:"not null default '' Varchar(255)"`
	// 城市
	City string `xorm:"not null default '' Varchar(255)"`
	// 城区
	District string `xorm:"not null default '' Varchar(255)"`
	// 街道
	Street string `xorm:"not null default '' Varchar(255)"`
	// 城市编码
	CityCode string `xorm:"not null default '' Varchar(255)"`
	// 区域编码，身份证前4位
	AdCode string `xorm:"not null default '' Varchar(255)"`
	// POI
	PoiName string `xorm:"not null default '' Varchar(255)"`
	// AOI
	AoiName string `xorm:"not null default '' Varchar(255)"`
}

func GenerateLid(longitude, latitude float64) string {
	return fmt.Sprintf("%.6f_%.6f", longitude, latitude)
}

func QueryVideoLocations(session *xorm.Session, vid string) ([]*Location, error) {
	locations := make([]*Location, 0)
	err := session.Join("LEFT", "video_location", "location.lid = video_location.lid").
		Where("video_location.vid = ?", vid).Find(&locations)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func ConvertLocationModel2Pb(locations []*Location) []*schedule.Location {
	pbLocations := make([]*schedule.Location, len(locations))
	for i := range locations {
		pbLocations[i] = &schedule.Location{
			Latitude:  locations[i].Latitude,
			Longitude: locations[i].Longitude,
			Category:  int32(locations[i].Category),
			Address:   locations[i].Address,
			Country:   locations[i].Country,
			Province:  locations[i].Province,
			City:      locations[i].City,
			District:  locations[i].District,
			Street:    locations[i].Street,
			CityCode:  locations[i].CityCode,
			AdCode:    locations[i].AdCode,
			PoiName:   locations[i].PoiName,
			AoiName:   locations[i].AoiName,
		}
	}

	return pbLocations
}

func ConvertPb2LocationModel(locations []*schedule.Location) []*Location {
	pbLocations := make([]*Location, len(locations))
	for i := range locations {
		pbLocations[i] = &Location{
			Latitude:  locations[i].Latitude,
			Longitude: locations[i].Longitude,
			Category:  int(locations[i].Category),
			Address:   locations[i].Address,
			Country:   locations[i].Country,
			Province:  locations[i].Province,
			City:      locations[i].City,
			District:  locations[i].District,
			Street:    locations[i].Street,
			CityCode:  locations[i].CityCode,
			AdCode:    locations[i].AdCode,
			PoiName:   locations[i].PoiName,
			AoiName:   locations[i].AoiName,
		}
		pbLocations[i].Lid = GenerateLid(pbLocations[i].Longitude, pbLocations[i].Latitude)
	}

	return pbLocations
}
