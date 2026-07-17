package model

import (
	"time"
	"github.com/google/uuid"
)

type Organization struct {
	ID uuid.UUID `json:"-" db:"id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at,omitempty" db:"disabled_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type OrgMember struct {
	ID uuid.UUID `json:"-" db:"id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	Role string `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
type Team struct {
	ID uuid.UUID `json:"-" db:"id"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type TeamMember struct {
	ID uuid.UUID `json:"-" db:"id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	Role string `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type Group struct {
	ID uuid.UUID `json:"-" db:"id"`
	GroupID uuid.UUID `json:"group_id" db:"group_id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type GroupMember struct {
	ID uuid.UUID `json:"-" db:"id"`
	OrgID uuid.UUID `json:"org_id" db:"org_id"`
	TeamID uuid.UUID `json:"team_id" db:"team_id"`
	GroupID uuid.UUID `json:"group_id" db:"group_id"`
	AccountID uuid.UUID `json:"account_id" db:"account_id"`
	Role string `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
