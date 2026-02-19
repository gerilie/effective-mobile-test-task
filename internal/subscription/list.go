package subscription

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/postgres"
	"go.uber.org/zap"
)

// @Summary		List subscriptions
// @Description	Get paginated list of subscriptions with optional filters
// @Tags			subscription
// @ID				list-subscription
// @Produce		json
// @Param			page	query		int			true	"Page number (1-based)"
// @Param			limit	query		int			true	"Items per page (max: 100)"
// @Success		200		{object}	ListSubResp	"List subscriptions"
// @Failure		400		{string}	string		"Bad request"
// @Failure		500		{string}	string		"Internal server error"
// @Router			/subscriptions [get].
func (s *server) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(ctx, "bad request", zap.Error(err))

		return
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(ctx, "bad request", zap.Error(err))

		return
	}

	list, err := s.service.list(ctx, ListSubReq{
		page:  page,
		limit: limit,
	})
	if err != nil {
		httputil.HandleErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, list); err != nil {
		return
	}

	log.Info(
		ctx,
		"list of subscriptions sent",
		zap.Int("page", page),
		zap.Int("limit", limit),
	)
}

func (s *service) list(ctx context.Context, list ListSubReq) (ListSubResp, error) {
	list.offset = (list.page - 1) * list.limit

	if list.limit <= 0 {
		list.limit = 20
	}
	if list.limit > 100 {
		list.limit = 100
	}

	return s.repo.list(ctx, list)
}

func (r *pgRepository) list(ctx context.Context, list ListSubReq) (ListSubResp, error) {
	log := logger.FromContext(ctx)

	sqlStr, args, err := r.builder.Select("id, service_name, price, user_id, start_date, end_date").
		From("subscriptions").
		Limit(uint64(list.limit)).
		Offset(uint64(list.offset)).
		ToSql()
	if err != nil {
		log.Error(ctx, postgres.ErrBuildingQuery.Error(), zap.Error(err))

		return ListSubResp{}, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

		return ListSubResp{}, err
	}
	defer rows.Close()

	var subs []SubResp
	for rows.Next() {
		var startDate, endDate sql.NullTime
		var sub SubResp
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &startDate, &endDate); err != nil {
			log.Error(ctx, postgres.ErrReadingRow.Error(), zap.Error(err))

			return ListSubResp{}, err
		}

		if startDate.Valid {
			sub.StartDate = startDate.Time.Format("01-2006")
		}
		if endDate.Valid {
			endStr := endDate.Time.Format("01-2006")
			sub.EndDate = &endStr
		}

		subs = append(subs, sub)
	}

	log.Info(ctx, "query executed", zap.String("query", sqlStr), zap.Any("args", args))

	return subs, nil
}
