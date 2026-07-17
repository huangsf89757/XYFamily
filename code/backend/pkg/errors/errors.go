package errors

import "fmt"

// ErrCode is a 5-digit business error code (ADR-013).
type ErrCode int

const (
	// Auth errors 10xxx
	ErrTokenMissing     ErrCode = 101001
	ErrTokenInvalid     ErrCode = 101002
	ErrTokenRevoked     ErrCode = 101003
	ErrRefreshInvalid   ErrCode = 101004
	ErrRefreshRevoked   ErrCode = 101005
	ErrRefreshExpired   ErrCode = 101006
	ErrAccountNotFound  ErrCode = 101007
	ErrPasswordWrong    ErrCode = 101008
	ErrAccountDeactivated ErrCode = 101009
	ErrRefreshReplayed  ErrCode = 101010
	ErrLoginLocked       ErrCode = 104290
	ErrRefreshRateLimit ErrCode = 104291
	ErrResetRateLimit   ErrCode = 104292

	// Verification code errors 110xxx
	ErrCodeWrong   ErrCode = 110001
	ErrCodeExpired ErrCode = 110002
	ErrCodeUsed     ErrCode = 110003
	ErrCodeMaxAttempt ErrCode = 110004

	// Verification code rate limit 114xxx
	ErrCodeTargetRate ErrCode = 114290
	ErrCodeIPRate     ErrCode = 114291
	ErrCodeAccountRate ErrCode = 114293

	// Account errors 20xxx
	ErrPasswordWeak       ErrCode = 200001
	ErrPhoneRegistered    ErrCode = 200002
	ErrEmailRegistered    ErrCode = 200003
	ErrPasswordSame       ErrCode = 200005
	ErrDeactivating       ErrCode = 200007
	ErrOldPasswordWrong   ErrCode = 200011

	// General 80xxx
	ErrBadRequest  ErrCode = 100001
	ErrInternal    ErrCode = 800001
	ErrNotFound     ErrCode = 800002
	ErrForbidden    ErrCode = 800003
)

func (e ErrCode) Int() int { return int(e) }

func (e ErrCode) Error() string { return fmt.Sprintf("error_code_%d", e) }
