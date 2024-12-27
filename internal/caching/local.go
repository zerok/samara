package caching

import (
	"context"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

type LocalCache struct {
	stringCache *expirable.LRU[string, string]
	threadCache *expirable.LRU[string, *bsky.FeedGetPostThread_Output]
}

func NewLocalCache() *LocalCache {
	stringCache := expirable.NewLRU[string, string](1000, nil, time.Minute)
	threadCache := expirable.NewLRU[string, *bsky.FeedGetPostThread_Output](10, nil, time.Minute*5)
	return &LocalCache{
		stringCache: stringCache,
		threadCache: threadCache,
	}
}
func (c *LocalCache) GetString(ctx context.Context, key string) (string, bool, error) {
	result, ok := c.stringCache.Get(key)
	return result, ok, nil
}

func (c *LocalCache) SetString(ctx context.Context, key string, val string, options ...SetOption) error {
	c.stringCache.Add(key, val)
	return nil
}

func (c *LocalCache) GetThread(ctx context.Context, key string) (*bsky.FeedGetPostThread_Output, bool, error) {
	thread, ok := c.threadCache.Get(key)
	return thread, ok, nil
}

func (c *LocalCache) SetThread(ctx context.Context, key string, val *bsky.FeedGetPostThread_Output, options ...SetOption) error {
	c.threadCache.Add(key, val)
	return nil
}

func (c *LocalCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok, err := c.GetString(ctx, key)
	return ok, err
}
