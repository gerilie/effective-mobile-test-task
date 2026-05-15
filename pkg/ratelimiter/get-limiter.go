package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"
)

// GetLimiter returns a rate limiter associated with the provided IP address.
//
// If a limiter for the IP does not already exist, a new limiter is created,
// stored, and returned.
//
// The client's last activity timestamp is updated on every call.
//
// If ip is empty, a new standalone limiter is returned without storing it
// in the internal IP map.
func (l *ipRateLimiter) GetLimiter(ip string) *rate.Limiter {
	if ip == "" {
		return rate.NewLimiter(l.R, l.B)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	c := l.ips[ip]
	now := time.Now()

	if c == nil {
		c = &client{
			limiter:  rate.NewLimiter(l.R, l.B),
			lastSeen: now,
		}

		l.ips[ip] = c
	} else {
		c.lastSeen = now
	}

	return c.limiter
}
