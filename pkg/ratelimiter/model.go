package ratelimiter

import (
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter defines a basic interface for rate limiter implementations
// that support cleanup of internal state.
type RateLimiter interface {
	CleanUp()
}

// IPRateLimiter extends RateLimiter with IP-based rate limiting capabilities.
type IPRateLimiter interface {
	RateLimiter

	// GetLimiter returns a rate limiter associated with the given IP address.
	// If the IP is empty, a new independent limiter is created.
	GetLimiter(ip string) *rate.Limiter

	// GetIPCount returns the number of tracked IPs in the limiter.
	GetIPCount() int
}

// rateLimiter contains shared configuration for rate limiting.
type rateLimiter struct {
	mu *sync.RWMutex
	r  rate.Limit
	b  int
}

// ipRateLimiter implements IPRateLimiter and stores per-IP limiters.
type ipRateLimiter struct {
	rateLimiter

	ips map[string]*rate.Limiter
}

// NewIPRateLimiter creates a new IP-based rate limiter.
//
// If r or b are less than or equal to zero, they are replaced with safe defaults.
func NewIPRateLimiter(r rate.Limit, b int) IPRateLimiter {
	if r <= 0 {
		r = rate.Limit(1)
	}
	if b <= 0 {
		b = 1
	}

	return &ipRateLimiter{
		rateLimiter: rateLimiter{
			mu: &sync.RWMutex{},
			r:  r,
			b:  b,
		},
		ips: make(map[string]*rate.Limiter),
	}
}

// GetLimiter returns a rate limiter for the given IP address.
// If the limiter does not exist, it is created and stored.
//
// If ip is empty, a new standalone limiter is returned.
func (l *ipRateLimiter) GetLimiter(ip string) *rate.Limiter {
	if ip == "" {
		return rate.NewLimiter(l.r, l.b)
	}

	l.mu.RLock()
	limiter, ok := l.ips[ip]
	l.mu.RUnlock()

	if ok {
		return limiter
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok = l.ips[ip]
	if !ok {
		limiter = rate.NewLimiter(l.r, l.b)
		l.ips[ip] = limiter
	}

	return limiter
}

// CleanUp resets all stored IP rate limiters.
func (l *ipRateLimiter) CleanUp() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ips = make(map[string]*rate.Limiter)
}

// GetIPCount returns the number of IP addresses currently tracked by the limiter.
func (l *ipRateLimiter) GetIPCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.ips)
}
