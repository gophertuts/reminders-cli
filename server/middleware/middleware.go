package middleware

import (
	"net/http"
)

// New creates a new instance of Middleware chains
func New(ms ...func(h http.Handler) http.Handler) *Middleware {
	return &Middleware{
		functions: ms,
	}
}

// Middleware represents HTTP router middleware type
type Middleware struct {
	functions []func(h http.Handler) http.Handler
}

// Then runs the request through the middleware chain, then serves it
func (m *Middleware) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for i := range m.functions {
		// each middleware returns an http.Handler
		// h = M[n-1](...(M[0](controller)))
		h = m.functions[len(m.functions)-1-i](h)
	}
	return h
}
