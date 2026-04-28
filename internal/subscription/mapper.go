package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

func subToModel(ctx context.Context, dto SubReq) (sub, error) {
	log := logger.FromContext(ctx)

	startDate, err := time.Parse(dtoDateLayout, dto.StartDate)
	if err != nil {
		return sub{}, fmt.Errorf("invalid format %s: %w", dto.StartDate, err)
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		parsed, err := time.Parse(dtoDateLayout, *dto.EndDate)
		if err != nil {
			return sub{}, fmt.Errorf("invalid format %s: %w", *dto.EndDate, err)
		}

		endDate = &parsed
	}

	log.Info(ctx, "subscription mapped to domain")

	return sub{
		serviceName: dto.ServiceName,
		price:       dto.Price,
		userID:      dto.UserID,
		startDate:   startDate,
		endDate:     endDate,
	}, nil
}

func subToDTO(ctx context.Context, sub sub) (SubResp, error) {
	log := logger.FromContext(ctx)

	var endDate *string
	if sub.endDate != nil {
		parsed := sub.endDate.Format(dtoDateLayout)
		endDate = &parsed
	}

	log.Info(ctx, "subscription mapped to dto")

	return SubResp{
		ID:          sub.id,
		ServiceName: sub.serviceName,
		Price:       sub.price,
		UserID:      sub.userID,
		StartDate:   sub.startDate.Format(dtoDateLayout),
		EndDate:     endDate,
	}, nil
}

func subListToModel(ctx context.Context, dto SubListReq) subList {
	log := logger.FromContext(ctx)

	log.Info(ctx, "subscription list mapped to model")

	return subList{
		serviceName: dto.ServiceName,
		userID:      dto.UserID,
		page:        dto.Page,
		limit:       dto.Limit,
		offset:      (dto.Page - 1) * dto.Limit,
	}
}

func subListToDTO(ctx context.Context, subs []sub) (SubListResp, error) {
	log := logger.FromContext(ctx)
	subListResp := make(SubListResp, 0, len(subs))

	for _, sub := range subs {
		subResp, err := subToDTO(ctx, sub)
		if err != nil {
			return nil, fmt.Errorf("subscription to dto: %w", err)
		}

		subListResp = append(subListResp, subResp)
	}

	log.Info(ctx, "subscription list mapped to DTO")

	return subListResp, nil
}

func subSumToModel(ctx context.Context, dto SubSumReq) (subSum, error) {
	log := logger.FromContext(ctx)

	startDate, err := time.Parse(dtoDateLayout, dto.StartDate)
	if err != nil {
		return subSum{}, fmt.Errorf("invalid format %s: %w", dto.StartDate, err)
	}

	endDate, err := time.Parse(dtoDateLayout, dto.EndDate)
	if err != nil {
		return subSum{}, fmt.Errorf("invalid format %s: %w", dto.EndDate, err)
	}

	log.Info(ctx, "subscription summa mapped to model")

	return subSum{
		serviceName: dto.ServiceName,
		userID:      dto.UserID,
		startDate:   startDate,
		endDate:     endDate,
		totalPrice:  0,
	}, nil
}

func subSumToDTO(ctx context.Context, sum subSum) SubSumResp {
	log := logger.FromContext(ctx)

	log.Info(ctx, "subscription summa mapped to dto")

	return SubSumResp{
		TotalPrice: sum.totalPrice,
	}
}

func updateSubToModel(ctx context.Context, dto UpdateSubReq) (updateSub, error) {
	log := logger.FromContext(ctx)

	model := updateSub{
		id:          dto.ID,
		serviceName: dto.ServiceName,
		price:       dto.Price,
		userID:      dto.UserID,
	}

	if dto.StartDate != nil {
		startDate, err := time.Parse(dtoDateLayout, *dto.StartDate)
		if err != nil {
			return updateSub{}, fmt.Errorf("invalid format %T: %w", dto.StartDate, err)
		}

		model.startDate = &startDate
	}

	if dto.EndDate != nil {
		endDate, err := time.Parse(dtoDateLayout, *dto.EndDate)
		if err != nil {
			return updateSub{}, fmt.Errorf("invalid format %T: %w", dto.EndDate, err)
		}

		model.endDate = &endDate
	}

	log.Info(ctx, "updating subscription mapped to model")

	return model, nil
}
