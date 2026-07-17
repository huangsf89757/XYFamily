package model

import (
	"time"
	"github.com/google/uuid"
)
type Role struct{ID uuid.UUID `json:"-" db:"id"`;RoleKey string `json:"role_key" db:"role_key"`;Name string `json:"name" db:"name"`;Level int `json:"level" db:"level"`;Scope string `json:"scope" db:"scope"`;Description string `json:"description" db:"description"`;CreatedAt time.Time `json:"created_at" db:"created_at"`}
type PermissionPoint struct{ID uuid.UUID `json:"-" db:"id"`;PermissionKey string `json:"permission_key" db:"permission_key"`;Module string `json:"module" db:"module"`;Description string `json:"description" db:"description"`;IsPublic bool `json:"is_public" db:"is_public"`;CreatedAt time.Time `json:"created_at" db:"created_at"`}
type RolePermission struct{ID uuid.UUID `json:"-" db:"id"`;RoleKey string `json:"role_key" db:"role_key"`;PermissionKey string `json:"permission_key" db:"permission_key"`;CreatedAt time.Time `json:"created_at" db:"created_at"`}
type SystemConfig struct{ID uuid.UUID `json:"-" db:"id"`;ConfigKey string `json:"config_key" db:"config_key"`;ConfigValue string `json:"config_value" db:"config_value"`;Description string `json:"description" db:"description"`;UpdatedAt time.Time `json:"updated_at" db:"updated_at"`;UpdatedBy *uuid.UUID `json:"updated_by,omitempty" db:"updated_by"`}
