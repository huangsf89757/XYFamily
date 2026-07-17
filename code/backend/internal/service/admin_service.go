package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"xyfamily/internal/model"
	"xyfamily/internal/repository"
	"xyfamily/pkg/bcrypt"
	appErr "xyfamily/pkg/errors"
)

type AdminService struct {
	accountRepo *repository.AccountRepository
	rbacRepo    *repository.RBACRepository
	tenantRepo  *repository.TenantRepository
	auditRepo   *repository.AuditRepository
	cacheRepo   *repository.CacheRepository
	pool        interface{ Exec(context.Context, string, ...interface{}) (interface{}, error) }
}

func NewAdminService(ar *repository.AccountRepository, rr *repository.RBACRepository, tr *repository.TenantRepository, aur *repository.AuditRepository, cr *repository.CacheRepository) *AdminService {
	return &AdminService{accountRepo: ar, rbacRepo: rr, tenantRepo: tr, auditRepo: aur, cacheRepo: cr}
}

type InitReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type InitResp struct {
	AccountID            string `json:"account_id"`
	RequirePasswordChange bool  `json:"require_password_change"`
}

func (s *AdminService) Init(ctx context.Context, req *InitReq) (*InitResp, error) {
	count, _ := s.rbacRepo.CountSuperAdmins(ctx)
	if count > 0 { return nil, fmt.Errorf("%w: super admin exists", appErr.ErrSuperAdminExists) }
	hash, err := bcrypt.Hash(req.Password, 12)
	if err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	acctID := uuid.New()
	emailHash := req.Email
	acct := &model.Account{AccountID: acctID, EmailHash: emailHash, PasswordHash: hash, Nickname: "SuperAdmin", Status: model.AccountStatusActive}
	if err := s.accountRepo.Create(ctx, acct); err != nil { return nil, fmt.Errorf("%w: internal", appErr.ErrInternal) }
	return &InitResp{AccountID: acctID.String(), RequirePasswordChange: true}, nil
}

type ConfigResp struct {
	TokenAccessTTL          int    `json:"token.access_ttl"`
	TokenRefreshTTL         int    `json:"token.refresh_ttl"`
	LoginRateLimitThreshold int    `json:"login.rate_limit.threshold"`
	VerificationCodeTTL     int    `json:"verification_code.ttl"`
	AuditRetentionDays      int    `json:"audit.retention_days"`
	DeactivateGraceDays     int    `json:"account.deactivate_grace_days"`
	LoginAutoRegister       bool   `json:"login.auto_register"`
}

func (s *AdminService) GetConfig(ctx context.Context) (*ConfigResp, error) {
	resp := &ConfigResp{TokenAccessTTL: 1800, TokenRefreshTTL: 604800, LoginRateLimitThreshold: 5, VerificationCodeTTL: 300, AuditRetentionDays: 365, DeactivateGraceDays: 30, LoginAutoRegister: false}
	return resp, nil
}

type UpdateConfigReq struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Confirm     bool   `json:"confirm"`
}

func (s *AdminService) UpdateConfig(ctx context.Context, req *UpdateConfigReq) error {
	if !req.Confirm { return fmt.Errorf("%w: confirm required", appErr.ErrConfigInvalid) }
	return fmt.Errorf("%w: not implemented", appErr.ErrInternal)
}

type ForceDowngradeReq struct {
	ScopeType string `json:"scope_type"`
	ScopeID   string `json:"scope_id"`
	AccountID string `json:"account_id"`
	Confirm   bool   `json:"confirm"`
}
type ForceDowngradeResp struct {
	AccountID  string `json:"account_id"`
	BeforeRole string `json:"before_role"`
	AfterRole  string `json:"after_role"`
}

func (s *AdminService) ForceDowngrade(ctx context.Context, req *ForceDowngradeReq) (*ForceDowngradeResp, error) {
	if !req.Confirm { return nil, fmt.Errorf("%w: confirm required", appErr.ErrConfigInvalid) }
	acctID, _ := uuid.Parse(req.AccountID)
	beforeRole := "organization_core_admin"
	if req.ScopeType == "organization" {
		orgID, _ := uuid.Parse(req.ScopeID)
		_ = s.tenantRepo.UpdateOrgMemberRole(ctx, orgID, acctID, "regular_member")
	}
	return &ForceDowngradeResp{AccountID: req.AccountID, BeforeRole: beforeRole, AfterRole: "regular_member"}, nil
}

func (s *AdminService) GlobalAuditList(ctx context.Context, orgID, accountID *uuid.UUID, actionType, result string, start, end time.Time, page, size int) ([]model.AuditLog, int64, error) {
	return s.auditRepo.ListGlobal(ctx, orgID, accountID, actionType, result, start, end, page, size)
}

func (s *AdminService) AuditDetail(ctx context.Context, id uuid.UUID) (*model.AuditLog, error) {
	return s.auditRepo.GetByID(ctx, id)
}

var _ = strconv.Itoa
