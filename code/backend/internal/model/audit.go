package model
import (
	"time"
	"github.com/google/uuid"
)
type AuditLog struct{ID uuid.UUID `json:"-" db:"id"`;EventID uuid.UUID `json:"event_id" db:"event_id"`;AccountID *uuid.UUID `json:"account_id,omitempty" db:"account_id"`;OrgID *uuid.UUID `json:"org_id,omitempty" db:"org_id"`;ActionDomain string `json:"action_domain" db:"action_domain"`;ActionType string `json:"action_type" db:"action_type"`;TargetType string `json:"target_type,omitempty" db:"target_type"`;TargetID *uuid.UUID `json:"target_id,omitempty" db:"target_id"`;Result string `json:"result,omitempty" db:"result"`;FailureReason string `json:"failure_reason,omitempty" db:"failure_reason"`;LoginMethod string `json:"login_method,omitempty" db:"login_method"`;Details []byte `json:"details,omitempty" db:"details"`;TraceID string `json:"trace_id,omitempty" db:"trace_id"`;IPAddress string `json:"ip_address,omitempty" db:"ip_address"`;UserAgent string `json:"user_agent,omitempty" db:"user_agent"`;CreatedAt time.Time `json:"created_at" db:"created_at"`}
