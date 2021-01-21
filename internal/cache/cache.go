package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tizi-local/llib/log"
	"github.com/tizi-local/lvodQuery/config"
)

type CacheService struct {
	*redis.Client
	*log.Logger
}

var (
	s   *CacheService
	mux sync.RWMutex
)

func InitCacheService(c *config.RedisConfig, logger *log.Logger) *CacheService {
	mux.Lock()
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password, // no password set
		DB:           0,          // use default DB
		MinIdleConns: 10,
	})
	res, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("res: %s, err: %v \n", res, err)
	}
	//logger.Infof("res %s, err %v", res, err)
	s = &CacheService{rdb, logger}
	mux.Unlock()

	return s
}

func Default() *CacheService {
	return s
}

const (
	commonPrefix = "vod"
)

func SAdd(ctx context.Context, key string, value []byte) (int64, error) {
	return Default().SAdd(ctx, key, value).Result()
}
func SNum(ctx context.Context, key string) (int64, error) {
	return Default().SCard(ctx, key).Result()
}
func SPop(ctx context.Context, key string) (string, error) {
	return Default().SPop(ctx, key).Result()
}
func SExpire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return Default().Expire(ctx, key, ttl).Result()
}
func SGet(ctx context.Context, key string) (string, error) {
	return Default().SRandMember(ctx, key).Result()
}
func ZAdd(ctx context.Context,key string,value *redis.Z)(int64,error){
	return Default().ZAdd(ctx,key,value).Result()
}
func ZRange(ctx context.Context,key string,container [][]byte,start,stop int64)error{
	return Default().ZRange(ctx,key,start,stop).ScanSlice(container)
}
func Exist(ctx context.Context,key string)int64{
	return Default().Exists(ctx,key).Val()
}
func ZNum(ctx context.Context,key string)(int64,error){
	return Default().ZCard(ctx,key).Result()
}