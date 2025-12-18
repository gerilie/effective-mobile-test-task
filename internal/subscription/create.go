package subscription

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

// @Summary		Create subscription
// @Description	ID will be generated.
// @Description	Price must be > 0.
// @Description	User_id must be uuid.
// @Description	Date format: 01-2006. End_date is optional.
// @Tags			subscription
// @ID				create-subscription
// @Accept			json
// @Produce		json
// @Param			sub	body		SubReq	true	"Subscription"
// @Success		201	{object}	SubResp	"Created subscription"
// @Failure		400	{string}	string	"Bad request"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/subscriptions [post].
func (s *server) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var sub SubReq
	if err := httputil.DecodeJSON(ctx, r, &sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := s.validateSub(ctx, &sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	resp, err := s.service.create(ctx, sub)
	if err != nil {
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

	w.Header().Set("Location", fmt.Sprint("/subscriptions/", resp.ID))
	if err := httputil.WriteJSON(ctx, w, http.StatusCreated, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Info(ctx, "subscription created")
}

func (s *service) create(ctx context.Context, sub SubReq) (SubResp, error) {
	return s.repo.create(ctx, sub)
}

func (r *pgRepository) create(ctx context.Context, sub SubReq) (SubResp, error) {
	log := logger.FromContext(ctx)

	sql, args, err := r.builder.Insert("subscriptions").
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
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

		return resp, err
	}
	log.Info(ctx, "query executed", zap.String("query", sql), zap.Any("args", args))

	return resp, nil
}
