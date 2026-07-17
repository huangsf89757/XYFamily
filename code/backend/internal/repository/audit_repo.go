package repository

import (
	"time"
	"context"
	"fmt"

	"github.com/google/uuid"

	"xyfamily/internal/model"
)

type AuditRepository struct {
	db *DB
}

func NewAuditRepository(db *DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	query := "INSERT INTO audit_logs (event_id, account_id, org_id, action_domain, action_type, target_type, target_id, result, failure_reason, login_method, details, trace_id, ip_address, user_agent) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
	_, err := r.db.Pool.Exec(ctx, query,
		log.EventID, log.AccountID, log.OrgID, log.ActionDomain, log.ActionType,
		log.TargetType, log.TargetID, log.Result, log.FailureReason, log.LoginMethod,
		log.Details, log.TraceID, log.IPAddress, log.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *AuditRepository) ListByAccountID(ctx context.Context, accountID uuid.UUID, page, pageSize int) ([]model.AuditLog, int64, error) {
	offset := (page - 1) * pageSize
	var total int64
	err := r.db.Pool.QueryRow(ctx, "SELECT count(*) FROM audit_logs WHERE account_id = $1", accountID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	rows, err := r.db.Pool.Query(ctx, "SELECT id, event_id, account_id, org_id, action_domain, action_type, target_type, target_id, result, failure_reason, login_method, details, trace_id, ip_address, user_agent, created_at FROM audit_logs WHERE account_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", accountID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()
	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		if err := rows.Scan(&l.ID, &l.EventID, &l.AccountID, &l.OrgID, &l.ActionDomain, &l.ActionType, &l.TargetType, &l.TargetID, &l.Result, &l.FailureReason, &l.LoginMethod, &l.Details, &l.TraceID, &l.IPAddress, &l.UserAgent, &l.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan audit log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}
func (r *AuditRepository) ListGlobal(ctx context.Context, orgID, accountID *uuid.UUID, actionType, result string, start, end time.Time, page, size int) ([]model.AuditLog, int64, error) {
	offset := (page-1)*size
	query := "SELECT id, event_id, account_id, org_id, action_domain, action_type, target_type, target_id, result, failure_reason, login_method, details, trace_id, ip_address, user_agent, created_at FROM audit_logs WHERE 1=1"
	args := []interface{}{}
	n := 1
	if orgID != nil { query += fmt.Sprintf(" AND org_id = $%d", n); args = append(args, *orgID); n++ }
	if accountID != nil { query += fmt.Sprintf(" AND account_id = $%d", n); args = append(args, *accountID); n++ }
	if actionType != "" { query += fmt.Sprintf(" AND action_type = $%d", n); args = append(args, actionType); n++ }
	if result != "" { query += fmt.Sprintf(" AND result = $%d", n); args = append(args, result); n++ }
	if !start.IsZero() { query += fmt.Sprintf(" AND created_at >= $%d", n); args = append(args, start); n++ }
	if !end.IsZero() { query += fmt.Sprintf(" AND created_at <= $%d", n); args = append(args, end); n++ }
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, size, offset)
	var total int64
	countQuery := "SELECT count(*) FROM audit_logs WHERE 1=1"
	countArgs := []interface{}{}
	cn := 1
	if orgID != nil { countQuery += fmt.Sprintf(" AND org_id = $%d", cn); countArgs = append(countArgs, *orgID); cn++ }
	if accountID != nil { countQuery += fmt.Sprintf(" AND account_id = $%d", cn); countArgs = append(countArgs, *accountID); cn++ }
	_ = r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil { return nil, 0, fmt.Errorf("query audit logs: %w", err) }
	defer rows.Close()
	var logs []model.AuditLog
	for rows.Next() { var l model.AuditLog; if err := rows.Scan(&l.ID,&l.EventID,&l.AccountID,&l.OrgID,&l.ActionDomain,&l.ActionType,&l.TargetType,&l.TargetID,&l.Result,&l.FailureReason,&l.LoginMethod,&l.Details,&l.TraceID,&l.IPAddress,&l.UserAgent,&l.CreatedAt); err != nil { return nil, 0, err }; logs = append(logs, l) }
	return logs, total, nil
}
func (r *AuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLog, error) {
	row := r.db.Pool.QueryRow(ctx, "SELECT id, event_id, account_id, org_id, action_domain, action_type, target_type, target_id, result, failure_reason, login_method, details, trace_id, ip_address, user_agent, created_at FROM audit_logs WHERE id = $1", id)
	var l model.AuditLog; err := row.Scan(&l.ID,&l.EventID,&l.AccountID,&l.OrgID,&l.ActionDomain,&l.ActionType,&l.TargetType,&l.TargetID,&l.Result,&l.FailureReason,&l.LoginMethod,&l.Details,&l.TraceID,&l.IPAddress,&l.UserAgent,&l.CreatedAt)
	if err != nil { return nil, err }
	return &l, nil
}
func (r *AuditRepository) ListByOrgID(ctx context.Context, orgID uuid.UUID, accountID *uuid.UUID, actionType, result string, start, end time.Time, page, size int) ([]model.AuditLog, int64, error) {
	return r.ListGlobal(ctx, &orgID, accountID, actionType, result, start, end, page, size)
}
