package server

import "github.com/bluesky-social/indigo/xrpc"

type Configuration struct {
	AllowedRootAccountDIDs []string
	Client                 *xrpc.Client
}
