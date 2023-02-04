package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"GoBBS/usecase"
	"net/http"
	"strings"
)

type (
	// Auth 認証ミドルウェア
	// mockgen -source interface/middleware/auth_middleware.go -destination mock/mock_middleware/auth_middleware_mock.go
	Auth interface {
		VerifyAuth(middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc
	}

	// auth 認証ミドルウェア
	auth struct {
		uc usecase.User
	}
)

var _ Auth = (*auth)(nil)

// NewAuth 認証ミドルウェアを生成する
func NewAuth(usecase usecase.User) *auth {
	return &auth{uc: usecase}
}

// VerifyAuth 認証する
func (m *auth) VerifyAuth(next middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc {
	return func(c handlerctx.APIContext) error {
		header := c.RequestHeader()
		if len(header["Authorization"]) == 0 {
			c.WriteStatusCode(http.StatusUnauthorized)
			return nil
		}

		auth := header["Authorization"][0]
		if !strings.HasPrefix(auth, "Bearer ") {
			c.WriteStatusCode(http.StatusUnauthorized)
			return nil
		}

		token := strings.Replace(auth, "Bearer ", "", 1)
		if ok := m.uc.VerifyAuthorization(token); !ok {
			c.WriteStatusCode(http.StatusUnauthorized)
			return nil
		}

		return next(c)
	}
}
