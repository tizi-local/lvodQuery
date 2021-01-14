package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/lvodQuery/config"
	"github.com/tizi-local/lvodQuery/internal/db/models"
	"xorm.io/core"
)

var (
	engine *xorm.Engine
)

func InitDBEngine(c *config.DBConfig, l *log.Logger) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", c.Username, c.Password, c.Addr, c.DbName)
	fmt.Printf("connect to db: %s\n", dsn)
	e, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		l.Fatalln("init data failed")
		return
	}
	engine = e
	err = e.Sync2(new(models.VideoInfo), new(models.Poi),new(models.CommentFirst),new(models.CommentReply))
	if err != nil {
		fmt.Println(err)
	}
	// TODO enable in debug; disable in release
	engine.SetTableMapper(core.SnakeMapper{})
	dbLogger := xorm.NewSimpleLogger(l.Writer())
	dbLogger.SetLevel(core.LOG_DEBUG)
	engine.SetLogger(dbLogger)
}

func GetDb() *xorm.Engine {
	return engine
}
