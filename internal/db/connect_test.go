package db

import (
	"github.com/stretchr/testify/assert"
	ljson "github.com/tizi-local/llib/encoding/json"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/lvodQuery/config"
	"github.com/tizi-local/lvodQuery/pkg/models"

	//"github.com/tizi-local/lauthority/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitDBEngine(&config.DBConfig{
		Username: "root",
		Password: "123456",
		DbName:   "lvod",
		Addr:     "47.113.90.8:3306",
	}, log.New())
	os.Exit(m.Run())
}

func TestFind(t *testing.T) {
	s := GetDb().NewSession()
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = s.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	s := GetDb().NewSession()
	defer s.Close()
	err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	comment := &models.CommentReply{
		Vid:            "deb9ccc4-fcb9-4dc2-901f-d1172297b73b",
		Uid:            "38643039333434622d626238342d3463",
		CommentContent: "您可赶紧发货吧",
	}
	exist, err := s.InsertOne(comment)
	if err != nil {
		t.Error("error", err)
		return
	}
	t.Log("exist", exist, *comment)
	err = s.Commit()
	if err != nil {
		t.Fatal(err)
	}

}

func TestM2MFind(t *testing.T) {
	s := GetDb().NewSession()
	_, err := s.Insert([]*models.Location{
		{
			Lid:       models.GenerateLid(22.4939700, 114.0723000),
			Latitude:  114.0723000,
			Longitude: 22.4939700,
			Category:  1,
			Address:   "shenzhen 110-221",
			Country:   "china",
			Province:  "guangdong",
			City:      "shenzhen",
			District:  "futian",
			Street:    "fuqiang",
			CityCode:  "0755",
			AdCode:    "4000",
			PoiName:   "poi example2",
			AoiName:   "futian",
		},
		{
			Lid:       models.GenerateLid(22.4939614, 114.0728131),
			Latitude:  114.0728131,
			Longitude: 22.4939614,
			Category:  1,
			Address:   "shenzhen 110-222",
			Country:   "china",
			Province:  "guangdong",
			City:      "shenzhen",
			District:  "futian",
			Street:    "fuqiang",
			CityCode:  "0755",
			AdCode:    "4000",
			PoiName:   "poi example",
			AoiName:   "futian",
		},
	})
	assert.NoError(t, err, "insert error")
	defer s.Close()
}

func TestJoinGetDb(t *testing.T) {
	s := GetDb().NewSession()
	locations, err := models.QueryVideoLocations(s, "lv867f7622fe5d4cf972d801dcbcc777f4")
	assert.NoError(t, err, "query video locations error")
	defer s.Close()
	t.Log(ljson.MustMarshalIndentToString(locations))
}
