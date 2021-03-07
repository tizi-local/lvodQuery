package cache

import "strings"

const (
	VodKeyFeedSession = "vod_feed_session"
	VodKeyVideo       = "vod_video"

	VodCommentKeyFirst    = "comment_first"    //hash comment meta data
	VodCommentKeySecond   = "comment_second"   //hash sub comment meta data
	VodCommentKeyTimeline = "comment_timeline" // zset for comment time with commentid
	VodCommentScore       = "comment_score"    // zset for comment score with commentid
	VodCommentWriteMutex  = "comment_mutex"    // mutex for write comment data cache

)

func Key(prefix, key string) string {
	return strings.Join([]string{commonPrefix, prefix, key}, ":")
}
