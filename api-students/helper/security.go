package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Precomputed bcrypt cost 12 dummy hash for constant-time verification when user is not found.
// Hash of "dummy-password-for-constant-time-check"
const dummyBcryptHash = "$2a$12$6/7/Yv6q.X2X4zY7U3fF/.bV/L/H.nS6eD9B8a.x0U9v2N7g3q1yG"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashed), nil
}

// VerifyPassword membandingkan password dengan hash-nya.
func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// VerifyDummyPassword performs a dummy bcrypt check to prevent timing attacks / user enumeration.
func VerifyDummyPassword() {
	_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte("dummy-password-attempt"))
}

func RandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func SHA256Hex(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
