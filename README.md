# Samara: Bluesky comments for static websites

**Warning:** This is currently experimental and mostly written to work on my own website.
The API *will* change so use at your own risk.

Welcome to Samara, a simple server for fetching Bluesky thread information to be rendered by weblogs and other websites.

For a concrete example on how to run Samara and integrate it via HTMX into your site, please take a look at the `examples/simple-htmx` folder.

## Configuration options

```
--addr string                           Address to listen on (default "0.0.0.0:8080")
--allowed-origin strings                Allowed origins (CORS)
--allowed-root-account-handle strings   Allowed root account handles
--log-level string                      Log level (debug, info, warn, error) (default "warn")
--otel-config string                    Path to an OpenTelemetry configuration file
--version                               Print version information

```

## OpenTelemetry integration

Samara generates OpenTelemetry tracing data.
In order to activate this, specify an OpenTelemetry configuration file with the `--otel-config` flag.
You can find a sample file inside the `examples/simple-htmx` folder.

## The name

<https://en.wikipedia.org/wiki/Samara_(fruit)>
