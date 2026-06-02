package test

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHandler is a mock implementation of http.Handler used for testing HTTP handlers.
// It embeds testify/mock to allow method call expectations and assertions.
// When ServeHTTP is called, it records the call with the provided ResponseWriter and Request,
// then automatically writes an HTTP 200 OK status to the response.
type MockHandler struct {
	mock.Mock
}

// ServeHTTP satisfies the http.Handler interface.
// It records the call with the given ResponseWriter and Request for test verification,
// and writes a 200 OK status code to the response.
func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}
