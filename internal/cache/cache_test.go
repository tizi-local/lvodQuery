package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/tizi-local/lvodQuery/config"
	"github.com/tizi-local/lvodQuery/pkg/models"
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	str, _ := jsoniter.MarshalToString(&models.CommentFirst{
		Vid: "123",
	})
	res, err := Set(context.Background(), "abc", str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("res", res)
}

func TestMGet(t *testing.T) {
	res, err := MGet(context.Background(), "a", "b")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("res", res)
}

func TestMain(m *testing.M) {
	InitCacheService(&config.RedisConfig{
		Password: "",
		Addr:     "127.0.0.1:6379",
	}, nil)
	os.Exit(m.Run())
}

func TestZRange(t *testing.T) {
	//c := make([]string, 0)a
	res := Default().ZRange(context.Background(), "zset1", 0, 1).Val()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	t.Log("res", res)
}

func TestHMGet(t *testing.T) {
	str, _ := jsoniter.MarshalToString(&models.CommentFirst{
		Vid: "123",
	})
	hsetKey := "hash222"
	_, err := HSet(context.Background(), hsetKey, "abc", str)
	if err != nil {
		t.Fatal(err)
	}

	data, err := HMGet(context.Background(), hsetKey, "abc", "bbc")
	assert.NoError(t, err, "get error %v", err)
	for i := range data{
		d, ok := data[i].(string)
		if ok {
			t.Log("d:", d)
			c := models.CommentFirst{}
			err := jsoniter.UnmarshalFromString(d, &c)
			if err != nil{
				t.Error("unmarshal error", err.Error())
				continue
			}
			t.Log("c:", c)
		}
	}

}