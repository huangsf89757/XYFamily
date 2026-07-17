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

type AccountRepository struct {
	db *DB
}

func NewAccountRepository(db *DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, a *model.Account) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO accounts (account_id, phone_encrypted, phone_hash, email_encrypted, email_hash, username, password_hash, nickname, avatar, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		a.AccountID, a.PhoneEncrypted, a.PhoneHash, a.EmailEncrypted, a.EmailHash, a.Username, a.PasswordHash, a.Nickname, a.Avatar, a.Status)
	return err
}

func (r *AccountRepository) GetByPhoneHash(ctx context.Context, phoneHash string) (*model.Account, error) {
	return r.getByField(ctx, "phone_hash", phoneHash)
}

func (r *AccountRepository) GetByEmailHash(ctx context.Context, emailHash string) (*model.Account, error) {
	return r.getByField(ctx, "email_hash", emailHash)
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*model.Account, error) {
	return r.getByField(ctx, "username", username)
}

func (r *AccountRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.Account, error) {
	return r.getByField(ctx, "account_id", accountID)
}

func (r *AccountRepository) getByField(ctx context.Context, field string, val interface{}) (*model.Account, error) {
	query := fmt.Sprintf("SELECT id, account_id, phone_encrypted, phone_hash, email_encrypted, email_hash, username, password_hash, nickname, avatar, status, deactivated_at, last_login_at, created_at, updated_at FROM accounts WHERE %s = $1 AND deleted_at IS NULL", field)
	row := r.db.Pool.QueryRow(ctx, query, val)
	var a model.Account
	err := row.Scan(&a.ID, &a.AccountID, &a.PhoneEncrypted, &a.PhoneHash, &a.EmailEncrypted, &a.EmailHash, &a.Username, &a.PasswordHash, &a.Nickname, &a.Avatar, &a.Status, &a.DeactivatedAt, &a.LastLoginAt, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
	if err != nil { return nil, fmt.Errorf("scan account: %w", err) }
	return &a, nil
}

func (r *AccountRepository) UpdateLastLogin(ctx context.Context, accountID uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE accounts SET last_login_at = $1, updated_at = $1 WHERE account_id = $2 AND deleted_at IS NULL", time.Now(), accountID)
	return err
}

func (r *AccountRepository) UpdatePassword(ctx context.Context, accountID uuid.UUID, passwordHash string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE accounts SET password_hash = $1, updated_at = now() WHERE account_id = $2 AND deleted_at IS NULL", passwordHash, accountID)
	return err
}

func (r *AccountRepository) UpdateStatus(ctx context.Context, accountID uuid.UUID, status string, deactivatedAt *time.Time) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE accounts SET status = $1, deactivated_at = $2, updated_at = now() WHERE account_id = $3 AND deleted_at IS NULL", status, deactivatedAt, accountID)
	return err
}

func (r *AccountRepository) UpdateProfile(ctx context.Context, accountID uuid.UUID, nickname, avatar string) error {
	_, err := r.db.Pool.Exec(ctx, "UPDATE accounts SET nickname = $1, avatar = $2, updated_at = now() WHERE account_id = $3 AND deleted_at IS NULL", nickname, avatar, accountID)
	return err
}
