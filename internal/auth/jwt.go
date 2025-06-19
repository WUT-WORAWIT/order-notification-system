package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKeyBytes = []byte(os.Getenv("JWT_SECRET_KEY"))

// CustomClaims defines the structure of our JWT claims
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token with custom claims.
func GenerateToken(username string) (string, error) {
	if len(secretKeyBytes) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set in environment variables")
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	// expirationTime := time.Now().Add(5 * time.Minute) // For testing shorter expiration

	claims := &CustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   username, // Using Subject for username is common practice
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// VerifyToken validates the token string and returns the parsed token with CustomClaims.
func VerifyToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	if len(secretKeyBytes) == 0 {
		return nil, nil, fmt.Errorf("JWT_SECRET_KEY is not set in environment variables")
	}
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKeyBytes, nil
	})
	return token, claims, err
}
