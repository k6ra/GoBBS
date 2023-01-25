package middleware

import (
	"GoBBS/usecase"
	"net/http"
	"strings"
)

type (
	// Auth 認証ミドルウェア
	// mockgen -source interface/middleware/auth_middleware.go -destination mock/mock_middleware/auth_middleware_mock.go
	Auth interface {
		VerifyAuth(http.Header, string) bool
	}

	// auth 認証ミドルウェア
	auth struct {
		uc usecase.User
	}
)

// NewAuth 認証ミドルウェアを生成する
func NewAuth(usecase usecase.User) *auth {
	return &auth{uc: usecase}
}

// VerifyAuth 認証をチェックする
func (m *auth) VerifyAuth(header http.Header, userID string) bool {
	if len(header["Authorization"]) == 0 {
		return false
	}

	auth := header["Authorization"][0]
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}

	token := strings.Replace(auth, "Bearer ", "", 1)
	if ok := m.uc.VerifyAuthorization(token, userID); !ok {
		return false
	}

	return true
}
