package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
)

func TestGenerateAndValidate(t *testing.T) {
	secret := "test-secret-key"
	userID := int64(42)
	role := "admin"
	nickname := "江湖侠客"

	token, err := Generate(secret, userID, role, nickname, 24)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned empty token")
	}

	claims, err := Validate(secret, token)
	if err != nil {
		t.Fatalf("Validate() error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("Role = %q, want %q", claims.Role, role)
	}
	if claims.Nickname != nickname {
		t.Errorf("Nickname = %q, want %q", claims.Nickname, nickname)
	}
	if claims.Issuer != "jianghu-server" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "jianghu-server")
	}
}

func TestValidate_WrongSecret(t *testing.T) {
	token, _ := Generate("secret-a", 1, "user", "test", 24)
	_, err := Validate("secret-b", token)
	if err == nil {
		t.Fatal("Validate() with wrong secret should return error")
	}
}

func TestValidate_ExpiredToken(t *testing.T) {
	// Create a token that's already expired by directly crafting claims
	claims := Claims{
		UserID:   1,
		Role:     "user",
		Nickname: "test",
		StandardClaims: jwtlib.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			Issuer:    "jianghu-server",
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	_, err = Validate("secret", tokenStr)
	if err == nil {
		t.Fatal("Validate() with expired token should return error")
	}
}

func TestValidate_InvalidToken(t *testing.T) {
	_, err := Validate("secret", "not-a-valid-token")
	if err == nil {
		t.Fatal("Validate() with invalid token should return error")
	}
}

func TestValidate_EmptyToken(t *testing.T) {
	_, err := Validate("secret", "")
	if err == nil {
		t.Fatal("Validate() with empty token should return error")
	}
}
