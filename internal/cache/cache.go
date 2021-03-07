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

// set
func Set(ctx context.Context, key, value string) (string, error) {
	return Default().Set(ctx, key, value, -1).Result()
}

// multi get
func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return Default().MGet(ctx, keys...).Result()
}

// set and expire
func SetExpire(ctx context.Context, key, value string, ttl time.Duration) (string, error) {
	return Default().Set(ctx, key, value, ttl).Result()
}

// setNx and expire
func SetNXExpire(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	return Default().SetNX(ctx, key, value, ttl).Result()
}

func Del(ctx context.Context, keys ...string) (int64, error) {
	return Default().Del(ctx, keys...).Result()
}

// expire
func Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return Default().Expire(ctx, key, ttl).Result()
}

// exist
func Exist(ctx context.Context, key string) int64 {
	return Default().Exists(ctx, key).Val()
}

// set
func SAdd(ctx context.Context, key string, value []byte) (int64, error) {
	return Default().SAdd(ctx, key, value).Result()
}

func SCard(ctx context.Context, key string) (int64, error) {
	return Default().SCard(ctx, key).Result()
}
func SPop(ctx context.Context, key string) (string, error) {
	return Default().SPop(ctx, key).Result()
}

func SGet(ctx context.Context, key string) (string, error) {
	return Default().SRandMember(ctx, key).Result()
}

//zset
func ZAdd(ctx context.Context, key string, value *redis.Z) (int64, error) {
	return Default().ZAdd(ctx, key, value).Result()
}
func ZRange(ctx context.Context, key string, start, stop int64) []string {
	return Default().ZRange(ctx, key, start, stop).Val()
}

// zrevrange return [start, stop]
func ZREVRange(ctx context.Context, key string, start, stop int64) []string {
	return Default().ZRevRange(ctx, key, start, stop).Val()
}

func ZNum(ctx context.Context, key string) (int64, error) {
	return Default().ZCard(ctx, key).Result()
}

// list
func RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return Default().RPush(ctx, key, values).Result()
}

func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return Default().LRange(ctx, key, start, stop).Result()
}

func LLen(ctx context.Context, key string) (int64, error) {
	return Default().LLen(ctx, key).Result()
}

// hash
func HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return Default().HSet(ctx, key, values).Result()
}

func HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return Default().HMGet(ctx, key, fields...).Result()
}
