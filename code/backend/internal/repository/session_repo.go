package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"xyfamily/internal/model"
)

type SessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *model.Session) error {
	query := "INSERT INTO sessions (account_id, refresh_token_hash, device, client_ip, user_agent, expires_at, last_activity_at) VALUES ($1, $2, $3, $4, $5, $6, $6) RETURNING id"
	err := r.db.Pool.QueryRow(ctx, query,
		s.AccountID, s.RefreshTokenHash, s.Device, s.ClientIP, s.UserAgent, s.ExpiresAt,
	).Scan(&s.ID)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.Session, error) {
	query := "SELECT id, account_id, refresh_token_hash, device, client_ip, user_agent, expires_at, revoked_at, last_activity_at, created_at FROM sessions WHERE refresh_token_hash = $1"
	row := r.db.Pool.QueryRow(ctx, query, tokenHash)
	var s model.Session
	err := row.Scan(&s.ID, &s.AccountID, &s.RefreshTokenHash, &s.Device, &s.ClientIP, &s.UserAgent, &s.ExpiresAt, &s.RevokedAt, &s.LastActivityAt, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan session: %w", err)
	}
	return &s, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE sessions SET revoked_at = $1 WHERE id = $2 AND revoked_at IS NULL", time.Now(), id)
	return err
}

func (r *SessionRepository) RevokeAllByAccountID(ctx context.Context, accountID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE sessions SET revoked_at = $1 WHERE account_id = $2 AND revoked_at IS NULL", time.Now(), accountID)
	return err
}

func (r *SessionRepository) RevokeAllExceptCurrent(ctx context.Context, accountID uuid.UUID, currentSessionID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE sessions SET revoked_at = $1 WHERE account_id = $2 AND id != $3 AND revoked_at IS NULL", time.Now(), accountID, currentSessionID)
	return err
}
