package subscription

import (
	"context"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

// @Summary		Delete subscription
// @Description	Delete by subscription ID.
// @Tags			subscription
// @ID				delete-subscription
// @Param			id	path	string	true	"Subscription ID"
// @Success		204	"no content"
// @Failure		400	{string}	string	"Bad request"
// @Failure		404	{string}	string	"Not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/subscriptions/{id} [delete].
func (s *server) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	err := s.service.delete(ctx, id)
	if err != nil {
		httputil.HandleErrors(ctx, w, err)

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
		log.Error(ctx, postgres.ErrBuildingQuery.Error(), zap.Error(err))

		return err
	}
	var deletedID string
	row := r.db.QueryRow(ctx, sql, args...)
	if err := row.Scan(&deletedID); err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return err
	}

	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return nil
}
