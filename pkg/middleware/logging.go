package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

func Logging(next http.Handler, logCfg logger.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.NewWithConfig(logCfg)
		ctx = logger.WithLogger(ctx, log)
		defer deferfunc.Close(ctx, log.Stop, "error stopping logger")

		id := r.Header.Get(httputil.RequestID)
		if id == "" {
			log.Info(ctx, "empty request id, creating new one")
			id = uuid.NewString()
		}

		ctx = context.WithValue(ctx, logger.RequestID, id)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
