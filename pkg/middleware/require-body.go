package middleware

import (
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/logger"
)

func RequireBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if r.ContentLength == 0 {
			log.Error(ctx, "Empty request body error")
			http.Error(w, "request body is empty", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
