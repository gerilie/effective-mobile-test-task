package subscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// @Summary		Get subscription summa
// @Description	Get total_price of all subscriptions.
// @Tags			subscription
// @ID				get-subscription-sum
// @Produce		json
// @Param			start_date		query		string			true	"Date format: MM-YYYY"
// @Param			end_date		query		string			true	"Date format: MM-YYYY"
// @Param			service_name	query		string			false	"filter by service name"
// @Param			user_id			query		string			false	"filter by user ID"
// @Success		200				{object}	SubSumResp		"Subscription sum"
// @Failure		400				{object}	validation.Resp	"Bad request"
// @Failure		500				{string}	string			"Internal server error"
// @Router			/subscriptions/sum [get].
func (s *server) sum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	query := r.URL.Query()

	req := SubSumReq{
		StartDate:   query.Get("start_date"),
		EndDate:     query.Get("end_date"),
		ServiceName: query.Get("service_name"),
		UserID:      query.Get("user_id"),
	}

	if err := s.validate.StructCtx(ctx, req); err != nil {
		log.Error(ctx, "validate subscription summa", zap.Error(err))
		handleValidationErrors(ctx, w, err)

		return
	}

	resp, err := s.service.sum(ctx, req)
	if err != nil {
		handleServiceErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		log.Error(ctx, "write json", zap.Error(err))

		return
	}

	log.Info(
		ctx,
		"subscription summa retrieved",
		zap.String("start_date", req.StartDate),
		zap.String("end_date", req.EndDate),
	)
}

func (s *service) sum(ctx context.Context, dto SubSumReq) (SubSumResp, error) {
	model, err := subSumToModel(ctx, dto)
	if err != nil {
		return SubSumResp{}, fmt.Errorf("subscription summa to model: %w", err)
	}

	if model.startDate.After(model.endDate) {
		return SubSumResp{}, fmt.Errorf("%w: start date must be before end date", errDateOrder)
	}

	model, err = s.repo.sum(ctx, model)
	if err != nil {
		return SubSumResp{}, fmt.Errorf("sum subscription: %w", err)
	}

	return subSumToDTO(ctx, model), nil
}

func (r *pgRepository) sum(ctx context.Context, model subSum) (subSum, error) {
	log := logger.FromContext(ctx)

	qb := r.builder.Select().
		Column(squirrel.Expr(`
        COALESCE(SUM(
            (
                EXTRACT(YEAR FROM age(
                    LEAST(end_date, ?),
                    GREATEST(start_date, ?)
                )) * 12 +
                EXTRACT(MONTH FROM age(
                    LEAST(end_date, ?),
                    GREATEST(start_date, ?)
                ))
            ) * price
        ), 0)
    `,
			model.endDate, model.startDate,
			model.endDate, model.startDate)).
		From("subscriptions").
		Where(squirrel.LtOrEq{"start_date": model.endDate}).
		Where(squirrel.GtOrEq{"end_date": model.startDate})

	if model.serviceName != "" {
		qb = qb.Where(squirrel.Eq{"service_name": model.serviceName})
	}
	if model.userID != "" {
		qb = qb.Where(squirrel.Eq{"user_id": model.userID})
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return subSum{}, fmt.Errorf("build query: %w", err)
	}

	row := r.db.QueryRow(ctx, sql, args...)

	var sum subSum
	if err := row.Scan(&sum.totalPrice); err != nil {
		return subSum{}, fmt.Errorf("read row: %w", err)
	}

	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return sum, nil
}
