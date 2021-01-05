package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tizi-local/lvodQuery/internal/cache"
	tizi_local_proto_lvodQuery "github.com/tizi-local/lvodQuery/proto/vodQuery"
	"github.com/tizi-local/lvodQuery/config"
	"github.com/tizi-local/lvodQuery/internal/db"
	"github.com/tizi-local/lvodQuery/internal/rpc"
	"github.com/tizi-local/lvodQuery/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
)

var (
	COMMIT_SHA1 string
	BUILD_TIME  string
	VERSION     string
)

//flags
var (
	v          bool
	configPath string
)

var (
	appConfig config.Config
)

func parse() {
	flag.BoolVar(&v, "version", false, "show version")
	flag.StringVar(&configPath, "config", "/app/lvodquery.confg", "lvodquery config file")
	flag.Parse()
}

func main() {
	parse()
	if v {
		fmt.Println("COMMIT_SHA1:", COMMIT_SHA1)
		fmt.Println("VERSION:", VERSION)
		fmt.Println("BUILD_TIME:", BUILD_TIME)
		flag.Usage()
		return
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}
	if err := json.Unmarshal(configBytes, &appConfig); err != nil {
		fmt.Printf("read data:%s\n", string(configBytes))
		panic(err)
	}
	fmt.Println(appConfig.VodQueryRpc)
	//pprof
	go func() {
		err := http.ListenAndServe(appConfig.PprofAddr, nil)
		if err != nil {
			panic(err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())
	// init by config file
	logger := log.NewLogger(appConfig.Logger)
	db.InitDBEngine(appConfig.DB, logger)
	cache.InitCacheService(appConfig.Redis,logger)
	s := grpc.NewServer()
	conn, err := net.Listen("tcp", appConfig.VodQueryRpc.Addr)
	if err != nil {
		logger.Error("listen grpc port failed")
		return
	}
	vodQueryService := rpc.NewVodQueryService(logger)
	tizi_local_proto_lvodQuery.RegisterVodQueryServiceServer(s,vodQueryService)
	if err := s.Serve(conn); err != nil {
		logger.Error("serve auth rpc failed:", err.Error())
	}
}
