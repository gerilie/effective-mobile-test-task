package httputil

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func DecodeJSON[T any](ctx context.Context, w http.ResponseWriter, r *http.Request, dst *T) error {
	log := logger.FromContext(ctx)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrReadingRequestBody.Error(), http.StatusBadRequest)
		log.Error(ctx, ErrReadingRequestBody.Error(), zap.Error(err))

		return err
	}
	defer deferfunc.Close(ctx, r.Body.Close, ErrClosingRequestBody.Error())

	if err := json.Unmarshal(body, dst); err != nil {
		http.Error(w, ErrDecodingRequestBody.Error(), http.StatusBadRequest)
		log.Error(ctx, ErrDecodingRequestBody.Error(), zap.Error(err))

		return err
	}

	return nil
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, status int, v any) error {
	log := logger.FromContext(ctx)

	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, ErrEncodingResponseBody.Error(), http.StatusInternalServerError)
		log.Error(ctx, ErrEncodingResponseBody.Error(), zap.Error(err))

		return err
	}

	w.Header().Set(ContentType, JSON)
	w.WriteHeader(status)

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, ErrWritingResponseBody.Error(), http.StatusInternalServerError)
		log.Error(ctx, ErrWritingResponseBody.Error(), zap.Error(err))

		return err
	}

	return nil
}
