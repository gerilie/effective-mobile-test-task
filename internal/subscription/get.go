package subscription

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

// @Summary		Get subscription
// @Description	Get by subscription ID.
// @Description	If subscription is not found, returns 404.
// @Tags			subscription
// @ID				get-subscription
// @Produce		json
// @Param			id	path		string	true	"Subscription ID"
// @Success		200	{object}	SubResp	"Subscription"
// @Failure		400	{string}	string	"Bad request"
// @Failure		404	{string}	string	"Not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/subscriptions/{id} [get].
func (s *server) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	sub, err := s.service.get(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "subscription is not found", http.StatusNotFound)

			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Info(ctx, "subscription found and sent", zap.String("id", id))
}

func (s *service) get(ctx context.Context, id string) (SubResp, error) {
	return s.repo.get(ctx, id)
}

func (r *pgRepository) get(ctx context.Context, id string) (SubResp, error) {
	log := logger.FromContext(ctx)

	sqlStr, args, err := r.builder.Select("id, service_name, price, user_id, start_date, end_date").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Error(ctx, postgres.ErrBuildingQuery.Error(), zap.Error(err))

		return SubResp{}, err
	}

	row := r.db.QueryRow(ctx, sqlStr, args...)

	var startDate, endDate sql.NullTime
	var sub SubResp
	if err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &startDate, &endDate); err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return SubResp{}, err
	}

	if startDate.Valid {
		sub.StartDate = startDate.Time.Format("01-2006")
	}
	if endDate.Valid {
		endStr := endDate.Time.Format("01-2006")
		sub.EndDate = &endStr
	}

	log.Info(ctx, "query executed", zap.String("query", sqlStr), zap.Any("args", args))

	return sub, nil
}
