# Samara: Bluesky comments for static websites

**Warning:** This is currently experimental and mostly written to work on my own website.
The API *will* change so use at your own risk.

Welcome to Samara, a simple server for fetching Bluesky thread information to be rendered by weblogs and other websites.

For a concrete example on how to run Samara and integrate it via HTMX into your site, please take a look at the `examples/simple-htmx` folder.

## Configuration options

```
--addr string                           Address to listen on (default "0.0.0.0:80
--allowed-origin strings                Allowed origins (CORS)
--allowed-root-account-handle strings   Allowed root account handles
--base-url string                       Base URL to be used for linking inside the results (default "http://localhost:8080")
--log-level string                      Log level (debug, info, warn, error) (default "warn")
--otel-config string                    Path to an OpenTelemetry configuration file
--valkey-hostname string                Hostname of a valkey caching server
--version                               Print version information
```

## Caching

Samara by default caches various request to the Bluesky API to improve performance.
This cache is kept within the main process so it won't survive a restart of the application.
To get around this, Samara also allows caching through a [Valkey][] instance which you can configure using the `--valkey-hostname` flag.

## OpenTelemetry integration

Samara generates OpenTelemetry tracing data.
In order to activate this, specify an OpenTelemetry configuration file with the `--otel-config` flag.
You can find a sample file inside the `examples/simple-htmx` folder.

## The name

<https://en.wikipedia.org/wiki/Samara_(fruit)>

[valkey]: https://valkey.io/
