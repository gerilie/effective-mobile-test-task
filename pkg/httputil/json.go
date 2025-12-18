package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func DecodeJSON[T any](ctx context.Context, r *http.Request, dst *T) error {
	log := logger.FromContext(ctx)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(ctx, ErrReadingRequestBody.Error(), zap.Error(err))

		return fmt.Errorf("read body: %w", err)
	}
	defer deferfunc.Close(ctx, r.Body.Close, ErrClosingRequestBody.Error())

	if err := json.Unmarshal(body, dst); err != nil {
		log.Error(ctx, ErrDecodingRequestBody.Error(), zap.Error(err))

		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, status int, v any) error {
	log := logger.FromContext(ctx)

	data, err := json.Marshal(v)
	if err != nil {
		log.Error(ctx, ErrEncodingResponseBody.Error(), zap.Error(err))

		return err
	}

	w.Header().Set(ContentType, JSON)
	w.WriteHeader(status)

	_, err = w.Write(data)
	if err != nil {
		log.Error(ctx, ErrWritingResponseBody.Error(), zap.Error(err))
	}

	return err
}
