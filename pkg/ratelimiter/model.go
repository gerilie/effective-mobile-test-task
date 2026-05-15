package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter provides per-IP rate limiting functionality.
//
// Implementations are responsible for storing and managing
// independent rate limiters for different client IP addresses.
type IPRateLimiter interface {
	// GetLimiter returns a rate limiter associated with the given IP address.
	//
	// If a limiter for the provided IP does not exist, it is created
	// and stored internally.
	//
	// If ip is empty, a new standalone limiter is returned without storing it.
	GetLimiter(ip string) *rate.Limiter

	// StartCleanUp starts a background cleanup routine that periodically removes
	// inactive IP limiters from internal storage using the configured cleanup settings.
	StartCleanUp()
}

// client stores a rate limiter and metadata associated with a single IP address.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Config defines the configuration for an IP rate limiter.
type Config struct {
	// R defines the number of allowed events per second.
	R rate.Limit

	// B defines the maximum burst size.
	B int

	// CleanUpInterval defines how often inactive IP limiters are checked and removed.
	CleanUpInterval time.Duration

	// CleanUpMaxIdle defines how long an IP limiter may remain inactive
	// before it is removed from storage.
	CleanUpMaxIdle time.Duration
}

// ipRateLimiter implements IPRateLimiter using an in-memory map of IP addresses.
type ipRateLimiter struct {
	Config

	mu  *sync.Mutex
	ips map[string]*client
}

// NewIPRateLimiter creates a new in-memory IP rate limiter using the provided configuration.
//
// If cfg.r is less than or equal to zero, a default rate of 1 request per second is used.
//
// If cfg.b is less than or equal to zero, a default burst size of 1 is used.
func NewIPRateLimiter(cfg Config) IPRateLimiter {
	if cfg.R <= 0 {
		cfg.R = rate.Limit(1)
	}

	if cfg.B <= 0 {
		cfg.B = 1
	}

	return &ipRateLimiter{
		Config: Config{
			R: cfg.R,
			B: cfg.B,

			CleanUpInterval: cfg.CleanUpInterval,
			CleanUpMaxIdle:  cfg.CleanUpMaxIdle,
		},

		mu:  &sync.Mutex{},
		ips: make(map[string]*client),
	}
}
