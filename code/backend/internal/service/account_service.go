package service

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"xyfamily/internal/model"
	"xyfamily/internal/repository"
	"xyfamily/pkg/bcrypt"
	appErr "xyfamily/pkg/errors"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
	sessionRepo *repository.SessionRepository
	cacheRepo   *repository.CacheRepository
	auditRepo   *repository.AuditRepository
}

func NewAccountService(ar *repository.AccountRepository, sr *repository.SessionRepository, cr *repository.CacheRepository, aur *repository.AuditRepository) *AccountService {
	return &AccountService{accountRepo: ar, sessionRepo: sr, cacheRepo: cr, auditRepo: aur}
}

type ProfileResp struct {
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
}

func (s *AccountService) GetProfile(ctx context.Context, accountID uuid.UUID) (*ProfileResp, error) {
	account, err := s.accountRepo.GetByAccountID(ctx, accountID)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	if account == nil { return nil, fmt.Errorf("%w", appErr.ErrAccountNotFound) }
	phone := ""; if account.PhoneHash != "" { phone = maskPII(account.PhoneHash) }
	email := ""; if account.EmailHash != "" { email = maskPII(account.EmailHash) }
	return &ProfileResp{AccountID: account.AccountID.String(), Username: account.Username, Phone: phone, Email: email, Nickname: account.Nickname, Avatar: account.Avatar, Status: account.Status}, nil
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (s *AccountService) UpdateProfile(ctx context.Context, accountID uuid.UUID, req *UpdateProfileReq) error {
	return s.accountRepo.UpdateProfile(ctx, accountID, req.Nickname, req.Avatar)
}
type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (s *AccountService) ChangePassword(ctx context.Context, accountID uuid.UUID, req *ChangePasswordReq) error {
	account, err := s.accountRepo.GetByAccountID(ctx, accountID)
	if err != nil || account == nil { return fmt.Errorf("%w", appErr.ErrAccountNotFound) }
	if err := bcrypt.Compare(account.PasswordHash, req.OldPassword); err != nil { return fmt.Errorf("%w", appErr.ErrOldPasswordWrong) }
	if req.NewPassword == req.OldPassword { return fmt.Errorf("%w", appErr.ErrPasswordSame) }
	if len(req.NewPassword) < 8 { return fmt.Errorf("%w: too short", appErr.ErrPasswordWeak) }
	hash, err := bcrypt.Hash(req.NewPassword, 12)
	if err != nil { return fmt.Errorf("%w: internal", appErr.ErrInternal) }
	if err := s.accountRepo.UpdatePassword(ctx, accountID, hash); err != nil { return fmt.Errorf("%w: internal", appErr.ErrInternal) }
	_ = s.sessionRepo.RevokeAllByAccountID(ctx, accountID)
	return nil
}

func (s *AccountService) Deactivate(ctx context.Context, accountID uuid.UUID) (string, time.Time, error) {
	account, err := s.accountRepo.GetByAccountID(ctx, accountID)
	if err != nil || account == nil { return "", time.Time{}, fmt.Errorf("%w", appErr.ErrAccountNotFound) }
	if account.Status != model.AccountStatusActive { return "", time.Time{}, fmt.Errorf("%w: already deactivating", appErr.ErrDeactivating) }
	deactivatedAt := time.Now().Add(30 * 24 * time.Hour)
	if err := s.accountRepo.UpdateStatus(ctx, accountID, model.AccountStatusDeactivating, &deactivatedAt); err != nil { return "", time.Time{}, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	_ = s.sessionRepo.RevokeAllByAccountID(ctx, accountID)
	return model.AccountStatusDeactivating, deactivatedAt, nil
}

type UndeactivateReq struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Code   string `json:"code"`
}

func (s *AccountService) Undeactivate(ctx context.Context, accountID uuid.UUID, req *UndeactivateReq) error {
	account, err := s.accountRepo.GetByAccountID(ctx, accountID)
	if err != nil || account == nil { return fmt.Errorf("%w", appErr.ErrAccountNotFound) }
	if account.Status != model.AccountStatusDeactivating { return fmt.Errorf("%w: not in grace period", appErr.ErrDeactivating) }
	stored, _ := s.cacheRepo.GetVerificationCode(ctx, req.Target, "reset_password")
	if stored == "" || stored != req.Code { return fmt.Errorf("%w: code invalid", appErr.ErrCodeWrong) }
	_ = s.cacheRepo.DelVerificationCode(ctx, req.Target, "reset_password")
	if err := s.accountRepo.UpdateStatus(ctx, accountID, model.AccountStatusActive, nil); err != nil { return fmt.Errorf("%w: internal", appErr.ErrInternal) }
	return nil
}

func maskPII(s string) string { if len(s) <= 4 { return "****" }; return s[:3] + "****" + s[len(s)-4:] }
