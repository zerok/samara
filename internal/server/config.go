package server

import (
	"log/slog"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/zerok/samara/internal/caching"
)

type Configuration struct {
	AllowedRootAccountDIDs []string
	AllowedRootAccounts    []string
	Client                 *xrpc.Client
	Logger                 *slog.Logger
	BaseURL                string
	Cache                  caching.Cache
}
