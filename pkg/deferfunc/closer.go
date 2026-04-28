package deferfunc

import (
	"context"
	"errors"
	"syscall"

	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"go.uber.org/zap"
)

func Close(ctx context.Context, c func() error, errMsg string) {
	log := logger.FromContext(ctx)

	if err := c(); err != nil && !errors.Is(err, syscall.EINVAL) {
		log.Error(ctx, errMsg, zap.Error(err))
	}
}
