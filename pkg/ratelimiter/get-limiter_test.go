package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/yushafro/effective-mobile-tz/pkg/ratelimiter"
	"golang.org/x/time/rate"
)

func Test_ipRateLimiter_GetLimiter_Creating(t *testing.T) {
	t.Parallel()

	config := ratelimiter.Config{
		R:               rate.Limit(1),
		B:               1,
		CleanUpInterval: time.Minute,
		CleanUpMaxIdle:  time.Minute,
	}

	t.Run("the same limiter for the same ip", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")
		limiter2 := ipRateLimiter.GetLimiter("1.1.1.1")

		require.Same(t, limiter1, limiter2)
	})

	t.Run("different limiter for different ip", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")
		limiter2 := ipRateLimiter.GetLimiter("1.1.1.2")

		require.NotSame(t, limiter1, limiter2)
	})

	t.Run("standalone limiter for empty ip", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		limiter1 := ipRateLimiter.GetLimiter("")
		limiter2 := ipRateLimiter.GetLimiter("")

		require.NotSame(t, limiter1, limiter2)
	})
}

func Test_ipRateLimiter_GetLimiter_Rate(t *testing.T) {
	t.Parallel()

	config := ratelimiter.Config{
		R:               rate.Limit(1),
		B:               1,
		CleanUpInterval: time.Minute,
		CleanUpMaxIdle:  time.Minute,
	}

	ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

	limiter := ipRateLimiter.GetLimiter("1.1.1.1")

	require.True(t, limiter.Allow())
	require.False(t, limiter.Allow())
}
