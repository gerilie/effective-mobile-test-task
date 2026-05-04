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

// get handles HTTP request for retrieving a subscription by ID.
//
// It returns subscription data in JSON format.
func (s *server) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	id := r.PathValue("id")

	resp, err := s.service.get(ctx, id)
	if err != nil {
		handleServiceErrors(ctx, w, err)

		return
	}

	if err := httputil.WriteJSON(ctx, w, http.StatusOK, resp); err != nil {
		log.Error(ctx, "write json", zap.Error(err))

		return
	}

	log.Info(ctx, "subscription retrieved", zap.String("id", id))
}

func (s *service) get(ctx context.Context, id string) (SubResp, error) {
	model, err := s.repo.get(ctx, id)
	if err != nil {
		return SubResp{}, fmt.Errorf("get subscription: %w", err)
	}

	return subToDTO(ctx, model)
}

func (r *pgRepository) get(ctx context.Context, id string) (sub, error) {
	log := logger.FromContext(ctx)

	sqlStr, args, err := r.builder.Select("id, service_name, price, user_id, start_date, end_date").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
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
		&endDate,
	); err != nil {
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
