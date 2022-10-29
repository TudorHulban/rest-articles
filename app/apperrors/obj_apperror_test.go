package apperrors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringer(t *testing.T) {
	errApp := ErrorApplication{
		Area:    Areas[ErrorAreaInfra],
		Message: "this is error message",
		Code:    "INFRA-1234",
		OSExit:  &OSExitForDatabaseIssues,
	}

	fmt.Println(errApp)
}

func TestMarshaler(t *testing.T) {
	errApp := ErrorApplication{
		Area:    Areas[ErrorAreaInfra],
		Message: "this is error message",
		Code:    "INFRA-1234",
		OSExit:  &OSExitForDatabaseIssues,
	}

	raw, errMa := json.Marshal(errApp)
	require.NoError(t, errMa)

	fmt.Println(string(raw))
}