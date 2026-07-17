package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xyfamily/internal/model"
)

type TenantRepository struct {
	db *DB
}

func NewTenantRepository(db *DB) *TenantRepository { return &TenantRepository{db: db} }

func (r *TenantRepository) CreateOrganization(ctx context.Context, o *model.Organization) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO organizations (org_id, name, description, owner_id) VALUES ($1,$2,$3,$4)", o.OrgID, o.Name, o.Description, o.OwnerID)
	return err
}
func (r *TenantRepository) GetOrganizationByID(ctx context.Context, orgID uuid.UUID) (*model.Organization, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id,org_id,name,description,owner_id,created_at,updated_at,disabled_at,deleted_at FROM organizations WHERE org_id=$1 AND deleted_at IS NULL", orgID)
	var o model.Organization
	err := row.Scan(&o.ID,&o.OrgID,&o.Name,&o.Description,&o.OwnerID,&o.CreatedAt,&o.UpdatedAt,&o.DisabledAt,&o.DeletedAt)
	if err == pgx.ErrNoRows { return nil, nil }
	if err != nil { return nil, fmt.Errorf("scan org: %w", err) }
	return &o, nil
}
func (r *TenantRepository) UpdateOrganization(ctx context.Context, orgID uuid.UUID, name, desc string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE organizations SET name=$1, description=$2, updated_at=now() WHERE org_id=$3 AND deleted_at IS NULL", name, desc, orgID)
	return err
}
func (r *TenantRepository) DisableOrganization(ctx context.Context, orgID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE organizations SET disabled_at=now() WHERE org_id=$1 AND disabled_at IS NULL AND deleted_at IS NULL", orgID)
	return err
}
func (r *TenantRepository) EnableOrganization(ctx context.Context, orgID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE organizations SET disabled_at=NULL, updated_at=now() WHERE org_id=$1 AND disabled_at IS NOT NULL AND deleted_at IS NULL", orgID)
	return err
}
func (r *TenantRepository) CreateOrgMember(ctx context.Context, m *model.OrgMember) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO org_members (org_id, account_id, role) VALUES ($1,$2,$3) ON CONFLICT (org_id, account_id) DO NOTHING", m.OrgID, m.AccountID, m.Role)
	return err
}
func (r *TenantRepository) ListOrgMembers(ctx context.Context, orgID uuid.UUID, page, size int) ([]model.OrgMember, int64, error) {
	offset := (page-1)*size
	var total int64
	_ = r.db.Pool.QueryRow(ctx, "SELECT count(*) FROM org_members WHERE org_id=$1 AND deleted_at IS NULL", orgID).Scan(&total)
	rows, err := r.db.Pool.Query(ctx, "SELECT id,org_id,account_id,role,joined_at,created_at,updated_at FROM org_members WHERE org_id=$1 AND deleted_at IS NULL LIMIT $2 OFFSET $3", orgID, size, offset)
	if err != nil { return nil, 0, fmt.Errorf("query: %w", err) }
	defer rows.Close()
	var members []model.OrgMember
	for rows.Next() { var m model.OrgMember; if err := rows.Scan(&m.ID,&m.OrgID,&m.AccountID,&m.Role,&m.JoinedAt,&m.CreatedAt,&m.UpdatedAt); err != nil { return nil, 0, err }; members = append(members, m) }
	return members, total, nil
}
func (r *TenantRepository) UpdateOrgMemberRole(ctx context.Context, orgID, accountID uuid.UUID, role string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE org_members SET role=$1, updated_at=now() WHERE org_id=$2 AND account_id=$3 AND deleted_at IS NULL", role, orgID, accountID)
	return err
}
func (r *TenantRepository) RemoveOrgMember(ctx context.Context, orgID, accountID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE org_members SET deleted_at=now() WHERE org_id=$1 AND account_id=$2 AND deleted_at IS NULL", orgID, accountID)
	return err
}
func (r *TenantRepository) CountCoreAdmins(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64; err := r.db.Pool.QueryRow(ctx, "SELECT count(*) FROM org_members WHERE org_id=$1 AND role=$2 AND deleted_at IS NULL", orgID, "organization_core_admin").Scan(&count); return count, err
}
func (r *TenantRepository) CreateTeam(ctx context.Context, t *model.Team) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO teams (team_id, org_id, name, description, owner_id) VALUES ($1,$2,$3,$4,$5)", t.TeamID, t.OrgID, t.Name, t.Description, t.OwnerID)
	return err
}
func (r *TenantRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id,team_id,org_id,name,description,owner_id,created_at,updated_at,archived_at,deleted_at FROM teams WHERE team_id=$1 AND deleted_at IS NULL", teamID)
	var t model.Team; err := row.Scan(&t.ID,&t.TeamID,&t.OrgID,&t.Name,&t.Description,&t.OwnerID,&t.CreatedAt,&t.UpdatedAt,&t.ArchivedAt,&t.DeletedAt)
	if err == pgx.ErrNoRows { return nil, nil }; if err != nil { return nil, fmt.Errorf("scan team: %w", err) }; return &t, nil
}
func (r *TenantRepository) UpdateTeam(ctx context.Context, teamID uuid.UUID, name, desc string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE teams SET name=$1, description=$2, updated_at=now() WHERE team_id=$3 AND deleted_at IS NULL", name, desc, teamID); return err
}
func (r *TenantRepository) ArchiveTeam(ctx context.Context, teamID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE teams SET archived_at=now() WHERE team_id=$1 AND archived_at IS NULL AND deleted_at IS NULL", teamID); return err
}
func (r *TenantRepository) CreateGroup(ctx context.Context, g *model.Group) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO groups (group_id, org_id, team_id, name, description, owner_id) VALUES ($1,$2,$3,$4,$5,$6)", g.GroupID, g.OrgID, g.TeamID, g.Name, g.Description, g.OwnerID); return err
}
func (r *TenantRepository) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*model.Group, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id,group_id,org_id,team_id,name,description,owner_id,created_at,updated_at,deleted_at FROM groups WHERE group_id=$1 AND deleted_at IS NULL", groupID)
	var g model.Group; err := row.Scan(&g.ID,&g.GroupID,&g.OrgID,&g.TeamID,&g.Name,&g.Description,&g.OwnerID,&g.CreatedAt,&g.UpdatedAt,&g.DeletedAt)
	if err == pgx.ErrNoRows { return nil, nil }; if err != nil { return nil, fmt.Errorf("scan group: %w", err) }; return &g, nil
}
func (r *TenantRepository) UpdateGroup(ctx context.Context, groupID uuid.UUID, name, desc string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE groups SET name=$1, description=$2, updated_at=now() WHERE group_id=$3 AND deleted_at IS NULL", name, desc, groupID); return err
}
func (r *TenantRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE groups SET deleted_at=now() WHERE group_id=$1 AND deleted_at IS NULL", groupID); return err
}
func (r *TenantRepository) CreateInvitation(ctx context.Context, inv *model.Invitation) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO invitations (org_id, scope_type, scope_id, inviter_id, invitee_contact, role, token, status, expired_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)", inv.OrgID, inv.ScopeType, inv.ScopeID, inv.InviterID, inv.InviteeContact, inv.Role, inv.Token, inv.Status, inv.ExpiredAt)
	return err
}
func (r *TenantRepository) GetInvitationByToken(ctx context.Context, token string) (*model.Invitation, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id,org_id,scope_type,scope_id,inviter_id,invitee_contact,role,token,status,expired_at,created_at,accepted_at,deleted_at FROM invitations WHERE token=$1 AND deleted_at IS NULL", token)
	var inv model.Invitation
	err := row.Scan(&inv.ID,&inv.OrgID,&inv.ScopeType,&inv.ScopeID,&inv.InviterID,&inv.InviteeContact,&inv.Role,&inv.Token,&inv.Status,&inv.ExpiredAt,&inv.CreatedAt,&inv.AcceptedAt,&inv.DeletedAt)
	if err == pgx.ErrNoRows { return nil, nil }
	if err != nil { return nil, fmt.Errorf("scan invitation: %w", err) }
	return &inv, nil
}
func (r *TenantRepository) AcceptInvitation(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE invitations SET status=$1, accepted_at=now() WHERE id=$2 AND status=$3", "accepted", id, "pending")
	return err
}
