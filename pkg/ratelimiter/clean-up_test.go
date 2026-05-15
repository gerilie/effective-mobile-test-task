package ratelimiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/yushafro/effective-mobile-tz/pkg/ratelimiter"
	"golang.org/x/time/rate"
)

func Test_ipRateLimiter_StartCleanUp_CleanUp(t *testing.T) {
	t.Parallel()

	config := ratelimiter.Config{
		R:               rate.Limit(1),
		B:               1,
		CleanUpInterval: 10 * time.Millisecond,
		CleanUpMaxIdle:  10 * time.Millisecond,
	}

	t.Run("clean up inactive limiters", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)
		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")

		ipRateLimiter.StartCleanUp(context.Background())
		time.Sleep(15 * time.Millisecond)

		limiter2 := ipRateLimiter.GetLimiter("1.1.1.1")
		require.NotSame(t, limiter1, limiter2)
	})

	t.Run("not clean up active limiters", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)
		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")

		ipRateLimiter.StartCleanUp(context.Background())

		for range 5 {
			time.Sleep(5 * time.Millisecond)

			limiter2 := ipRateLimiter.GetLimiter("1.1.1.1")
			require.Same(t, limiter1, limiter2)
		}
	})
}

func Test_ipRateLimiter_StartCleanUp_Sync(t *testing.T) {
	t.Parallel()

	config := ratelimiter.Config{
		R:               rate.Limit(1),
		B:               1,
		CleanUpInterval: 10 * time.Millisecond,
		CleanUpMaxIdle:  10 * time.Millisecond,
	}

	t.Run("stop via context", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		ipRateLimiter.StartCleanUp(ctx)
		cancel()

		time.Sleep(15 * time.Millisecond)
		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")
		time.Sleep(15 * time.Millisecond)
		limiter2 := ipRateLimiter.GetLimiter("1.1.1.1")

		require.Same(t, limiter1, limiter2)
	})

	t.Run("StartCleanUp starts only once", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		ipRateLimiter.StartCleanUp(ctx)
		cancel()

		time.Sleep(5 * time.Millisecond)

		ipRateLimiter.StartCleanUp(ctx)
		limiter1 := ipRateLimiter.GetLimiter("1.1.1.1")

		time.Sleep(15 * time.Millisecond)
		limiter2 := ipRateLimiter.GetLimiter("1.1.1.1")

		require.Same(t, limiter1, limiter2)
	})

	t.Run("concurrent access", func(t *testing.T) {
		t.Parallel()

		ipRateLimiter := ratelimiter.NewIPRateLimiter(config)

		var limiter1, limiter2 *rate.Limiter
		var wg sync.WaitGroup
		wg.Go(func() {
			limiter1 = ipRateLimiter.GetLimiter("1.1.1.1")
		})
		wg.Go(func() {
			limiter2 = ipRateLimiter.GetLimiter("1.1.1.1")
		})
		wg.Wait()

		require.Same(t, limiter1, limiter2)
	})
}
