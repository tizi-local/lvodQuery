package feed

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/tizi-local/lvodQuery/pkg/models"
	"math/rand"
	"time"
)

const (
	FeedCountLimit = 100
	FeedPageSize   = 10

)

var (
	ErrorLastOfSession = errors.New("LastOfSession")
)

type FeedGenerator struct {
	idGenerator *snowflake.Node
}

type Feed struct {
	Videos  []*models.VideoInfo
	Session string
}

func NewFeedGenerator() *FeedGenerator {
	node, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return &FeedGenerator{
		idGenerator: node,
	}
}

func (fg *FeedGenerator) NewSession() string {
	return fmt.Sprintf("lls_%s", fg.idGenerator.Generate())
}

func (fg *FeedGenerator) NewFeed(videos []*models.VideoInfo) Feed {
	f := Feed{
		Videos: videos,
		Session: fg.NewSession(),
	}
	f.RandFeed()
	return f
}

func (f *Feed) RandFeed()  {
	f.Videos = sliceOutOfOrder(f.Videos)
}

func sliceOutOfOrder(in []*models.VideoInfo) []*models.VideoInfo {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := r.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}
