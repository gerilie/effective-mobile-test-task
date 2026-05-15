package ratelimiter

import "time"

// StartCleanUp starts a background cleanup routine that periodically removes
// inactive IP limiters from internal storage.
//
// Cleanup behavior is configured through the limiter Config:
//
//   - cleanUpInterval defines how often cleanup runs.
//   - cleanUpMaxIdle defines how long a limiter may remain inactive
//     before removal.
func (l *ipRateLimiter) StartCleanUp() {
	ticker := time.NewTicker(l.CleanUpInterval)

	go func() {
		for range ticker.C {
			l.cleanUp(l.CleanUpMaxIdle)
		}
	}()
}

// cleanUp removes IP limiters that have been inactive longer than maxIdle.
func (l *ipRateLimiter) cleanUp(maxIdle time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for ip, c := range l.ips {
		if time.Since(c.lastSeen) > maxIdle {
			delete(l.ips, ip)
		}
	}
}
