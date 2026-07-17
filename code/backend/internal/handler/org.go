package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/service"
	"xyfamily/pkg/response"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler(os *service.OrgService) *OrgHandler { return &OrgHandler{orgService: os} }

func (h *OrgHandler) Create(c *gin.Context) {
	var req service.CreateOrgReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 300001, "invalid request"); return }
	accountID, _ := c.Get("account_id"); creatorID, _ := uuid.Parse(accountID.(string))
	resp, err := h.orgService.Create(c.Request.Context(), &req, creatorID)
	if err != nil { response.Fail(c, 400, 300001, err.Error()); return }
	response.Created(c, resp)
}

func (h *OrgHandler) GetInfo(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("organization_id"))
	if err != nil { response.Fail(c, 400, 300001, "invalid org id"); return }
	resp, err := h.orgService.GetInfo(c.Request.Context(), orgID)
	if err != nil { response.Fail(c, 404, 300003, err.Error()); return }
	response.OK(c, resp)
}

func (h *OrgHandler) Update(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	var req service.UpdateOrgReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 300001, "invalid request"); return }
	if err := h.orgService.Update(c.Request.Context(), orgID, &req); err != nil { response.Fail(c, 400, 300004, err.Error()); return }
	response.OK(c, gin.H{"org_id": orgID, "name": req.Name, "description": req.Description})
}

func (h *OrgHandler) Disable(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	if err := h.orgService.Disable(c.Request.Context(), orgID); err != nil { response.Fail(c, 400, 300010, err.Error()); return }
	response.OK(c, gin.H{"org_id": orgID, "disabled_at": "now"})
}

func (h *OrgHandler) Enable(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	if err := h.orgService.Enable(c.Request.Context(), orgID); err != nil { response.Fail(c, 400, 300010, err.Error()); return }
	response.OK(c, gin.H{"org_id": orgID, "disabled_at": nil})
}

func (h *OrgHandler) Invite(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	var req service.InviteReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 300001, "invalid request"); return }
	accountID, _ := c.Get("account_id"); inviterID, _ := uuid.Parse(accountID.(string))
	resp, err := h.orgService.Invite(c.Request.Context(), orgID, inviterID, &req)
	if err != nil { response.Fail(c, 400, 300005, err.Error()); return }
	response.Created(c, resp)
}

func (h *OrgHandler) ListMembers(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }; if size < 1 || size > 100 { size = 20 }
	resp, err := h.orgService.ListMembers(c.Request.Context(), orgID, page, size)
	if err != nil { response.Fail(c, 500, 800001, err.Error()); return }
	response.OK(c, resp)
}

func (h *OrgHandler) AssignRole(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	accountID, _ := uuid.Parse(c.Param("account_id"))
	var req struct{ Role string `json:"role"` }
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 300001, "invalid request"); return }
	if err := h.orgService.AssignRole(c.Request.Context(), orgID, accountID, req.Role); err != nil { response.Fail(c, 400, 300007, err.Error()); return }
	response.OK(c, gin.H{"account_id": accountID, "role": req.Role})
}

func (h *OrgHandler) Downgrade(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	accountID, _ := uuid.Parse(c.Param("account_id"))
	if err := h.orgService.Downgrade(c.Request.Context(), orgID, accountID); err != nil { response.Fail(c, 400, 300008, err.Error()); return }
	response.OK(c, gin.H{"account_id": accountID, "role": "regular_member"})
}

func (h *OrgHandler) RemoveMember(c *gin.Context) {
	orgID, _ := uuid.Parse(c.Param("organization_id"))
	accountID, _ := uuid.Parse(c.Param("account_id"))
	if err := h.orgService.RemoveMember(c.Request.Context(), orgID, accountID); err != nil { response.Fail(c, 400, 300006, err.Error()); return }
	response.OK(c, nil)
}
