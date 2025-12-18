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

// @Summary		Get subscription sum
// @Description	Get total_price of all subscriptions.
// @Description	Filter by service_name AND/OR user_id.
// @Description	Filter by start_date AND end_date. Date format: 01-2006.
// @Tags			subscription
// @ID				get-subscription-sum
// @Accept			json
// @Produce		json
// @Param			subSum	body		SubSumReq	true	"Subscription sum"
// @Success		200		{object}	SubSumResp	"Subscription sum"
// @Failure		400		{string}	string		"Bad request"
// @Failure		500		{string}	string		"Internal server error"
// @Router			/subscriptions/sum [get].
func (s *server) sum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var subSum SubSumReq
	if err := httputil.DecodeJSON(ctx, r, &subSum); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := s.validateSubSum(ctx, &subSum); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := s.service.sum(ctx, subSum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			http.Error(w, pgErr.Message, http.StatusBadRequest)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Info(ctx, "total price calculated")
}

func (s *service) sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error) {
	return s.repo.sum(ctx, subSum)
}

func (r *pgRepository) sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error) {
	log := logger.FromContext(ctx)

	qb := r.builder.Select("SUM(price) as total_price").
		From("subscriptions").
		Where(squirrel.GtOrEq{"start_date": subSum.StartDate}).
		Where(squirrel.LtOrEq{"end_date": subSum.EndDate})

	if subSum.ServiceName != nil {
		qb = qb.Where(squirrel.Eq{"service_name": subSum.ServiceName}).GroupBy("service_name")
	}
	if subSum.UserID != nil {
		qb = qb.Where(squirrel.Eq{"user_id": subSum.UserID}).GroupBy("user_id")
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		log.Error(ctx, postgres.ErrBuildingQuery.Error(), zap.Error(err))

		return SubSumResp{}, err
	}

	row := r.db.QueryRow(ctx, sql, args...)

	var resp SubSumResp
	if err := row.Scan(&resp.TotalPrice); err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return SubSumResp{}, err
	}
	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return resp, nil
}
