package service

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"xyfamily/internal/model"
	"xyfamily/internal/repository"
	appErr "xyfamily/pkg/errors"
)

type TeamService struct {
	tenantRepo *repository.TenantRepository
	auditRepo  *repository.AuditRepository
}

func NewTeamService(tr *repository.TenantRepository, aur *repository.AuditRepository) *TeamService {
	return &TeamService{tenantRepo: tr, auditRepo: aur}
}

type CreateTeamReq struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	InitialAdmin string `json:"initial_admin"`
}
type CreateTeamResp struct {
	TeamID string `json:"team_id"`
	OrgID  string `json:"org_id"`
	Name   string `json:"name"`
}

func (s *TeamService) Create(ctx context.Context, orgID, creatorID uuid.UUID, req *CreateTeamReq) (*CreateTeamResp, error) {
	teamID := uuid.New()
	team := &model.Team{TeamID: teamID, OrgID: orgID, Name: req.Name, Description: req.Description, OwnerID: creatorID}
	if err := s.tenantRepo.CreateTeam(ctx, team); err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	return &CreateTeamResp{TeamID: teamID.String(), OrgID: orgID.String(), Name: req.Name}, nil
}

type TeamInfoResp struct {
	TeamID      string     `json:"team_id"`
	OrgID       string     `json:"org_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

func (s *TeamService) GetInfo(ctx context.Context, teamID uuid.UUID) (*TeamInfoResp, error) {
	t, err := s.tenantRepo.GetTeamByID(ctx, teamID)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	if t == nil { return nil, fmt.Errorf("%w: not found", appErr.ErrNotFound) }
	return &TeamInfoResp{TeamID: t.TeamID.String(), OrgID: t.OrgID.String(), Name: t.Name, Description: t.Description, ArchivedAt: t.ArchivedAt}, nil
}

func (s *TeamService) Update(ctx context.Context, teamID uuid.UUID, name, desc string) error {
	return s.tenantRepo.UpdateTeam(ctx, teamID, name, desc)
}

func (s *TeamService) Archive(ctx context.Context, teamID uuid.UUID) error {
	return s.tenantRepo.ArchiveTeam(ctx, teamID)
}

type CreateGroupReq struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	InitialAdmin string `json:"initial_admin"`
}
type CreateGroupResp struct {
	GroupID string `json:"group_id"`
	TeamID  string `json:"team_id"`
	Name    string `json:"name"`
}

func (s *TeamService) CreateGroup(ctx context.Context, orgID, teamID, creatorID uuid.UUID, req *CreateGroupReq) (*CreateGroupResp, error) {
	groupID := uuid.New()
	group := &model.Group{GroupID: groupID, OrgID: orgID, TeamID: teamID, Name: req.Name, Description: req.Description, OwnerID: creatorID}
	if err := s.tenantRepo.CreateGroup(ctx, group); err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	return &CreateGroupResp{GroupID: groupID.String(), TeamID: teamID.String(), Name: req.Name}, nil
}

func (s *TeamService) GetGroup(ctx context.Context, groupID uuid.UUID) (*model.Group, error) {
	g, err := s.tenantRepo.GetGroupByID(ctx, groupID)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	if g == nil { return nil, fmt.Errorf("%w: not found", appErr.ErrNotFound) }
	return g, nil
}

func (s *TeamService) UpdateGroup(ctx context.Context, groupID uuid.UUID, name, desc string) error {
	return s.tenantRepo.UpdateGroup(ctx, groupID, name, desc)
}

func (s *TeamService) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	return s.tenantRepo.DeleteGroup(ctx, groupID)
}
