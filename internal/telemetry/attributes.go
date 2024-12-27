package telemetry

import "go.opentelemetry.io/otel/attribute"

const ThreadCacheHitKey = attribute.Key("samara.thread.cache_hit")
const ThreadURIKey = attribute.Key("samara.thread.uri")
const AvatarDIDKey = attribute.Key("samara.avatar.did")
const AvatarUpstreamURLKey = attribute.Key("samara.avatar.upstream_url")
