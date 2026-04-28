package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
)

func DecodeJSON[T any](ctx context.Context, w http.ResponseWriter, r *http.Request, dst *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}
	defer deferfunc.Close(ctx, r.Body.Close, "close request body")

	if err := json.Unmarshal(body, dst); err != nil {
		var unmarshalErr *json.UnmarshalTypeError
		var syntaxErr *json.SyntaxError
		if errors.As(err, &unmarshalErr) || errors.As(err, &syntaxErr) {
			http.Error(
				w,
				"invalid request body: type mismatch or malformed JSON",
				http.StatusBadRequest,
			)

			return fmt.Errorf("decode request body: %w", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return fmt.Errorf("decode request body: %w", err)
	}

	return nil
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, status int, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return fmt.Errorf("encode request body: %w", err)
	}

	w.Header().Set(ContentType, JSON)
	w.WriteHeader(status)

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return fmt.Errorf("write response body: %w", err)
	}

	return nil
}
