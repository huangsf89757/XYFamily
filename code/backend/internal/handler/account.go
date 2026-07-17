package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/service"
	"xyfamily/pkg/response"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(as *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: as}
}

func (h *AccountHandler) GetProfile(c *gin.Context) {
	accountID, _ := c.Get("account_id")
	id, _ := uuid.Parse(accountID.(string))
	resp, err := h.accountService.GetProfile(c.Request.Context(), id)
	if err != nil { response.Fail(c, 404, 200010, err.Error()); return }
	response.OK(c, resp)
}

func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	var req service.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 100001, "invalid request"); return }
	accountID, _ := c.Get("account_id"); id, _ := uuid.Parse(accountID.(string))
	if err := h.accountService.UpdateProfile(c.Request.Context(), id, &req); err != nil { response.Fail(c, 400, 100001, err.Error()); return }
	response.OK(c, gin.H{"account_id": accountID, "nickname": req.Nickname, "avatar": req.Avatar})
}

func (h *AccountHandler) ChangePassword(c *gin.Context) {
	var req service.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 100001, "invalid request"); return }
	accountID, _ := c.Get("account_id"); id, _ := uuid.Parse(accountID.(string))
	if err := h.accountService.ChangePassword(c.Request.Context(), id, &req); err != nil { response.Fail(c, 400, 100001, err.Error()); return }
	response.OK(c, nil)
}

func (h *AccountHandler) Deactivate(c *gin.Context) {
	accountID, _ := c.Get("account_id"); id, _ := uuid.Parse(accountID.(string))
	status, deactivatedAt, err := h.accountService.Deactivate(c.Request.Context(), id)
	if err != nil { response.Fail(c, 400, 200007, err.Error()); return }
	response.OK(c, gin.H{"status": status, "deactivated_at": deactivatedAt.Format(time.RFC3339)})
}

func (h *AccountHandler) Undeactivate(c *gin.Context) {
	var req service.UndeactivateReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 100001, "invalid request"); return }
	accountID, _ := c.Get("account_id"); id, _ := uuid.Parse(accountID.(string))
	if err := h.accountService.Undeactivate(c.Request.Context(), id, &req); err != nil { response.Fail(c, 400, 100001, err.Error()); return }
	response.OK(c, gin.H{"status": "active"})
}
