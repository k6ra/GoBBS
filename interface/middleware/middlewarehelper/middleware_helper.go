package middlewarehelper

import (
	"GoBBS/interface/handler/handlerctx"
	"errors"
	"log"
	"net/http"
)

type (
	// Helper ミドルウェアヘルパー
	// mockgen -source interface/middleware/middleware_helper.go -destination mock/mock_middleware/middleware_helper_mock.go
	Helper interface {
		Apply(handlerctx.APIContextFactory, HandlerFunc, ...MiddlewareFunc) http.HandlerFunc
	}

	HandlerFunc    func(handlerctx.APIContext) error
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

var ErrContextKeyNotFound = errors.New("context key not found")

// Apply ハンドラーにミドルウェアを適用する
func Apply(ctxFactory handlerctx.APIContextFactory, h HandlerFunc, m ...MiddlewareFunc) http.HandlerFunc {
	f := h
	for i := len(m) - 1; i >= 0; i-- {
		f = m[i](f)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctxFactory(w, r)
		if err := f(c); err != nil {
			log.Printf("handlerfunc error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
