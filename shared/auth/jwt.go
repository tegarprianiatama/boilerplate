package auth

import (
	"go-boilerplate/shared/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a bcrypt hash with a plain-text password.
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// GenerateToken creates a signed JWT and returns it with its expiry time.
func GenerateToken(userID uint, email, role string) (string, time.Time, error) {
	cfg := config.GetConfig()
	now := time.Now()
	expiresAt := now.Add(time.Duration(cfg.JWTExpireHours) * time.Hour)

	claims := jwt.MapClaims{
		"sub":     strconv.FormatUint(uint64(userID), 10),
		"user_id": userID,
		"email":   email,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
