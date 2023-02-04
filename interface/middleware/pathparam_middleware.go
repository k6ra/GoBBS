package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"errors"
	"net/http"
	"strings"
)

type (
	// PathParam パスパラメータミドルウェア
	// mockgen -source interface/middleware/pathparam_middleware.go -destination mock/mock_middleware/pathparam_middleware_mock.go
	PathParam interface {
		Parse(middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc
	}

	// pathParam パスパラメータミドルウェア
	pathParam struct {
		pathPattern string
	}
)

const paramPrefix = ":"

// NewPathParam パスパラメータミドルウェアを生成する
func NewPathParam(pathPattern string) *pathParam {
	return &pathParam{
		pathPattern: pathPattern,
	}
}

// Parse パスパラメータを取得する
func (m *pathParam) Parse(next middlewarehelper.HandlerFunc) middlewarehelper.HandlerFunc {
	return func(c handlerctx.APIContext) error {
		// パスパラメータは1個のみ指定可能、2個以上指定された場合は不正な指定
		if strings.Count(m.pathPattern, paramPrefix) > 1 {
			return errors.New("pathparam.Parse: invalid path pattern")
		}

		splitPaths := strings.Split(c.URL().Path, "/")
		splitPathPatterns := strings.Split(m.pathPattern, "/")
		if len(splitPaths) != len(splitPathPatterns) {
			c.WriteStatusCode(http.StatusBadRequest)
			return nil
		}

		for i, v := range strings.Split(m.pathPattern, "/") {
			if strings.HasPrefix(v, ":") {
				c.SetPathParam(splitPaths[i])
				continue
			}
			if splitPaths[i] != splitPathPatterns[i] {
				c.WriteStatusCode(http.StatusBadRequest)
				return nil
			}
		}

		return next(c)
	}
}
