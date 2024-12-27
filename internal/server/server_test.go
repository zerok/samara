package server

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/stretchr/testify/require"
)

func TestThreadRendering(t *testing.T) {
	thread := getOrCreateTestThread(t)
	require.NotNil(t, thread)
	srv := New(Configuration{})

	t.Run("no error", func(t *testing.T) {
		_, _, err := srv.renderThread(context.Background(), thread, 0)
		require.NoError(t, err)
	})
}

const testThreadFile = "testdata/thread.json"

func getOrCreateTestThread(t *testing.T) *bsky.FeedDefs_ThreadViewPost {
	t.Helper()
	_, err := os.Stat(testThreadFile)
	var thread *bsky.FeedDefs_ThreadViewPost
	if os.IsNotExist(err) {
		threadURI := "at://did:plc:rcygvu3gobjenoognsjq3y4q/app.bsky.feed.post/3lar526s5u22s"
		client := xrpc.Client{
			Host: "https://public.api.bsky.app",
		}
		tout, err := bsky.FeedGetPostThread(context.Background(), &client, 5, 0, threadURI)
		require.NoError(t, err)
		thread = tout.Thread.FeedDefs_ThreadViewPost
		fp, err := os.Create(testThreadFile)
		require.NoError(t, err)
		defer fp.Close()
		enc := json.NewEncoder(fp)
		enc.SetIndent("", "  ")
		require.NoError(t, enc.Encode(tout.Thread))
	} else {
		fp, err := os.Open(testThreadFile)
		require.NoError(t, err)
		defer fp.Close()
		require.NoError(t, json.NewDecoder(fp).Decode(&thread))
	}
	return thread
}
