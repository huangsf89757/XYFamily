package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims defines the JWT Access Token payload (ADR-003, P2 JWT).
type Claims struct {
	AccountID string `json:"sub"`
	JTI       string `json:"jti"`
	OrgIDs    []string `json:"org_ids,omitempty"`
	Roles     map[string]string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair holds the access and refresh tokens returned to the client.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Manager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(secret, issuer string, accessTTL, refreshTTL int) *Manager {
	return &Manager{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  time.Duration(accessTTL) * time.Second,
		refreshTTL: time.Duration(refreshTTL) * time.Second,
	}
}

func (m *Manager) GenerateAccessToken(accountID string, orgIDs []string, roles map[string]string) (string, string, error) {
	now := time.Now()
	jti := uuid.New().String()
	claims := &Claims{
		AccountID: accountID,
		JTI:       jti,
		OrgIDs:    orgIDs,
		Roles:     roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, jti, nil
}

func (m *Manager) GenerateRefreshToken() (string, string, error) {
	token := uuid.New().String()
	hash := uuid.New().String()
	return token, hash, nil
}

func (m *Manager) ParseAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *Manager) AccessTTL() time.Duration { return m.accessTTL }
func (m *Manager) RefreshTTL() time.Duration { return m.refreshTTL }
