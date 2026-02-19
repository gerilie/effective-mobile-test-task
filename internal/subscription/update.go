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

// @Summary		Update subscription
// @Description	Update by subscription ID.
// @Tags			subscription
// @ID				update-subscription
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Subscription ID"
// @Param			sub	body		SubReq	true	"Subscription"
// @Success		200	{object}	SubResp	"Price must be \u003e 0.\nUser_id must be uuid."
// @Failure		400	{string}	string	"Bad request"
// @Failure		404	{string}	string	"Not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/subscriptions/{id} [put].
func (s *server) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	var sub SubReq
	if err := httputil.DecodeJSON(ctx, w, r, &sub); err != nil {
		return
	}

	if err := s.validateSub(ctx, &sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := s.service.update(ctx, id, sub)
	if err != nil {
		httputil.HandleErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		return
	}

	log.Info(ctx, "subscription updated", zap.String("id", id))
}

func (s *service) update(ctx context.Context, id string, sub SubReq) (SubResp, error) {
	return s.repo.update(ctx, id, sub)
}

func (r *pgRepository) update(ctx context.Context, id string, sub SubReq) (SubResp, error) {
	log := logger.FromContext(ctx)

	sql, args, err := r.builder.Update("subscriptions").
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("user_id", sub.UserID).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Error(ctx, postgres.ErrBuildingQuery.Error(), zap.Error(err))

		return SubResp{}, err
	}

	row := r.db.QueryRow(ctx, sql, args...)

	resp := SubResp{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
	if err := row.Scan(&resp.ID); err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return SubResp{}, err
	}

	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return resp, nil
}
