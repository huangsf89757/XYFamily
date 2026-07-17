package model

import (
	"time"
	"github.com/google/uuid"
)

const (
	AccountStatusActive = "active"
	AccountStatusDeactivating = "deactivating"
	AccountStatusDeactivated = "deactivated"
)

type Account struct {
	ID uuid.UUID `json:"-" db:"id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	PhoneEncrypted []byte `json:"-" db:"phone_encrypted"`
	PhoneHash string `json:"-" db:"phone_hash"`
	EmailEncrypted []byte `json:"-" db:"email_encrypted"`
	EmailHash string `json:"-" db:"email_hash"`
	Username string `json:"username" db:"username"`
	PasswordHash string `json:"-" db:"password_hash"`
	Nickname string `json:"nickname" db:"nickname"`
	Avatar string `json:"avatar" db:"avatar"`
	Status string `json:"status" db:"status"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty" db:"deactivated_at"`
	PreviousUsername *string `json:"-" db:"previous_username"`
	UsernameChangedAt *time.Time `json:"-" db:"username_changed_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type Session struct {
	ID uuid.UUID `json:"-" db:"id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	RefreshTokenHash string `json:"-" db:"refresh_token_hash"`
	Device string `json:"device,omitempty" db:"device"`
	ClientIP string `json:"-" db:"client_ip"`
	UserAgent string `json:"-" db:"user_agent"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	LastActivityAt *time.Time `json:"-" db:"last_activity_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Invitation struct {
	ID uuid.UUID `json:"id" db:"id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	ScopeType string `json:"scope_type" db:"scope_type"`
	ScopeID uuid.UUID `json:"scope_id" db:"scope_id"`
	InviterID uuid.UUID `json:"inviter_id" db:"inviter_id"`
	InviteeAccountID *uuid.UUID `json:"invitee_account_id,omitempty" db:"invitee_account_id"`
	InviteeContact string `json:"invitee_contact" db:"invitee_contact"`
	Role string `json:"role" db:"role"`
	Token string `json:"-" db:"token"`
	Status string `json:"status" db:"status"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty" db:"accepted_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
