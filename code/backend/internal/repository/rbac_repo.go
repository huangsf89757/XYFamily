package repository

import (
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
