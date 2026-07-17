package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"xyfamily/pkg/errors"
	"xyfamily/pkg/response"
	"xyfamily/internal/repository"
)

type RateLimitMiddleware struct {
	cache *repository.CacheRepository
	threshold int
	window    int
	lockout   int
}

func NewRateLimitMiddleware(cache *repository.CacheRepository, threshold, window, lockout int) *RateLimitMiddleware {
	return &RateLimitMiddleware{cache: cache, threshold: threshold, window: window, lockout: lockout}
}

func (m *RateLimitMiddleware) LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lockoutKey := fmt.Sprintf("lockout:%s", ip)
		locked, err := m.cache.IsRateLimited(c.Request.Context(), lockoutKey)
		if err == nil && locked {
			response.Fail(c, 429, int(errors.ErrLoginLocked), "too many failed attempts, please try again later")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *RateLimitMiddleware) RecordLoginFailure(c *gin.Context) {
	ip := c.ClientIP()
	failKey := fmt.Sprintf("login_fail:%s", ip)
	ok, err := m.cache.SetRateLimit(c.Request.Context(), failKey, m.threshold, m.window)
	if err == nil && !ok {
		lockoutKey := fmt.Sprintf("lockout:%s", ip)
		_ = m.cache.SetLockout(c.Request.Context(), lockoutKey, m.lockout)
	}
}
