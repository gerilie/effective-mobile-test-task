package ratelimiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	CleanUp()
}

type IPRateLimiter interface {
	RateLimiter

	GetLimiter(ip string) *rate.Limiter
	GetIPCount() int
}

type rateLimiter struct {
	mu *sync.RWMutex
	r  rate.Limit
	b  int
}

type ipRateLimiter struct {
	rateLimiter

	ips map[string]*rate.Limiter
}

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

func (l *ipRateLimiter) CleanUp() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ips = make(map[string]*rate.Limiter)
}

func (l *ipRateLimiter) GetIPCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.ips)
}
