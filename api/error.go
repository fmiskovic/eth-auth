package api

import (
	"encoding/json"
	"fmt"
)

// Error represents an error that can be returned by the API centralized error handling.
type Error struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal error  `json:"error"`
}

func newError(code int, message string, internal error) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Internal: internal,
	}
}

func (e Error) Error() string {
	jsonErr, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"message":"%s"}`, e.Message)
	}
	return string(jsonErr)
}
