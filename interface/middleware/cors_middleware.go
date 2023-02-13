package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"strconv"
)

type (
	// CORS CORSミドルウェア
	// mockgen -source interface/middleware/cors_middleware.go -destination mock/mock_middleware/cors_middleware_mock.go
	CORS interface {
		AddResponseHeader(middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc
	}

	// cors corsミドルウェア
	cors struct {
		allowOrigin  string
		allowMethods []string
		allowHeaders []string
		maxAge       int
	}
)

var _ CORS = (*cors)(nil)

// NewCORS corsミドルウェアを生成する
func NewCORS(allowOrigin string, allowMethods []string, allowHeaders []string, maxAge int) *cors {
	return &cors{
		allowOrigin:  allowOrigin,
		allowMethods: allowMethods,
		allowHeaders: allowHeaders,
		maxAge:       maxAge,
	}
}

// AddResponseHeader レスポンスヘッダを追加
func (cors *cors) AddResponseHeader(next middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc {
	return func(c handlerctx.APIContext) error {
		c.AddResponseHeader("Access-Control-Allow-Origin", cors.allowOrigin)

		for _, method := range cors.allowMethods {
			c.AddResponseHeader("Access-Control-Allow-Methods", method)
		}

		for _, header := range cors.allowHeaders {
			c.AddResponseHeader("Access-Control-Allow-Headers", header)
		}

		c.AddResponseHeader("Access-Control-Max-Age", strconv.Itoa(cors.maxAge))

		return next(c)
	}
}
