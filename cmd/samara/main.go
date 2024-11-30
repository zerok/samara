package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/zerok/samara/internal/server"
)

func main() {
	var addr string
	var allowedRootAccountHandles []string
	var allowedOrigins []string
	var logLevel string

	pflag.StringVar(&addr, "addr", "0.0.0.0:8080", "Address to listen on")
	pflag.StringSliceVar(&allowedRootAccountHandles, "allowed-root-account-handle", []string{}, "Allowed root account handles")
	pflag.StringSliceVar(&allowedOrigins, "allowed-origin", []string{}, "Allowed origins (CORS)")
	pflag.StringVar(&logLevel, "log-level", "warn", "Log level (debug, info, warn, error)")
	pflag.Parse()

	lvl := slog.LevelWarn
	switch logLevel {
	case "error":
		lvl = slog.LevelError
	case "info":
		lvl = slog.LevelInfo
	case "debug":
		lvl = slog.LevelDebug
	}

	httpSrv := http.Server{}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		cancel()
		httpSrv.Shutdown(context.Background())
	}()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

	cfg := server.Configuration{
		AllowedRootAccountDIDs: make([]string, 0, 5),
		AllowedRootAccounts:    make([]string, 0, 5),
		Logger:                 logger,
	}

	// Resolve handles to DIDs
	client := &xrpc.Client{}
	client.Host = "https://public.api.bsky.app"

	for _, handle := range allowedRootAccountHandles {
		result, err := atproto.IdentityResolveHandle(ctx, client, handle)
		if err != nil {
			logger.ErrorContext(ctx, "failed to resolve handle", "handle", handle, "err", err)
			os.Exit(1)
		}
		cfg.AllowedRootAccountDIDs = append(cfg.AllowedRootAccountDIDs, result.Did)
		cfg.AllowedRootAccounts = append(cfg.AllowedRootAccounts, handle)
	}

	srv := server.New(cfg)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedHeaders: []string{"hx-current-url", "hx-request"},
	})
	httpSrv.Addr = addr
	httpSrv.Handler = corsHandler.Handler(srv)
	logger.Info("starting server", "addr", addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		logger.Error("server startup failed", "err", err)
		os.Exit(1)
	}
}
