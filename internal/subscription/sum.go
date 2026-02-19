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

// @Summary		Get subscription summa
// @Description	Get total_price of all subscriptions.
// @Tags			subscription
// @ID				get-subscription-sum
// @Produce		json
// @Param			start_date		query		string		true	"Date format: MM-YYYY"
// @Param			end_date		query		string		true	"Date format: MM-YYYY"
// @Param			service_name	query		string		false	"filter by service name"
// @Param			user_id			query		string		false	"filter by user ID"
// @Success		200				{object}	SubSumResp	"Subscription sum"
// @Failure		400				{string}	string		"Bad request"
// @Failure		500				{string}	string		"Internal server error"
// @Router			/subscriptions/sum [get].
func (s *server) sum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	query := r.URL.Query()

	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	subSum := SubSumReq{
		startDate:   startDate,
		endDate:     endDate,
		serviceName: serviceName,
		userID:      userID,
	}

	if err := s.validateSubSum(ctx, &subSum); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(ctx, err.Error(), zap.Error(err))

		return
	}

	resp, err := s.service.sum(ctx, subSum)
	if err != nil {
		httputil.HandleErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		return
	}

	log.Info(ctx, "total price calculated")
}

func (s *service) sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error) {
	return s.repo.sum(ctx, subSum)
}

func (r *pgRepository) sum(ctx context.Context, subSum SubSumReq) (SubSumResp, error) {
	log := logger.FromContext(ctx)

	qb := r.builder.Select(`
        COALESCE(SUM(
            (EXTRACT(YEAR FROM age(end_date, start_date)) * 12 + 
             EXTRACT(MONTH FROM age(end_date, start_date))) * price
        ), 0)
    `).
		From("subscriptions").
		Where(squirrel.GtOrEq{"start_date::date": subSum.startDate}).
		Where(squirrel.LtOrEq{"end_date::date": subSum.endDate})

	if subSum.serviceName != "" {
		qb = qb.Where(squirrel.Eq{"service_name": subSum.serviceName})
	}
	if subSum.userID != "" {
		qb = qb.Where(squirrel.Eq{"user_id": subSum.userID})
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return SubSumResp{}, err
	}

	row := r.db.QueryRow(ctx, sql, args...)

	var sum SubSumResp
	if err := row.Scan(&sum.TotalPrice); err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return SubSumResp{}, err
	}

	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return sum, nil
}
