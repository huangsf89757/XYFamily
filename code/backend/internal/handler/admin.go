package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/service"
	"xyfamily/pkg/response"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(as *service.AdminService) *AdminHandler { return &AdminHandler{adminService: as} }

func (h *AdminHandler) Init(c *gin.Context) {
	var req service.InitReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 100001, "invalid request"); return }
	resp, err := h.adminService.Init(c.Request.Context(), &req)
	if err != nil { response.Fail(c, 400, 800010, err.Error()); return }
	response.Created(c, resp)
}

func (h *AdminHandler) GetConfig(c *gin.Context) {
	resp, err := h.adminService.GetConfig(c.Request.Context())
	if err != nil { response.Fail(c, 500, 800001, err.Error()); return }
	response.OK(c, resp)
}

func (h *AdminHandler) UpdateConfig(c *gin.Context) {
	var req service.UpdateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 800011, "invalid request"); return }
	if err := h.adminService.UpdateConfig(c.Request.Context(), &req); err != nil { response.Fail(c, 400, 800011, err.Error()); return }
	response.OK(c, gin.H{"config_key": req.ConfigKey, "config_value": req.ConfigValue, "updated_at": "now"})
}

func (h *AdminHandler) ForceDowngrade(c *gin.Context) {
	var req service.ForceDowngradeReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 800011, "invalid request"); return }
	resp, err := h.adminService.ForceDowngrade(c.Request.Context(), &req)
	if err != nil { response.Fail(c, 400, 800012, err.Error()); return }
	response.OK(c, resp)
}

func (h *AdminHandler) GlobalAuditList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }; if size < 1 || size > 100 { size = 20 }
	var orgID, accountID *uuid.UUID
	if oid := c.Query("org_id"); oid != "" { id, _ := uuid.Parse(oid); orgID = &id }
	if aid := c.Query("account_id"); aid != "" { id, _ := uuid.Parse(aid); accountID = &id }
	actionType := c.Query("action_type")
	result := c.Query("result")
	var start, end time.Time
	if s := c.Query("start"); s != "" { start, _ = time.Parse(time.RFC3339, s) }
	if e := c.Query("end"); e != "" { end, _ = time.Parse(time.RFC3339, e) }
	logs, total, err := h.adminService.GlobalAuditList(c.Request.Context(), orgID, accountID, actionType, result, start, end, page, size)
	if err != nil { response.Fail(c, 500, 800001, err.Error()); return }
	response.OK(c, gin.H{"items": logs, "total": total})
}

func (h *AdminHandler) AuditDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { response.Fail(c, 400, 700002, "invalid id"); return }
	log, err := h.adminService.AuditDetail(c.Request.Context(), id)
	if err != nil { response.Fail(c, 404, 700003, err.Error()); return }
	response.OK(c, log)
}
