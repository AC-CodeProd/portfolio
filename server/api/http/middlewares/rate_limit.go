package middlewares

import (
	"net/http"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/logger"
	"portfolio/shared"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors     map[string]*visitor
	mu           sync.RWMutex
	rate         time.Duration
	limit        int
	cleanupTimer *time.Timer
	logger       *logger.Logger
}

type visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

func NewRateLimiter(rate time.Duration, limit int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		limit:    limit,
		logger:   nil,
	}

	return rl
}

func (rl *RateLimiter) SetLogger(logger *logger.Logger) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.logger = logger
}

func (rl *RateLimiter) scheduleCleanup() {
	if rl.cleanupTimer != nil {
		rl.cleanupTimer.Stop()
	}
	rl.cleanupTimer = time.AfterFunc(5*time.Minute, func() {
		rl.cleanup()
	})
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	before := len(rl.visitors)
	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			close(v.limiter)
			delete(rl.visitors, ip)
		}
	}

	if len(rl.visitors) > 0 {
		rl.scheduleCleanup()
	}

	if rl.logger != nil && before > len(rl.visitors) {
		rl.logger.Debug("Rate limiter cleaned %d inactive visitors", before-len(rl.visitors))
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mu.RLock()
		v, exists := rl.visitors[ip]
		rl.mu.RUnlock()

		if !exists {
			rl.mu.Lock()
			if v, exists = rl.visitors[ip]; !exists {
				rl.visitors[ip] = &visitor{
					limiter:  make(chan struct{}, rl.limit),
					lastSeen: time.Now(),
				}
				v = rl.visitors[ip]

				if len(rl.visitors) == 1 {
					rl.scheduleCleanup()
				}
			}
			rl.mu.Unlock()
		}

		select {
		case v.limiter <- struct{}{}:
			next.ServeHTTP(w, r)

			time.AfterFunc(rl.rate, func() {
				select {
				case <-v.limiter:
				default:
					// Channel already closed (visitor removed)
				}
			})

			rl.mu.Lock()
			if v, exists := rl.visitors[ip]; exists {
				v.lastSeen = time.Now()
			}
			rl.mu.Unlock()

		default:
			domainErr := domain.NewRateLimitError("Too many requests")
			apiError := utils.DomainErrorToAPIError(domainErr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(domainErr.HTTPStatus())

			response := struct {
				Errors []*shared.APIError `json:"errors"`
			}{
				Errors: []*shared.APIError{apiError},
			}

			utils.JSONResponse(w, domainErr.HTTPStatus(), response)
		}
	})
}
