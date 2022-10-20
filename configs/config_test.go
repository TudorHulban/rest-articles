package configs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := NewParsedConfig("../.env")
	require.NoError(t, err)
	require.Equal(t, _HTTPServerPort, cfg.HTTPServerPort)
}
