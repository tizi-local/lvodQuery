package feed

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/tizi-local/lvodQuery/pkg/models"
)

const (
	FeedCountLimit = 1000
	FeedPageSize   = 10
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

func (f *FeedGenerator) GenerateSession() string {
	return fmt.Sprintf("lls_%s", f.idGenerator.Generate())
}

func (f *FeedGenerator) GenerateFeed() Feed {
	return Feed{}
}
