package apperrors

import (
	"encoding/json"
	"fmt"
)

type ErrorArea string

type ErrorApplication struct {
	Area      ErrorArea
	AreaError error
	Code      string
	OSExit    *int
}

var _ error = &ErrorApplication{}

func (e ErrorApplication) String() string {
	if e.OSExit == nil {
		return fmt.Sprintf("Area: %s\nCode: %s\nMessage: %s\n", e.Area, e.Code, e.AreaError)
	}

	return fmt.Sprintf("Area: %s\nCode: %s\nMessage: %s\nOS Exit Code: %d", e.Area, e.Code, e.AreaError, *e.OSExit)
}

func (e ErrorApplication) Error() string {
	return e.String()
}

func (e ErrorApplication) MarshalJSON() ([]byte, error) {
	res := make(map[string]string, 3)

	res["area"] = string(e.Area)
	res["errormessage"] = e.AreaError.Error()
	res["errorcode"] = e.Code

	return json.Marshal(&res)
}
