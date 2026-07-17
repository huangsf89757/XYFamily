package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xyfamily/internal/service"
	"xyfamily/pkg/response"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(ts *service.TeamService) *TeamHandler { return &TeamHandler{teamService: ts} }

func (h *TeamHandler) Create(c *gin.Context) {
	orgID, _ := uuid.Parse(c.GetHeader("X-Organization-ID"))
	accountID, _ := c.Get("account_id"); creatorID, _ := uuid.Parse(accountID.(string))
	var req service.CreateTeamReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 400001, "invalid request"); return }
	resp, err := h.teamService.Create(c.Request.Context(), orgID, creatorID, &req)
	if err != nil { response.Fail(c, 400, 400001, err.Error()); return }
	response.Created(c, resp)
}

func (h *TeamHandler) GetInfo(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("team_id"))
	if err != nil { response.Fail(c, 400, 400001, "invalid team id"); return }
	resp, err := h.teamService.GetInfo(c.Request.Context(), teamID)
	if err != nil { response.Fail(c, 404, 400003, err.Error()); return }
	response.OK(c, resp)
}

func (h *TeamHandler) Update(c *gin.Context) {
	teamID, _ := uuid.Parse(c.Param("team_id"))
	var req struct{ Name string `json:"name"`; Description string `json:"description"` }
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 400001, "invalid request"); return }
	if err := h.teamService.Update(c.Request.Context(), teamID, req.Name, req.Description); err != nil { response.Fail(c, 400, 400004, err.Error()); return }
	response.OK(c, gin.H{"team_id": teamID, "name": req.Name, "description": req.Description})
}

func (h *TeamHandler) Archive(c *gin.Context) {
	teamID, _ := uuid.Parse(c.Param("team_id"))
	if err := h.teamService.Archive(c.Request.Context(), teamID); err != nil { response.Fail(c, 400, 400001, err.Error()); return }
	response.OK(c, gin.H{"team_id": teamID, "archived_at": "now"})
}

func (h *TeamHandler) CreateGroup(c *gin.Context) {
	orgID, _ := uuid.Parse(c.GetHeader("X-Organization-ID"))
	teamID, _ := uuid.Parse(c.Param("team_id"))
	accountID, _ := c.Get("account_id"); creatorID, _ := uuid.Parse(accountID.(string))
	var req service.CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 400001, "invalid request"); return }
	resp, err := h.teamService.CreateGroup(c.Request.Context(), orgID, teamID, creatorID, &req)
	if err != nil { response.Fail(c, 400, 400009, err.Error()); return }
	response.Created(c, resp)
}

func (h *TeamHandler) GetGroup(c *gin.Context) {
	groupID, _ := uuid.Parse(c.Param("group_id"))
	g, err := h.teamService.GetGroup(c.Request.Context(), groupID)
	if err != nil { response.Fail(c, 404, 400003, err.Error()); return }
	response.OK(c, g)
}

func (h *TeamHandler) UpdateGroup(c *gin.Context) {
	groupID, _ := uuid.Parse(c.Param("group_id"))
	var req struct{ Name string `json:"name"`; Description string `json:"description"` }
	if err := c.ShouldBindJSON(&req); err != nil { response.Fail(c, 400, 400001, "invalid request"); return }
	if err := h.teamService.UpdateGroup(c.Request.Context(), groupID, req.Name, req.Description); err != nil { response.Fail(c, 400, 400004, err.Error()); return }
	response.OK(c, gin.H{"group_id": groupID, "name": req.Name, "description": req.Description})
}

func (h *TeamHandler) DeleteGroup(c *gin.Context) {
	groupID, _ := uuid.Parse(c.Param("group_id"))
	if err := h.teamService.DeleteGroup(c.Request.Context(), groupID); err != nil { response.Fail(c, 400, 400001, err.Error()); return }
	response.OK(c, nil)
}
