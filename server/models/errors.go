package models

import "fmt"

// WrapError wraps a plain error into a custom error
func WrapError(customErr string, originalErr error) error {
	return fmt.Errorf("%s: %v", customErr, originalErr)
}
