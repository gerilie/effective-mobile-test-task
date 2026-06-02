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

// InternalServerErrorMsg is the default error message for internal server errors.
const InternalServerErrorMsg = "internal server error"

// HandleDefaultErrors maps common application errors to HTTP responses
// and logs them using the logger stored in the context.
func HandleDefaultErrors(ctx context.Context, w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError
	log := logger.FromContext(ctx)

	switch {
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		log.Warn("request canceled", zap.Error(err))
		http.Error(w, "request canceled", http.StatusGatewayTimeout)

	case errors.Is(err, sql.ErrNoRows):
		log.Warn("resource not found", zap.Error(err))
		http.Error(w, "resource not found", http.StatusNotFound)

	case errors.As(err, &pgErr):
		log.Error("database error",
			zap.String("code", pgErr.Code),
			zap.String("detail", pgErr.Detail),
			zap.String("table", pgErr.TableName),
			zap.Error(err))
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)

	default:
		log.Error("unexpected error", zap.Error(err))
		http.Error(w, InternalServerErrorMsg, http.StatusInternalServerError)
	}
}
