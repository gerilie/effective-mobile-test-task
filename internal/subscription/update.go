package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// @Summary		Update subscription
// @Description	Update by subscription ID.
// @Tags			subscription
// @ID				update-subscription
// @Accept			json
// @Produce		json
// @Param			id	path		string			true	"Subscription ID"
// @Param			sub	body		UpdateSubReq	true	"Subscription"
// @Success		200	{object}	SubResp			"User ID must be uuid"
// @Failure		400	{object}	validation.Resp	"Bad request"
// @Failure		404	{string}	string			"Not found"
// @Failure		500	{string}	string			"Internal server error"
// @Router			/subscriptions/{id} [patch].
func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	var req UpdateSubReq
	if err := httputil.DecodeJSON(ctx, w, r, &req); err != nil {
		log.Error(ctx, "decode json", zap.Error(err))

		return
	}

	if err := s.validate.StructCtx(ctx, req); err != nil {
		log.Error(ctx, "validate update subscription", zap.Error(err))
		handleValidationErrors(ctx, w, err)

		return
	}

	resp, err := s.service.update(ctx, id, req)
	if err != nil {
		handleServiceErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		log.Error(ctx, "write json", zap.Error(err))

		return
	}

	log.Info(ctx, "subscription updated", zap.String("id", id))
}

func (s *service) update(ctx context.Context, id string, dto UpdateSubReq) (SubResp, error) {
	model, err := updateSubToModel(ctx, dto)
	if err != nil {
		return SubResp{}, fmt.Errorf("subscription update to model: %w", err)
	}

	isDatesSet := model.startDate != nil && model.endDate != nil
	if isDatesSet && model.startDate.After(*model.endDate) {
		return SubResp{}, fmt.Errorf("%w: start date must be before end date", errDateOrder)
	}

	sub, err := s.repo.update(ctx, id, model)
	if err != nil {
		return SubResp{}, fmt.Errorf("update subscription: %w", err)
	}

	return subToDTO(ctx, sub)
}

func (r *pgRepository) update(ctx context.Context, id string, model updateSub) (sub, error) {
	log := logger.FromContext(ctx)

	qb := r.builder.Update("subscriptions").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, service_name, price, user_id, start_date, end_date")

	if model.serviceName != nil {
		qb = qb.Set("service_name", model.serviceName)
	}
	if model.price != nil {
		qb = qb.Set("price", model.price)
	}
	if model.userID != nil {
		qb = qb.Set("user_id", model.userID)
	}
	if model.startDate != nil {
		qb = qb.Set("start_date", model.startDate)
	}
	if model.endDate != nil {
		qb = qb.Set("end_date", model.endDate)
	}

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return sub{}, fmt.Errorf("build query: %w", err)
	}

	row := r.db.QueryRow(ctx, sqlStr, args...)

	var startDate, endDate sql.NullTime
	var sub sub
	if err := row.Scan(
		&sub.id,
		&sub.serviceName,
		&sub.price,
		&sub.userID,
		&startDate,
		&endDate); err != nil {
		return sub, fmt.Errorf("read row: %w", err)
	}

	if startDate.Valid {
		sub.startDate = startDate.Time
	}
	if endDate.Valid {
		sub.endDate = &endDate.Time
	}

	log.Info(ctx, "query executed", zap.String("query", sqlStr), zap.Any("args", args))

	return sub, nil
}
