package httputil

import "errors"

var (
	ErrRequestCanceled = errors.New("request canceled")

	ErrReadingRequestBody  = errors.New("reading request body")
	ErrWritingResponseBody = errors.New("writing response body")
	ErrClosingRequestBody  = errors.New("closing request body")

	ErrEncodingResponseBody = errors.New("encoding response body")
	ErrDecodingRequestBody  = errors.New("decoding request body")
)
