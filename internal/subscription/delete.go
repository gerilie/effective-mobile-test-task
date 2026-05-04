package subscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// delete handles HTTP request for deleting a subscription by ID.
//
// It returns 204 No Content on success.
func (s *server) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	err := s.service.delete(ctx, id)
	if err != nil {
		handleServiceErrors(ctx, w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info(ctx, "subscription deleted", zap.String("id", id))
}

func (s *service) delete(ctx context.Context, id string) error {
	return s.repo.delete(ctx, id)
}

func (r *pgRepository) delete(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)

	sql, args, err := r.builder.Delete("subscriptions").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	row := r.db.QueryRow(ctx, sql, args...)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("read row: %w", err)
	}

	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return nil
}
