package errors

import (
	"testing"
)

func TestErrCodeValues(t *testing.T) {
	tests := []struct {
		name string
		code  ErrCode
		expected int
	}{
		{"token missing", ErrTokenMissing, 101001},
		{"token invalid", ErrTokenInvalid, 101002},
		{"token revoked", ErrTokenRevoked, 101003},
		{"account not found", ErrAccountNotFound, 101007},
		{"password wrong", ErrPasswordWrong, 101008},
		{"account deactivated", ErrAccountDeactivated, 101009},
		{"refresh replayed", ErrRefreshReplayed, 101010},
		{"login locked", ErrLoginLocked, 104290},
		{"code wrong", ErrCodeWrong, 110001},
		{"code expired", ErrCodeExpired, 110002},
		{"password weak", ErrPasswordWeak, 200001},
		{"phone registered", ErrPhoneRegistered, 200002},
		{"email registered", ErrEmailRegistered, 200003},
		{"bad request", ErrBadRequest, 100001},
		{"internal", ErrInternal, 800001},
		{"not found", ErrNotFound, 800002},
		{"forbidden", ErrForbidden, 800003},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code.Int() != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, tt.code.Int())
			}
		})
	}
}

func TestErrCodeError(t *testing.T) {
	e := ErrTokenInvalid
	s := e.Error()
	if s == "" {
		t.Fatal("Error() should not return empty string")
	}
}
