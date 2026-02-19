package httputil

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrRequestCanceled = errors.New("request canceled")

	ErrReadingRequestBody  = errors.New("reading request body")
	ErrWritingResponseBody = errors.New("writing response body")
	ErrClosingRequestBody  = errors.New("closing request body")

	ErrEncodingResponseBody = errors.New("encoding response body")
	ErrDecodingRequestBody  = errors.New("decoding request body")
)

func HandleErrors(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
