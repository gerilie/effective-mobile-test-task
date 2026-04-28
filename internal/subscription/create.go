package subscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// @Summary	Create subscription
// @Tags		subscription
// @ID			create-subscription
// @Accept		json
// @Produce	json
// @Param		sub	body		SubReq			true	"User ID must be uuid\nDate format: MM-YYYY"
// @Success	201	{object}	SubResp			"Created subscription"
// @Failure	400	{object}	validation.Resp	"Bad request"
// @Failure	404	{string}	string			"Not found"
// @Failure	500	{string}	string			"Internal server error"
// @Router		/subscriptions [post].
func (s *server) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req SubReq
	if err := httputil.DecodeJSON(ctx, w, r, &req); err != nil {
		log.Error(ctx, "decode json", zap.Error(err))

		return
	}

	if err := s.validate.StructCtx(ctx, req); err != nil {
		log.Error(ctx, "validate subscription", zap.Error(err))
		handleValidationErrors(ctx, w, err)

		return
	}

	resp, err := s.service.create(ctx, req)
	if err != nil {
		handleServiceErrors(ctx, w, err)

		return
	}

	w.Header().Set("Location", fmt.Sprint("/subscriptions/", resp.ID))
	if err := httputil.WriteJSON(ctx, w, http.StatusCreated, resp); err != nil {
		log.Error(ctx, "write response", zap.Error(err))

		return
	}

	log.Info(ctx, "subscription created")
}

func (s *service) create(ctx context.Context, dto SubReq) (SubResp, error) {
	model, err := subToModel(ctx, dto)
	if err != nil {
		return SubResp{}, fmt.Errorf("subscription to model: %w", err)
	}

	if model.endDate != nil && model.startDate.After(*model.endDate) {
		return SubResp{}, fmt.Errorf("%w: start date must be before end date", errDateOrder)
	}

	model, err = s.repo.create(ctx, model)
	if err != nil {
		return SubResp{}, fmt.Errorf("create subscription: %w", err)
	}

	return subToDTO(ctx, model)
}

func (r *pgRepository) create(ctx context.Context, model sub) (sub, error) {
	log := logger.FromContext(ctx)

	sql, args, err := r.builder.Insert("subscriptions").
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(model.serviceName, model.price, model.userID, model.startDate, model.endDate).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return sub{}, fmt.Errorf("build query: %w", err)
	}

	row := r.db.QueryRow(ctx, sql, args...)
	if err := row.Scan(&model.id); err != nil {
		return sub{}, fmt.Errorf("read row: %w", err)
	}
	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return model, nil
}
