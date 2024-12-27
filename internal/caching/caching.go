package caching

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/valkey-io/valkey-go"
)

type Cache interface {
	SetString(ctx context.Context, key string, value string, options ...SetOption) error
	GetString(ctx context.Context, key string) (string, bool, error)
	Exists(ctx context.Context, key string) (bool, error)
	GetThread(ctx context.Context, key string) (*bsky.FeedGetPostThread_Output, bool, error)
	SetThread(ctx context.Context, key string, thread *bsky.FeedGetPostThread_Output, options ...SetOption) error
}

type SetOption func(opt *setopt)

func WithTTL(ttl time.Duration) func(opt *setopt) {
	return func(opt *setopt) {
		opt.TTL = ttl
	}
}

type setopt struct {
	TTL time.Duration
}

type ValkeyCache struct {
	valkeyClient      valkey.Client
	defaultExpiration time.Duration
	prefix            string
}

type Configuration struct {
	DefaultExpiration time.Duration
	Prefix            string
	ValkeyClient      valkey.Client
}

func NewValkeyCache(cfg Configuration) *ValkeyCache {
	return &ValkeyCache{
		valkeyClient:      cfg.ValkeyClient,
		prefix:            cfg.Prefix,
		defaultExpiration: cfg.DefaultExpiration,
	}
}

func (c *ValkeyCache) GetString(ctx context.Context, key string) (string, bool, error) {
	cmd := c.valkeyClient.B().Get().Key(key).Build()
	result := c.valkeyClient.Do(ctx, cmd)
	if err := result.Error(); err != nil {
		return "", false, err
	}
	if !result.IsCacheHit() {
		return "", false, nil
	}
	return result.String(), true, nil
}

func (c *ValkeyCache) SetString(ctx context.Context, key string, val string, options ...SetOption) error {
	opts := setopt{}
	for _, opt := range options {
		opt(&opts)
	}
	var cmd valkey.Completed
	if opts.TTL > 0 {
		cmd = c.valkeyClient.B().Set().Key(key).Value(val).Ex(opts.TTL).Build()
	} else {
		cmd = c.valkeyClient.B().Set().Key(key).Value(val).Build()
	}
	result := c.valkeyClient.Do(ctx, cmd)
	return result.Error()
}

func (c *ValkeyCache) GetThread(ctx context.Context, key string) (*bsky.FeedGetPostThread_Output, bool, error) {
	cmd := c.valkeyClient.B().Get().Key(key).Build()
	result := c.valkeyClient.Do(ctx, cmd)
	var obj *bsky.FeedGetPostThread_Output
	if err := result.Error(); err != nil {
		return obj, false, err
	}
	if !result.IsCacheHit() {
		return obj, false, nil
	}
	if err := result.DecodeJSON(&obj); err != nil {
		return obj, false, err
	}
	return obj, true, nil
}

func (c *ValkeyCache) SetThread(ctx context.Context, key string, val *bsky.FeedGetPostThread_Output, options ...SetOption) error {
	opts := setopt{}
	for _, opt := range options {
		opt(&opts)
	}
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	var cmd valkey.Completed
	if opts.TTL > 0 {
		cmd = c.valkeyClient.B().Set().Key(key).Value(string(raw)).Ex(opts.TTL).Build()
	} else {
		cmd = c.valkeyClient.B().Set().Key(key).Value(string(raw)).Build()
	}
	result := c.valkeyClient.Do(ctx, cmd)
	return result.Error()
}

func (c *ValkeyCache) Exists(ctx context.Context, key string) (bool, error) {
	cmd := c.valkeyClient.B().Exists().Key(key).Build()
	return c.valkeyClient.Do(ctx, cmd).AsBool()
}
