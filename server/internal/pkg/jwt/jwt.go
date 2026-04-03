package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"`
	jwtlib.StandardClaims
}

// Generate creates a new JWT token.
func Generate(secret string, userID int64, role, nickname string, expireHours int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Role:     role,
		Nickname: nickname,
		StandardClaims: jwtlib.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "jianghu-server",
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Validate parses and validates a JWT token.
func Validate(secret, tokenStr string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
