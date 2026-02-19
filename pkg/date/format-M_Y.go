package date

import (
	"context"
	"errors"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

var ErrInvalidFormat = errors.New("invalid date format")

func FormatDateToPGDate(ctx context.Context, d string) (string, error) {
	log := logger.FromContext(ctx)

	date, err := time.Parse("01-2006", d)
	if err != nil {
		log.Error(ctx, "date parsing", zap.Error(err))

		return "", ErrInvalidFormat
	}

	d = date.Format("2006-01-02")
	log.Info(ctx, "date formatted", zap.String("date", d))

	return d, nil
}
