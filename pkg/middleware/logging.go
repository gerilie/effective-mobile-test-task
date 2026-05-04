package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/yushafro/effective-mobile-tz/pkg/deferfunc"
	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

// Logging is a middleware function that logs HTTP requests and responses.
// It takes a logger.Logger and an http.Handler as arguments.
// It returns an http.Handler.
//
// The Logging middleware adds the request ID to the context and logs the request and response.
// It also handles errors and stops the logger after the request is processed.
func Logging(next http.Handler, log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.WithLogger(r.Context(), log)
		defer deferfunc.Close(ctx, log.Stop, "error stopping logger")

		id := r.Header.Get(httputil.RequestID)
		if id == "" {
			log.Info(ctx, "empty request id, creating new one")
			id = uuid.NewString()
		}

		ctx = logger.WithRequestID(ctx, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
