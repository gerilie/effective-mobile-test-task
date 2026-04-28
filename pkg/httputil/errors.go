package httputil

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func HandleDefaultErrors(ctx context.Context, w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError
	log := logger.FromContext(ctx)

	switch {
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		log.Warn(ctx, "request canceled", zap.Error(err))
		http.Error(w, "request canceled", http.StatusGatewayTimeout)

	case errors.Is(err, sql.ErrNoRows):
		log.Warn(ctx, "resource not found", zap.Error(err))
		http.Error(w, "resource not found", http.StatusNotFound)

	case errors.As(err, &pgErr):
		log.Error(ctx, "database error",
			zap.String("code", pgErr.Code),
			zap.String("detail", pgErr.Detail),
			zap.String("table", pgErr.TableName),
			zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

	default:
		log.Error(ctx, "unexpected error", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
