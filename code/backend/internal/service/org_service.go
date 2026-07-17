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

type OrgService struct {
	tenantRepo *repository.TenantRepository
	cacheRepo  *repository.CacheRepository
	auditRepo  *repository.AuditRepository
}

func NewOrgService(tr *repository.TenantRepository, cr *repository.CacheRepository, aur *repository.AuditRepository) *OrgService {
	return &OrgService{tenantRepo: tr, cacheRepo: cr, auditRepo: aur}
}

type CreateOrgReq struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	InitialAdmin string `json:"initial_admin"`
}

type CreateOrgResp struct {
	OrgID   string `json:"org_id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (s *OrgService) Create(ctx context.Context, req *CreateOrgReq, creatorID uuid.UUID) (*CreateOrgResp, error) {
	orgID := uuid.New()
	org := &model.Organization{OrgID: orgID, Name: req.Name, Description: req.Description, OwnerID: creatorID}
	if err := s.tenantRepo.CreateOrganization(ctx, org); err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	member := &model.OrgMember{OrgID: orgID, AccountID: creatorID, Role: "organization_core_admin"}
	_ = s.tenantRepo.CreateOrgMember(ctx, member)
	return &CreateOrgResp{OrgID: orgID.String(), Name: req.Name, OwnerID: creatorID.String()}, nil
}
type OrgInfoResp struct {
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}
func (s *OrgService) GetInfo(ctx context.Context, orgID uuid.UUID) (*OrgInfoResp, error) {
	org, err := s.tenantRepo.GetOrganizationByID(ctx, orgID)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	if org == nil { return nil, fmt.Errorf("%w: not found", appErr.ErrNotFound) }
	return &OrgInfoResp{OrgID: org.OrgID.String(), Name: org.Name, Description: org.Description, OwnerID: org.OwnerID.String(), CreatedAt: org.CreatedAt}, nil
}
type UpdateOrgReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
func (s *OrgService) Update(ctx context.Context, orgID uuid.UUID, req *UpdateOrgReq) error {
	return s.tenantRepo.UpdateOrganization(ctx, orgID, req.Name, req.Description)
}
func (s *OrgService) Disable(ctx context.Context, orgID uuid.UUID) error { return s.tenantRepo.DisableOrganization(ctx, orgID) }
func (s *OrgService) Enable(ctx context.Context, orgID uuid.UUID) error { return s.tenantRepo.EnableOrganization(ctx, orgID) }
type InviteReq struct { Invitee string `json:"invitee"`; Role string `json:"role"` }
type InviteResp struct { InvitationID string `json:"invitation_id"`; Status string `json:"status"` }
func (s *OrgService) Invite(ctx context.Context, orgID, inviterID uuid.UUID, req *InviteReq) (*InviteResp, error) {
	invID := uuid.New(); token := uuid.New().String()
	inv := &model.Invitation{ID: invID, OrgID: orgID, ScopeType: "organization", ScopeID: orgID, InviterID: inviterID, InviteeContact: req.Invitee, Role: req.Role, Token: token, Status: "pending", ExpiredAt: time.Now().Add(7*24*time.Hour)}
	if err := s.tenantRepo.CreateInvitation(ctx, inv); err != nil { return nil, fmt.Errorf("%w: invite failed", appErr.ErrInternal) }
	return &InviteResp{InvitationID: invID.String(), Status: "pending"}, nil
}
type MemberItem struct { AccountID string `json:"account_id"`; Role string `json:"role"`; JoinedAt time.Time `json:"joined_at"` }
type ListMembersResp struct { Items []MemberItem `json:"items"`; Total int64 `json:"total"` }
func (s *OrgService) ListMembers(ctx context.Context, orgID uuid.UUID, page, size int) (*ListMembersResp, error) {
	members, total, err := s.tenantRepo.ListOrgMembers(ctx, orgID, page, size)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	items := make([]MemberItem, 0, len(members))
	for _, m := range members { items = append(items, MemberItem{AccountID: m.AccountID.String(), Role: m.Role, JoinedAt: m.JoinedAt}) }
	return &ListMembersResp{Items: items, Total: total}, nil
}
func (s *OrgService) AssignRole(ctx context.Context, orgID, accountID uuid.UUID, role string) error { return s.tenantRepo.UpdateOrgMemberRole(ctx, orgID, accountID, role) }
func (s *OrgService) Downgrade(ctx context.Context, orgID, accountID uuid.UUID) error {
	count, _ := s.tenantRepo.CountCoreAdmins(ctx, orgID)
	if count <= 1 { return fmt.Errorf("%w: last admin", appErr.ErrForbidden) }
	return s.tenantRepo.UpdateOrgMemberRole(ctx, orgID, accountID, "regular_member")
}
func (s *OrgService) RemoveMember(ctx context.Context, orgID, accountID uuid.UUID) error {
	count, _ := s.tenantRepo.CountCoreAdmins(ctx, orgID)
	if count <= 1 { return fmt.Errorf("%w: last admin", appErr.ErrForbidden) }
	return s.tenantRepo.RemoveOrgMember(ctx, orgID, accountID)
}
