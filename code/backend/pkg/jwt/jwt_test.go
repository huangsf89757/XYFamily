package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	m := NewManager("test-secret-key", "test-issuer", 1800, 604800)
	orgIDs := []string{"org-1", "org-2"}
	roles := map[string]string{"org-1": "organization_core_admin"}

	token, jti, err := m.GenerateAccessToken("account-1", orgIDs, roles)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}
	if jti == "" {
		t.Fatal("jti should not be empty")
	}

	claims, err := m.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken failed: %v", err)
	}
	if claims.AccountID != "account-1" {
		t.Errorf("expected account_id=account-1, got %s", claims.AccountID)
	}
	if claims.JTI != jti {
		t.Errorf("jti mismatch: expected %s, got %s", jti, claims.JTI)
	}
	if len(claims.OrgIDs) != 2 {
		t.Errorf("expected 2 org_ids, got %d", len(claims.OrgIDs))
	}
	if claims.Roles["org-1"] != "organization_core_admin" {
		t.Errorf("unexpected role: %s", claims.Roles["org-1"])
	}
}

func TestParseInvalidToken(t *testing.T) {
	m := NewManager("test-secret-key", "test-issuer", 1800, 604800)
	_, err := m.ParseAccessToken("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestParseTokenWithWrongSecret(t *testing.T) {
	m1 := NewManager("secret-1", "test", 1800, 604800)
	m2 := NewManager("secret-2", "test", 1800, 604800)
	token, _, _ := m1.GenerateAccessToken("acct-1", nil, nil)
	_, err := m2.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected error with wrong secret")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	m := NewManager("secret", "issuer", 1800, 604800)
	token, hash, err := m.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if token == "" || hash == "" {
		t.Fatal("token and hash should not be empty")
	}
	if token == hash {
		t.Fatal("token and hash should be different")
	}
}

func TestAccessTTL(t *testing.T) {
	m := NewManager("secret", "issuer", 1800, 604800)
	if m.AccessTTL() != 1800*time.Second {
		t.Errorf("expected 1800s, got %v", m.AccessTTL())
	}
	if m.RefreshTTL() != 604800*time.Second {
		t.Errorf("expected 604800s, got %v", m.RefreshTTL())
	}
}
