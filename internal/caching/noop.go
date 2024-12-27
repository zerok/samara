package caching

import (
	"context"

	"github.com/bluesky-social/indigo/api/bsky"
)

type NoopCache struct{}

func NewNoopCache() *NoopCache {
	return &NoopCache{}
}

func (c *NoopCache) GetString(_ context.Context, _ string) (string, bool, error) {
	return "", false, nil
}

func (c *NoopCache) SetString(_ context.Context, _ string, _ string, _ ...SetOption) error {
	return nil
}

func (c *NoopCache) Exists(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (c *NoopCache) SetThread(_ context.Context, _ string, _ *bsky.FeedGetPostThread_Output, _ ...SetOption) error {
	return nil
}

func (c *NoopCache) GetThread(_ context.Context, _ string) (*bsky.FeedGetPostThread_Output, bool, error) {
	return nil, false, nil
}
