package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"xyfamily/internal/middleware"
	"xyfamily/internal/service"
	"xyfamily/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
	rateLimit   *middleware.RateLimitMiddleware
}

func NewAuthHandler(authService *service.AuthService, rateLimit *middleware.RateLimitMiddleware) *AuthHandler {
	return &AuthHandler{authService: authService, rateLimit: rateLimit}
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req service.SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, "invalid request")
		return
	}
	resp, err := h.authService.SendVerificationCode(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, err.Error())
		return
	}
	response.OK(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, "invalid request")
		return
	}
	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, err.Error())
		return
	}
	response.Created(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, "invalid request")
		return
	}
	resp, err := h.authService.Login(c.Request.Context(), &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		h.rateLimit.RecordLoginFailure(c)
		response.Fail(c, http.StatusUnauthorized, 101008, err.Error())
		return
	}
	response.OK(c, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req service.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, "invalid request")
		return
	}
	resp, err := h.authService.RefreshToken(c.Request.Context(), &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 101004, err.Error())
		return
	}
	response.OK(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req service.LogoutReq
	_ = c.ShouldBindJSON(&req)
	accountID, _ := c.Get("account_id")
	jti, _ := c.Get("jti")
	accountIDStr, _ := accountID.(string)
	jtiStr, _ := jti.(string)
	_ = h.authService.Logout(c.Request.Context(), &req, accountIDStr, jtiStr)
	response.OK(c, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req service.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, "invalid request")
		return
	}
	err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 100001, err.Error())
		return
	}
	response.OK(c, nil)
}
