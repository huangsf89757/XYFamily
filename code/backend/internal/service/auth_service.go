package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"xyfamily/internal/model"
	"xyfamily/internal/repository"
	"xyfamily/pkg/bcrypt"
	appErr "xyfamily/pkg/errors"
	"xyfamily/pkg/jwt"
	"xyfamily/pkg/logger"
	"xyfamily/pkg/config"
	"go.uber.org/zap"
)

type AuthService struct {
	accountRepo *repository.AccountRepository
	sessionRepo *repository.SessionRepository
	cacheRepo   *repository.CacheRepository
	rbacRepo    *repository.RBACRepository
	auditRepo   *repository.AuditRepository
	jwtMgr      *jwt.Manager
	cfg         *config.Config
}

func NewAuthService(
	accountRepo *repository.AccountRepository,
	sessionRepo *repository.SessionRepository,
	cacheRepo *repository.CacheRepository,
	rbacRepo *repository.RBACRepository,
	auditRepo *repository.AuditRepository,
	jwtMgr *jwt.Manager,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		accountRepo: accountRepo, sessionRepo: sessionRepo, cacheRepo: cacheRepo, rbacRepo: rbacRepo, auditRepo: auditRepo, jwtMgr: jwtMgr, cfg: cfg}
}

type SendCodeReq struct {
	Target string `json:"target"`
	Type   string `json:"type"`
}

type SendCodeResp struct {
	ExpiresIn int `json:"expires_in"`
}

func (s *AuthService) SendVerificationCode(ctx context.Context, req *SendCodeReq) (*SendCodeResp, error) {
	validTypes := map[string]bool{"register": true, "login": true, "reset_password": true}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("%w: invalid type", appErr.ErrBadRequest)
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := s.cacheRepo.SetVerificationCode(ctx, req.Target, req.Type, code, s.cfg.Security.VerificationCodeTTL); err != nil {
		logger.Get().Error("set verification code failed", zap.Error(err))
		return nil, fmt.Errorf("%w: internal error", appErr.ErrInternal)
	}
	logger.Get().Info("verification code sent", zap.String("target", req.Target), zap.String("type", req.Type), zap.String("code_dev", code))
	return &SendCodeResp{ExpiresIn: s.cfg.Security.VerificationCodeTTL}, nil
}

