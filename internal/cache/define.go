package cache

import "strings"

const (
	VodKeyFeedSession = "vod_feed_session"
	VodKeyVideo = "vod_video"
)

func Key(prefix, key string) string {
	return strings.Join([]string{commonPrefix, prefix, key}, ":")
}
