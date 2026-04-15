package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestSetupOTel(t *testing.T) {
	t.Run("no-config", func(t *testing.T) {
		shutdown, err := setupOTelSDK(t.Context(), "")
		require.NoError(t, err)
		shutdown(t.Context())
	})
	t.Run("invalid-path", func(t *testing.T) {
		_, err := setupOTelSDK(t.Context(), "/does/not-exist")
		require.Error(t, err)
	})

	t.Run("valid-config", func(t *testing.T) {
		dir := t.TempDir()
		os.WriteFile(
			filepath.Join(dir, "test.conf"),
			[]byte(`file_format: "0.3"
disabled: false
meter_provider:
  readers:
  - periodic:
      exporter:
        console: {}
`),
			0600,
		)
		shutdown, err := setupOTelSDK(t.Context(), filepath.Join(dir, "test.conf"))
		require.NoError(t, err)
		defer shutdown(t.Context())
		meter := otel.Meter("test")
		counter, err := meter.Int64Counter("counter")
		require.NoError(t, err)
		counter.Add(t.Context(), 1)
	})
}
