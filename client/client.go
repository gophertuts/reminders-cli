package client

import "fmt"

// wrapError wraps a plain error into a custom error
func wrapError(customErr string, originalErr error) error {
	return fmt.Errorf("%s: %v", customErr, originalErr)
}
