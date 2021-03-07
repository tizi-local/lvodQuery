package feed

import (
	"fmt"
	ljson "github.com/tizi-local/llib/encoding/json"
	"github.com/tizi-local/lvodQuery/pkg/models"
	"testing"
)

func TestFeed_RandFeed(t *testing.T) {
	f := Feed{
		Videos: []*models.VideoInfo{
		},
		Session: "lvsession111",
	}
	for i := 0; i < 10;i++{
		f.Videos = append(f.Videos, &models.VideoInfo{
			Vid: fmt.Sprintf("%d", i),
		})
	}
	f.RandFeed()
	t.Log(ljson.MustMarshalIndentToString(f))
}
