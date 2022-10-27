package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDBConnection(t *testing.T) {
	conn, errCo := GetDBConnection()
	require.NoError(t, errCo)
	require.NotNil(t, conn)
}
