package db

import (
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
