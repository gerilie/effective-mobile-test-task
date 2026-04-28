package middleware

import (
	"fmt"
	"math"
	"net/http"

	"github.com/yushafro/effective-mobile-tz/pkg/httputil"
	"github.com/yushafro/effective-mobile-tz/pkg/logger"
	"github.com/yushafro/effective-mobile-tz/pkg/ratelimiter"
)

func RateLimiter(next http.Handler, l ratelimiter.IPRateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		ip := httputil.GetClientIP(r)
		if ip == "" {
			log.Error(ctx, "get client ip")
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		limiter := l.GetLimiter(ip)
		if limiter == nil {
			log.Error(ctx, "get limiter")
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		limit := fmt.Sprintf("%d", limiter.Burst())
		tokens := limiter.Tokens()
		remaining := fmt.Sprintf("%.0f", math.Floor(tokens))

		w.Header().Set("X-RateLimit-Limit", limit)
		w.Header().Set("X-RateLimit-Remaining", remaining)

		if !limiter.Allow() {
			log.Error(ctx, "too many requests")
			http.Error(w, "too many requests", http.StatusTooManyRequests)

			return
		}

		next.ServeHTTP(w, r)
	})
}
