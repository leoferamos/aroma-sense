package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	return []byte(secret)
}()

// GenerateJWT creates a JWT token for a given public_id and role
func GenerateJWT(publicID string, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  publicID,
		"role": role,
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
		"iss":  "aroma-sense-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT validates and returns the claims of the token
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
