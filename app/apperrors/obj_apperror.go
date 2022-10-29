package apperrors

import (
	"encoding/json"
	"fmt"
)

type ErrorArea string

type ErrorApplication struct {
	Area    ErrorArea
	Message string
	Code    string
	OSExit  *int
}

var _ error = &ErrorApplication{}

func (e ErrorApplication) String() string {
	if e.OSExit == nil {
		return fmt.Sprintf("Area: %s\nCode: %s\nMessage: %s\n", e.Area, e.Code, e.Message)
	}

	return fmt.Sprintf("Area: %s\nCode: %s\nMessage: %s\nOS Exit Code: %d", e.Area, e.Code, e.Message, *e.OSExit)
}

func (e ErrorApplication) Error() string {
	return e.String()
}

func (e ErrorApplication) MarshalJSON() ([]byte, error) {
	res := make(map[string]string)

	res["area"] = string(e.Area)
	res["errormessage"] = e.Message
	res["errorcode"] = e.Code

	return json.Marshal(&res)
}
