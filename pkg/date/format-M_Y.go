package date

import (
	"context"
	"time"

	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func FormatM_Y(ctx context.Context, d *string) error {
	log := logger.FromContext(ctx)

	if d == nil {
		return nil
	}

	date, err := time.Parse("01-2006", *d)
	if err != nil {
		log.Error(ctx, "date parsing", zap.Error(err))

		return err
	}

	*d = date.Format("2006-01-02")
	log.Info(ctx, "date formatted", zap.String("date", *d))

	return nil
}
