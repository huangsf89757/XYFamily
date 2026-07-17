package repository

import (
	"xyfamily/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

)

type RBACRepository struct {
	db *DB
}

func NewRBACRepository(db *DB) *RBACRepository {
	return &RBACRepository{db: db}
}

func (r *RBACRepository) GetPermissionsByRoleKey(ctx context.Context, roleKey string) ([]string, error) {
	query := "SELECT permission_key FROM role_permissions WHERE role_key = $1"
	rows, err := r.db.Pool.Query(ctx, query, roleKey)
	if err != nil {
		return nil, fmt.Errorf("query permissions: %w", err)
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *RBACRepository) GetOrgRoleByAccountID(ctx context.Context, accountID, orgID uuid.UUID) (string, error) {
	query := "SELECT role FROM org_members WHERE account_id = $1 AND org_id = $2 AND deleted_at IS NULL"
	row := r.db.Pool.QueryRow(ctx, query, accountID, orgID)
	var role string
	err := row.Scan(&role)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("scan role: %w", err)
	}
	return role, nil
}

func (r *RBACRepository) GetOrgIDsByAccountID(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error) {
	query := "SELECT org_id FROM org_members WHERE account_id = $1 AND deleted_at IS NULL"
	rows, err := r.db.Pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("query org ids: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan org id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *RBACRepository) GetAllPermissions(ctx context.Context) ([]string, error) {
	query := "SELECT permission_key FROM permission_points ORDER BY permission_key"
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all permissions: %w", err)
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, nil
}
func (r *RBACRepository) GetTeamRoleByAccountID(ctx context.Context, accountID, teamID uuid.UUID) (string, error) {
	query := "SELECT role FROM team_members WHERE account_id = $1 AND team_id = $2 AND deleted_at IS NULL"
	row := r.db.Pool.QueryRow(ctx, query, accountID, teamID)
	var role string; err := row.Scan(&role)
	if err == pgx.ErrNoRows { return "", nil }; if err != nil { return "", fmt.Errorf("scan role: %w", err) }; return role, nil
}
func (r *RBACRepository) GetGroupRoleByAccountID(ctx context.Context, accountID, groupID uuid.UUID) (string, error) {
	query := "SELECT role FROM group_members WHERE account_id = $1 AND group_id = $2 AND deleted_at IS NULL"
	row := r.db.Pool.QueryRow(ctx, query, accountID, groupID)
	var role string; err := row.Scan(&role)
	if err == pgx.ErrNoRows { return "", nil }; if err != nil { return "", fmt.Errorf("scan role: %w", err) }; return role, nil
}
func (r *RBACRepository) GetTeamMembersByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.TeamMember, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT id, org_id, team_id, account_id, role, joined_at, created_at, updated_at FROM team_members WHERE account_id = $1 AND deleted_at IS NULL", accountID)
	if err != nil { return nil, fmt.Errorf("query team members: %w", err) }
	defer rows.Close()
	var members []model.TeamMember
	for rows.Next() { var m model.TeamMember; if err := rows.Scan(&m.ID,&m.OrgID,&m.TeamID,&m.AccountID,&m.Role,&m.JoinedAt,&m.CreatedAt,&m.UpdatedAt); err != nil { return nil, err }; members = append(members, m) }
	return members, nil
}
func (r *RBACRepository) GetGroupMembersByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.GroupMember, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT id, org_id, team_id, group_id, account_id, role, joined_at, created_at, updated_at FROM group_members WHERE account_id = $1 AND deleted_at IS NULL", accountID)
	if err != nil { return nil, fmt.Errorf("query group members: %w", err) }
	defer rows.Close()
	var members []model.GroupMember
	for rows.Next() { var m model.GroupMember; if err := rows.Scan(&m.ID,&m.OrgID,&m.TeamID,&m.GroupID,&m.AccountID,&m.Role,&m.JoinedAt,&m.CreatedAt,&m.UpdatedAt); err != nil { return nil, err }; members = append(members, m) }
	return members, nil
}
func (r *RBACRepository) CountSuperAdmins(ctx context.Context) (int64, error) {
	var count int64; err := r.db.Pool.QueryRow(ctx, "SELECT count(*) FROM org_members WHERE role = $1 AND deleted_at IS NULL", "super_admin").Scan(&count); return count, err
}
