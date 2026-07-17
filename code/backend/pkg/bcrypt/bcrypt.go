package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash returns a bcrypt hash of the password with the given cost.
// Default cost is 12 (NFR-SEC-001).
func Hash(password string, cost int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on match, non-nil error otherwise.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
