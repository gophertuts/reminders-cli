package models

import "fmt"

// HTTPError represents an http error to be returned to the client
type HTTPError struct {
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

// WrapError wraps a plain error into a custom error
func WrapError(customErr string, originalErr error) error {
	return fmt.Errorf("%s: %v", customErr, originalErr)
}
