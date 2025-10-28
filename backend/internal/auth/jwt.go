package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultIssuer        = "aroma-sense-api"
	minSecretLength      = 32
	defaultExpiryMins    = 15
	clockSkewSeconds     = 5
	refreshTokenBytes    = 32 // 256 bits
	refreshTokenDuration = 7 * 24 * time.Hour // 7 days
)

var (
	secretOnce sync.Once
	secret     []byte
	secretErr  error
)

func loadSecret() ([]byte, error) {
	secretOnce.Do(func() {
		s := os.Getenv("JWT_SECRET")
		if s == "" {
			secretErr = errors.New("JWT_SECRET not set")
			return
		}
		if len(s) < minSecretLength {
			secretErr = errors.New("JWT_SECRET too short; must be >= 32 bytes")
			return
		}
		secret = []byte(s)
	})
	return secret, secretErr
}

// CustomClaims use RegisteredClaims and a role field.
type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for the given user public ID and role.
func GenerateJWT(publicID string, role string) (string, error) {
	sec, err := loadSecret()
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   publicID,
			Issuer:    defaultIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultExpiryMins * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sec)
}

// ParseJWT validates and parses a JWT token string, returning the custom claims.
func ParseJWT(tokenStr string) (*CustomClaims, error) {
	sec, err := loadSecret()
	if err != nil {
		return nil, err
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return sec, nil
	}

	claims := &CustomClaims{}
	parser := jwt.NewParser(jwt.WithLeeway(time.Second * clockSkewSeconds))
	token, err := parser.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// Validate Issuer
	if claims.Issuer != defaultIssuer {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// Check Subject field
	if claims.Subject == "" {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// GenerateRefreshToken generates a cryptographically secure random refresh token.
// Returns the raw token (to send to client) and expiration time.
func GenerateRefreshToken() (string, time.Time, error) {
	bytes := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	expiresAt := time.Now().Add(refreshTokenDuration)
	return token, expiresAt, nil
}

// HashRefreshToken creates a SHA-256 hash of the refresh token for DB storage.
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