func (s *AuthService) verifyCode(ctx context.Context, target, codeType, code string) error {
	stored, err := s.cacheRepo.GetVerificationCode(ctx, target, codeType)
	if err != nil {
		return fmt.Errorf("%w: internal", appErr.ErrInternal)
	}
	if stored == "" {
		return fmt.Errorf("%w: code expired or not found", appErr.ErrCodeExpired)
	}
	if stored != code {
		return fmt.Errorf("%w: wrong code", appErr.ErrCodeWrong)
	}
	_ = s.cacheRepo.DelVerificationCode(ctx, target, codeType)
	return nil
}
type RegisterReq struct {
	Type     string `json:"type"`
	Target   string `json:"target"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type AuthResp struct {
	AccountID     string `json:"account_id"`
	Nickname      string `json:"nickname"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	Status        string `json:"status,omitempty"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterReq) (*AuthResp, error) {
	if err := s.verifyCode(ctx, req.Target, "register", req.Code); err != nil {
		return nil, err
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("%w: password too short", appErr.ErrPasswordWeak)
	}
	var phoneHash, emailHash string
	if req.Type == "phone" {
		phoneHash = hashPII(req.Target)
		existing, _ := s.accountRepo.GetByPhoneHash(ctx, phoneHash)
		if existing != nil {
			return nil, fmt.Errorf("%w", appErr.ErrPhoneRegistered)
		}
	} else {
		emailHash = hashPII(req.Target)
		existing, _ := s.accountRepo.GetByEmailHash(ctx, emailHash)
		if existing != nil {
			return nil, fmt.Errorf("%w", appErr.ErrEmailRegistered)
		}
	}
	passHash, err := bcrypt.Hash(req.Password, s.cfg.Security.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("%w: hash password", appErr.ErrInternal)
	}
	nick := fmt.Sprintf("user_%s", req.Target[len(req.Target)-6:])
	account := &model.Account{
		AccountID:    uuid.New(),
		PhoneHash:    phoneHash,
		EmailHash:    emailHash,
		PasswordHash: passHash,
		Nickname:     nick,
		Status:       model.AccountStatusActive,
	}
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("%w: create account", appErr.ErrInternal)
	}
	tokenPair, refreshToken, err := s.generateTokenPair(ctx, account)
	if err != nil {
		return nil, err
	}
	s.writeAuditLog(ctx, account.AccountID, "login", "auth.register", "success", "", "")
	return &AuthResp{
		AccountID: account.AccountID.String(), Nickname: nick,
		AccessToken: tokenPair.AccessToken, RefreshToken: refreshToken,
		ExpiresIn: s.cfg.JWT.AccessTTL, Status: account.Status,
	}, nil
}
type LoginReq struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifier_type"`
	CredentialType string `json:"credential_type"`
	Password       string `json:"password"`
	Code           string `json:"code"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginReq, clientIP, userAgent string) (*AuthResp, error) {
	var account *model.Account
	var err error
	switch req.IdentifierType {
	case "phone":
		account, err = s.accountRepo.GetByPhoneHash(ctx, hashPII(req.Identifier))
	case "email":
		account, err = s.accountRepo.GetByEmailHash(ctx, hashPII(req.Identifier))
	case "username":
		account, err = s.accountRepo.GetByUsername(ctx, req.Identifier)
	default:
		return nil, fmt.Errorf("%w: invalid identifier_type", appErr.ErrBadRequest)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: internal", appErr.ErrInternal)
	}
	if account == nil {
		s.writeAuditLog(ctx, uuid.Nil, "login", "auth.login", "failed", "account_not_found", "")
		return nil, fmt.Errorf("%w", appErr.ErrAccountNotFound)
	}
	if account.Status == model.AccountStatusDeactivated {
		return nil, fmt.Errorf("%w", appErr.ErrAccountDeactivated)
	}
	if req.CredentialType == "password" {
		if err := bcrypt.Compare(account.PasswordHash, req.Password); err != nil {
			s.writeAuditLog(ctx, account.AccountID, "login", "auth.login", "failed", "password_error", "password")
			return nil, fmt.Errorf("%w", appErr.ErrPasswordWrong)
		}
	} else if req.CredentialType == "code" {
		if err := s.verifyCode(ctx, req.Identifier, "login", req.Code); err != nil {
			return nil, err
		}
	}
	_ = s.accountRepo.UpdateLastLogin(ctx, account.AccountID)
	tokenPair, refreshToken, err := s.generateTokenPair(ctx, account)
	if err != nil {
		return nil, err
	}
	s.createSession(ctx, account.AccountID, refreshToken, clientIP, userAgent)
	s.writeAuditLog(ctx, account.AccountID, "login", "auth.login", "success", "", "password")
	return &AuthResp{
		AccountID: account.AccountID.String(), Nickname: account.Nickname,
		AccessToken: tokenPair.AccessToken, RefreshToken: refreshToken,
		ExpiresIn: s.cfg.JWT.AccessTTL, Status: account.Status,
	}, nil
}
type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
}

func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshReq, clientIP, userAgent string) (*RefreshResp, error) {
	session, err := s.sessionRepo.GetByTokenHash(ctx, hashToken(req.RefreshToken))
	if err != nil {
		return nil, fmt.Errorf("%w: internal", appErr.ErrInternal)
	}
	if session == nil {
		return nil, fmt.Errorf("%w: invalid refresh token", appErr.ErrRefreshInvalid)
	}
	if session.RevokedAt != nil {
		_ = s.sessionRepo.RevokeAllByAccountID(ctx, session.AccountID)
		return nil, fmt.Errorf("%w: refresh token replay detected", appErr.ErrRefreshReplayed)
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("%w: refresh token expired", appErr.ErrRefreshExpired)
	}
	_ = s.sessionRepo.Revoke(ctx, session.ID)
	account, err := s.accountRepo.GetByAccountID(ctx, session.AccountID)
	if err != nil || account == nil {
		return nil, fmt.Errorf("%w: account not found", appErr.ErrAccountNotFound)
	}
	tokenPair, newRefreshToken, err := s.generateTokenPair(ctx, account)
	if err != nil {
		return nil, err
	}
	s.createSession(ctx, account.AccountID, newRefreshToken, clientIP, userAgent)
	return &RefreshResp{AccessToken: tokenPair.AccessToken, RefreshToken: newRefreshToken, ExpiresIn: s.cfg.JWT.AccessTTL}, nil
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Logout(ctx context.Context, req *LogoutReq, accountID, jti string) error {
	if req.RefreshToken != "" {
		session, _ := s.sessionRepo.GetByTokenHash(ctx, hashToken(req.RefreshToken))
		if session != nil {
			_ = s.sessionRepo.Revoke(ctx, session.ID)
		}
	}
	if jti != "" {
		_ = s.cacheRepo.BlacklistToken(ctx, jti, s.cfg.JWT.AccessTTL)
	}
	s.writeAuditLog(ctx, uuid.MustParse(accountID), "login", "auth.logout", "success", "", "")
	return nil
}

type ResetPasswordReq struct {
	Type        string `json:"type"`
	Target      string `json:"target"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordReq) error {
	if err := s.verifyCode(ctx, req.Target, "reset_password", req.Code); err != nil {
		return err
	}
	if len(req.NewPassword) < 8 {
		return fmt.Errorf("%w: password too short", appErr.ErrPasswordWeak)
	}
	var account *model.Account
	if req.Type == "phone" {
		account, _ = s.accountRepo.GetByPhoneHash(ctx, hashPII(req.Target))
	} else {
		account, _ = s.accountRepo.GetByEmailHash(ctx, hashPII(req.Target))
	}
	if account == nil {
		return nil
	}
	passHash, err := bcrypt.Hash(req.NewPassword, s.cfg.Security.BcryptCost)
	if err != nil {
		return fmt.Errorf("%w: internal", appErr.ErrInternal)
	}
	if err := s.accountRepo.UpdatePassword(ctx, account.AccountID, passHash); err != nil {
		return fmt.Errorf("%w: internal", appErr.ErrInternal)
	}
	_ = s.sessionRepo.RevokeAllByAccountID(ctx, account.AccountID)
	s.writeAuditLog(ctx, account.AccountID, "account", "account.password.reset", "success", "", "")
	return nil
}
func (s *AuthService) generateTokenPair(ctx context.Context, account *model.Account) (*jwt.TokenPair, string, error) {
	orgIDs, _ := s.rbacRepo.GetOrgIDsByAccountID(ctx, account.AccountID)
	roles := make(map[string]string)
	for _, oid := range orgIDs {
		role, _ := s.rbacRepo.GetOrgRoleByAccountID(ctx, account.AccountID, oid)
		if role == "" {
			role = "regular_member"
		}
		roles[oid.String()] = role
	}
	orgIDStrs := make([]string, len(orgIDs))
	for i, oid := range orgIDs {
		orgIDStrs[i] = oid.String()
	}
	accessToken, _, err := s.jwtMgr.GenerateAccessToken(account.AccountID.String(), orgIDStrs, roles)
	if err != nil {
		return nil, "", err
	}
	refreshToken, tokenHash, _ := s.jwtMgr.GenerateRefreshToken()
	return &jwt.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: s.cfg.JWT.AccessTTL}, tokenHash, nil
}

func (s *AuthService) createSession(ctx context.Context, accountID uuid.UUID, refreshToken, clientIP, userAgent string) {
	tokenHash := hashToken(refreshToken)
	session := &model.Session{
		AccountID:        accountID,
		RefreshTokenHash: tokenHash,
		ClientIP:         clientIP,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(s.jwtMgr.RefreshTTL()),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.Get().Error("create session failed", zap.Error(err))
	}
}

func (s *AuthService) writeAuditLog(ctx context.Context, accountID uuid.UUID, domain, actionType, result, failureReason, loginMethod string) {
	log := &model.AuditLog{
		EventID:       uuid.New(),
		AccountID:     &accountID,
		ActionDomain:  domain,
		ActionType:    actionType,
		Result:        result,
		FailureReason: failureReason,
		LoginMethod:   loginMethod,
	}
	if accountID == uuid.Nil {
		log.AccountID = nil
	}
	if err := s.auditRepo.Create(ctx, log); err != nil {
		logger.Get().Error("write audit log failed", zap.Error(err))
	}
}

func hashPII(input string) string {
	h := fmt.Sprintf("%x", []byte(input))
	return h
}

func hashToken(token string) string {
	return fmt.Sprintf("%x", []byte(token))
}
