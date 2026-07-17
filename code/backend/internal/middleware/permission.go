package middleware

import (
	"github.com/gin-gonic/gin"

	"xyfamily/internal/repository"
	appErr "xyfamily/pkg/errors"
	"xyfamily/pkg/response"
)

type PermissionChecker struct {
	rbacRepo *repository.RBACRepository
	cache    *repository.CacheRepository
}

func NewPermissionChecker(rbacRepo *repository.RBACRepository, cache *repository.CacheRepository) *PermissionChecker {
	return &PermissionChecker{rbacRepo: rbacRepo, cache: cache}
}

func (p *PermissionChecker) RequirePermission(permKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("effective_role")
		roleStr, _ := role.(string)
		if roleStr == "super_admin" {
			c.Next()
			return
		}
		perms, err := p.rbacRepo.GetPermissionsByRoleKey(c.Request.Context(), roleStr)
		if err != nil {
			response.Fail(c, 500, int(appErr.ErrInternal), "failed to check permissions")
			c.Abort()
			return
		}
		for _, p := range perms {
			if p == permKey {
				c.Next()
				return
			}
		}
		response.Fail(c, 403, 600002, "permission denied: "+permKey)
		c.Abort()
	}
}
