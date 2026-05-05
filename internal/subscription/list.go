package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

// @Summary		List subscriptions
// @Description	Get paginated list of subscriptions with optional filters
// @Tags			subscription
// @ID				subscription-list
// @Produce		json
// @Param			page			query		int			true	"Page number (1-based)"
// @Param			limit			query		int			true	"Items per page (max: 100)"
// @Param			service_name	query		string		false	"filter by service name"
// @Param			user_id			query		string		false	"filter by user ID"
// @Success		200				{object}	SubListResp	"List subscriptions"
// @Failure		400				{string}	string		"Bad request"
// @Failure		500				{string}	string		"Internal server error"
// @Router			/subscriptions [get].
func (s *server) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	query := r.URL.Query()

	serviceName := query.Get("service_name")
	userID := query.Get("user_id")

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		log.Error(ctx, "bad request", zap.Error(err))
		http.Error(w, "bad request: invalid page", http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		log.Error(ctx, "bad request", zap.Error(err))
		http.Error(w, "bad request: invalid limit", http.StatusBadRequest)

		return
	}

	req := SubListReq{
		ServiceName: serviceName,
		UserID:      userID,
		Page:        page,
		Limit:       limit,
	}

	if err := s.validate.StructCtx(ctx, req); err != nil {
		log.Error(ctx, "validate subscription list", zap.Error(err))
		handleValidationErrors(ctx, w, err)

		return
	}

	resp, err := s.service.list(ctx, req)
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
		"subscriptions listed",
		zap.Int("page", page),
		zap.Int("limit", limit),
	)
}

func (s *service) list(ctx context.Context, dto SubListReq) (SubListResp, error) {
	model := subListToModel(ctx, dto)

	subs, err := s.repo.list(ctx, model)
	if err != nil {
		return SubListResp{}, fmt.Errorf("list subscriptions: %w", err)
	}

	return subListToDTO(ctx, subs)
}

func (r *pgRepository) list(ctx context.Context, model subList) ([]sub, error) {
	log := logger.FromContext(ctx)

	qb := r.builder.Select("id, service_name, price, user_id, start_date, end_date").
		From("subscriptions").
		Limit(uint64(model.limit)).
		Offset(uint64(model.offset)).
		OrderBy("id")

	if model.serviceName != "" {
		qb = qb.Where(squirrel.Eq{"service_name": model.serviceName})
	}
	if model.userID != "" {
		qb = qb.Where(squirrel.Eq{"user_id": model.userID})
	}

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return []sub{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return []sub{}, fmt.Errorf("read rows: %w", err)
	}
	defer rows.Close()

	var subs []sub
	for rows.Next() {
		var startDate, endDate sql.NullTime
		var sub sub
		if err := rows.Scan(
			&sub.id,
			&sub.serviceName,
			&sub.price,
			&sub.userID,
			&startDate,
			&endDate,
		); err != nil {
			return subs, fmt.Errorf("read row: %w", err)
		}

		if startDate.Valid {
			sub.startDate = startDate.Time
		}
		if endDate.Valid {
			sub.endDate = &endDate.Time
		}

		subs = append(subs, sub)
	}

	log.Info(ctx, "query executed", zap.String("query", sqlStr), zap.Any("args", args))

	return subs, nil
}
