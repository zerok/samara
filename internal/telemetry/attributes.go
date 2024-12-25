package telemetry

import "go.opentelemetry.io/otel/attribute"

const ThreadCacheHitKey = attribute.Key("samara.thread.cache_hit")
const ThreadURIKey = attribute.Key("samara.thread.uri")
