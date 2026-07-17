package middleware

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/repository"
	"xyfamily/pkg/logger"
	"go.uber.org/zap"
)

type AuditMiddleware struct {
	cacheRepo *repository.CacheRepository
}

func NewAuditMiddleware(cache *repository.CacheRepository) *AuditMiddleware {
	return &AuditMiddleware{cacheRepo: cache}
}

func (m *AuditMiddleware) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		go m.pushEvent(c)
	}
}

func (m *AuditMiddleware) pushEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accountID, _ := c.Get("account_id")
	orgID, _ := c.Get("org_id")
	var acctIDStr, orgIDStr string
	if accountID != nil { acctIDStr, _ = accountID.(string) }
	if orgID != nil { orgIDStr = orgID.(uuid.UUID).String() }
	event := map[string]interface{}{
		"event_id":      uuid.New().String(),
		"account_id":    acctIDStr,
		"org_id":        orgIDStr,
		"action_domain":  "operation",
		"action_type":    c.Request.Method + ":" + c.FullPath(),
		"target_type":    "",
		"target_id":     "",
		"result":        "success",
		"trace_id":      c.GetString("trace_id"),
		"ip_address":    c.ClientIP(),
		"user_agent":     c.GetHeader("User-Agent"),
	}
	if c.Writer.Status() >= 400 {
		event["result"] = "failed"
		event["failure_reason"] = "http_error"
	}
	details, _ := json.Marshal(map[string]interface{}{"method": c.Request.Method, "path": c.Request.URL.Path, "status": c.Writer.Status()})
	event["details"] = string(details)
	if err := m.cacheRepo.PushAuditEvent(ctx, event); err != nil {
		logger.Get().Error("push audit event failed", zap.Error(err))
	}
}
