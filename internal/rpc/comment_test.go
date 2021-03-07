package rpc

import "testing"

func TestParseCommentId(t *testing.T) {
	id := ParseCommentId("3458764513820540934")
	t.Log(id)
}
