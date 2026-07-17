package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"xyfamily/internal/repository"
	"xyfamily/pkg/response"
)

type HealthHandler struct {
	db    *repository.DB
	cache *repository.RedisClient
}

func NewHealthHandler(db *repository.DB, cache *repository.RedisClient) *HealthHandler {
	return &HealthHandler{db: db, cache: cache}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	status := gin.H{"status": "ok"}
	allOK := true

	if err := h.db.Pool.Ping(ctx); err != nil {
		status["database"] = "unhealthy"
		allOK = false
	} else {
		status["database"] = "healthy"
	}

	if err := h.cache.Client.Ping(ctx).Err(); err != nil {
		status["redis"] = "unhealthy"
		allOK = false
	} else {
		status["redis"] = "healthy"
	}

	if allOK {
		c.JSON(http.StatusOK, response.Response{Code: 0, Message: "ok", Data: status})
	} else {
		c.JSON(http.StatusServiceUnavailable, response.Response{Code: 800001, Message: "unhealthy", Data: status})
	}
}
