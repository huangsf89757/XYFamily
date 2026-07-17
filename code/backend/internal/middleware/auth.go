package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"xyfamily/pkg/errors"
	"xyfamily/pkg/jwt"
	"xyfamily/pkg/response"
	"xyfamily/internal/repository"
)

type AuthMiddleware struct {
	jwtMgr *jwt.Manager
	cache  *repository.CacheRepository
}

func NewAuthMiddleware(jwtMgr *jwt.Manager, cache *repository.CacheRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtMgr: jwtMgr, cache: cache}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Fail(c, 401, int(errors.ErrTokenMissing), "missing token")
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 401, int(errors.ErrTokenMissing), "invalid auth header")
			c.Abort()
			return
		}
		tokenStr := parts[1]
		claims, err := m.jwtMgr.ParseAccessToken(tokenStr)
		if err != nil {
			response.Fail(c, 401, int(errors.ErrTokenInvalid), "invalid or expired token")
			c.Abort()
			return
		}
		blacklisted, err := m.cache.IsBlacklisted(c.Request.Context(), claims.JTI)
		if err != nil {
			response.Fail(c, 500, int(errors.ErrInternal), "internal error")
			c.Abort()
			return
		}
		if blacklisted {
			response.Fail(c, 401, int(errors.ErrTokenRevoked), "token has been revoked")
			c.Abort()
			return
		}
		c.Set("account_id", claims.AccountID)
		c.Set("jti", claims.JTI)
		c.Set("org_ids", claims.OrgIDs)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
