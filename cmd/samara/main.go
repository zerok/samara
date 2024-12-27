package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/valkey-io/valkey-go"
	"github.com/zerok/samara/internal/caching"
	"github.com/zerok/samara/internal/server"
	config "go.opentelemetry.io/contrib/config/v0.3.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var version string

func main() {
	var addr string
	var otelConfigFile string
	var allowedRootAccountHandles []string
	var allowedOrigins []string
	var logLevel string
	var showVersion bool
	var baseURL string
	var valkeyHostname string

	pflag.StringVar(&addr, "addr", "0.0.0.0:8080", "Address to listen on")
	pflag.StringSliceVar(&allowedRootAccountHandles, "allowed-root-account-handle", []string{}, "Allowed root account handles")
	pflag.StringSliceVar(&allowedOrigins, "allowed-origin", []string{}, "Allowed origins (CORS)")
	pflag.StringVar(&logLevel, "log-level", "warn", "Log level (debug, info, warn, error)")
	pflag.StringVar(&otelConfigFile, "otel-config", "", "Path to an OpenTelemetry configuration file")
	pflag.BoolVar(&showVersion, "version", false, "Print version information")
	pflag.StringVar(&baseURL, "base-url", "http://localhost:8080", "Base URL to be used for linking inside the results")
	pflag.StringVar(&valkeyHostname, "valkey-hostname", "", "Hostname of a valkey caching server")
	pflag.Parse()

	if version == "" {
		version = "unknown"
	}

	if showVersion {
		fmt.Println(version)
		return
	}

	lvl := slog.LevelWarn
	switch logLevel {
	case "error":
		lvl = slog.LevelError
	case "info":
		lvl = slog.LevelInfo
	case "debug":
		lvl = slog.LevelDebug
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

	var cache caching.Cache
	var err error
	if valkeyHostname != "" {
		valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{valkeyHostname}})
		if err != nil {
			logger.ErrorContext(ctx, "valkey connection failed", "err", err)
			os.Exit(1)
		}
		defer valkeyClient.Close()
		cache = caching.NewValkeyCache(caching.Configuration{
			ValkeyClient:      valkeyClient,
			DefaultExpiration: time.Minute * 5,
		})
	} else {
		cache = caching.NewLocalCache()
	}

	httpSrv := http.Server{}

	shutdownOtel, err := setupOTelSDK(ctx, otelConfigFile)
	if err != nil {
		logger.ErrorContext(ctx, "otel setup failed", "err", err)
		os.Exit(1)
	}
	defer shutdownOtel(context.Background())

	go func() {
		<-ctx.Done()
		cancel()
		httpSrv.Shutdown(context.Background())
	}()

	cfg := server.Configuration{
		AllowedRootAccountDIDs: make([]string, 0, 5),
		AllowedRootAccounts:    make([]string, 0, 5),
		Logger:                 logger,
		BaseURL:                baseURL,
		Cache:                  cache,
	}

	// Resolve handles to DIDs
	client := &xrpc.Client{}
	client.Host = "https://public.api.bsky.app"

	for _, handle := range allowedRootAccountHandles {
		did, found, _ := cache.GetString(ctx, "handle-did:"+handle)
		if !found {
			result, err := atproto.IdentityResolveHandle(ctx, client, handle)
			if err != nil {
				logger.ErrorContext(ctx, "failed to resolve handle", "handle", handle, "err", err)
				os.Exit(1)
			}
			did = result.Did
			cache.SetString(ctx, "handle-did:"+handle, result.Did)
		}
		cfg.AllowedRootAccountDIDs = append(cfg.AllowedRootAccountDIDs, did)
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

func setupOTelSDK(ctx context.Context, configPath string) (shutdown func(context.Context) error, err error) {
	opts := make([]config.ConfigurationOption, 0, 2)
	opts = append(opts, config.WithContext(ctx))

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		otelConfig, err := config.ParseYAML(data)
		if err != nil {
			return nil, err
		}

		// We want to override the resource definition to keep name and version
		// managed by the build process:
		schemaURL := semconv.SchemaURL
		res := config.Resource{
			SchemaUrl: &schemaURL,
			Attributes: []config.AttributeNameValue{
				{
					Name:  "service.name",
					Value: "samara",
				},
				{
					Name:  "service.version",
					Value: version,
				},
			},
		}
		otelConfig.Resource = &res
		opts = append(opts, config.WithOpenTelemetryConfiguration(*otelConfig))
	}

	sdk, err := config.NewSDK(opts...)
	if err != nil {
		return sdk.Shutdown, err
	}
	shutdown = sdk.Shutdown

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider := sdk.TracerProvider()
	otel.SetTracerProvider(tracerProvider)

	meterProvider := sdk.MeterProvider()
	otel.SetMeterProvider(meterProvider)
	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
