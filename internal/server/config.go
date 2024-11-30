package server

import (
	"log/slog"

	"github.com/bluesky-social/indigo/xrpc"
)

type Configuration struct {
	AllowedRootAccountDIDs []string
	AllowedRootAccounts    []string
	Client                 *xrpc.Client
	Logger                 *slog.Logger
}
