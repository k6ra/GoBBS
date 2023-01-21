package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	// Token トークン
	// mockgen -source interface/security/jwt_token.go -destination mock/mock_security/jwt_token_mock.go
	Token interface {
		Generate(string, time.Time) (string, error)
		Verify(string, string) bool
	}

	// jwtToken jwtトークン
	jwtToken struct {
		secretKey     string
		signingMethod jwt.SigningMethod
	}
)

var _ Token = (*jwtToken)(nil)

// NewJWTToken jwtトークンを生成する
func NewJWTToken(secretKey string) *jwtToken {
	return &jwtToken{
		secretKey:     secretKey,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// Generate トークンを生成する
func (j *jwtToken) Generate(userID string, now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"exp": now.Add(time.Hour * 1).Unix(),
		"id":  userID,
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", errors.Wrap(err, "Generate error")
	}

	return tokenString, nil
}

// Verify トークンを検証する
func (j *jwtToken) Verify(tokenString string, userID string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if !(token.Method.Alg() == j.signingMethod.Alg()) {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	if claims["id"] != userID {
		return false
	}

	return true
}
