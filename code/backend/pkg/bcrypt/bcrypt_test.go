package bcrypt

import (
	"testing"
)

func TestHashAndCompare(t *testing.T) {
	password := "MyStrong@Pass123"
	hash, err := Hash(password, 10)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == password {
		t.Fatal("hash should differ from password")
	}
	err = Compare(hash, password)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
}

func TestCompareWrongPassword(t *testing.T) {
	hash, _ := Hash("correct-password", 10)
	err := Compare(hash, "wrong-password")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestHashUniqueness(t *testing.T) {
	h1, _ := Hash("same-password", 10)
	h2, _ := Hash("same-password", 10)
	if h1 == h2 {
		t.Fatal("hashes should be different due to random salt")
	}
}
