package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/repository"
	appErr "xyfamily/pkg/errors"
	"xyfamily/pkg/response"
)

type MembershipValidator struct {
	rbacRepo *repository.RBACRepository
	cache    *repository.CacheRepository
}

func NewMembershipValidator(rbacRepo *repository.RBACRepository, cache *repository.CacheRepository) *MembershipValidator {
	return &MembershipValidator{rbacRepo: rbacRepo, cache: cache}
}

func (m *MembershipValidator) ValidateScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIDStr, _ := c.Get("account_id")
		accountID, err := uuid.Parse(accountIDStr.(string))
		if err != nil {
			response.Fail(c, 401, int(appErr.ErrTokenInvalid), "invalid account")
			c.Abort()
			return
		}
		roles, _ := c.Get("roles")
		rolesMap, _ := roles.(map[string]string)
		if _, isSuperAdmin := rolesMap["super"]; isSuperAdmin {
			c.Set("effective_role", "super_admin")
			c.Set("scope_type", "global")
			c.Next()
			return
		}
		orgIDsRaw, _ := c.Get("org_ids")
		orgIDStrs, _ := orgIDsRaw.([]string)
		var orgID uuid.UUID
		orgHeader := c.GetHeader("X-Organization-ID")
		if orgHeader != "" {
			orgID, err = uuid.Parse(orgHeader)
			if err != nil {
				response.Fail(c, 400, int(appErr.ErrBadRequest), "invalid X-Organization-ID")
				c.Abort()
				return
			}
		} else if len(orgIDStrs) == 1 {
			orgID, _ = uuid.Parse(orgIDStrs[0])
		} else if len(orgIDStrs) == 0 {
			response.Fail(c, 400, int(appErr.ErrBadRequest), "no organization context")
			c.Abort()
			return
		} else {
			response.Fail(c, 400, int(appErr.ErrBadRequest), "X-Organization-ID required for multi-org user")
			c.Abort()
			return
		}
		teamHeader := c.GetHeader("X-Team-ID")
		groupHeader := c.GetHeader("X-Group-ID")
		var scopeType, effectiveRole string
		var scopeID uuid.UUID
		if groupHeader != "" {
			scopeType = "group"
			scopeID, _ = uuid.Parse(groupHeader)
			role, err := m.rbacRepo.GetGroupRoleByAccountID(c.Request.Context(), accountID, scopeID)
			if err == nil && role != "" {
				effectiveRole = role
			} else {
				if teamHeader != "" {
					teamID, _ := uuid.Parse(teamHeader)
					role, err = m.rbacRepo.GetTeamRoleByAccountID(c.Request.Context(), accountID, teamID)
					if err == nil && role != "" {
						effectiveRole = m.downgradeRole(role, "team")
					}
			}
				if effectiveRole == "" {
					role, err = m.rbacRepo.GetOrgRoleByAccountID(c.Request.Context(), accountID, orgID)
					if err == nil && role != "" {
						effectiveRole = m.downgradeRole(role, "group")
					}
			}
		}
		} else if teamHeader != "" {
			scopeType = "team"
			scopeID, _ = uuid.Parse(teamHeader)
			role, err := m.rbacRepo.GetTeamRoleByAccountID(c.Request.Context(), accountID, scopeID)
			if err == nil && role != "" {
				effectiveRole = role
			} else {
				role, err = m.rbacRepo.GetOrgRoleByAccountID(c.Request.Context(), accountID, orgID)
				if err == nil && role != "" {
					effectiveRole = m.downgradeRole(role, "team")
				}
			}
		} else {
			scopeType = "organization"
			scopeID = orgID
			role, err := m.rbacRepo.GetOrgRoleByAccountID(c.Request.Context(), accountID, orgID)
			if err == nil && role != "" {
				effectiveRole = role
			}
		}
		if effectiveRole == "" {
			response.Fail(c, 403, 600001, "not a member of this scope")
			c.Abort()
			return
		}
		c.Set("scope_type", scopeType)
		c.Set("scope_id", scopeID)
		c.Set("org_id", orgID)
		c.Set("effective_role", effectiveRole)
		c.Next()
	}
}

func (m *MembershipValidator) downgradeRole(role, targetScope string) string {
	prefix := strings.SplitN(role, "_", 2)
	if len(prefix) < 2 {
		return "regular_member"
	}
	return targetScope + "_ordinary_admin"
}

func (m *MembershipValidator) CheckScope(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		actual, _ := c.Get("scope_type")
		if actual != expected && actual != "global" {
			response.Fail(c, 400, int(appErr.ErrBadRequest), fmt.Sprintf("scope mismatch: expected %s got %v", expected, actual))
			c.Abort()
			return
		}
		c.Next()
	}
}
