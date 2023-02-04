package middlewarehelper

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/mock/mock_handler/mock_handlerctx"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctxFactory handlerctx.APIContextFactory
		h          HandlerFunc
		m          []MiddlewareFunc
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "正常ケース",
			args: args{
				ctxFactory: func(w http.ResponseWriter, r *http.Request) handlerctx.APIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().SetPathParam(gomock.Any()),
						mock.EXPECT().PathParam(),
						mock.EXPECT().WriteResponseJSON(gomock.Any(), gomock.Any()).Return(nil),
					)
					return mock
				},
				h: func() HandlerFunc {
					return func(c handlerctx.APIContext) error {
						pathParam := c.PathParam()
						if err := c.WriteResponseJSON(
							http.StatusOK,
							fmt.Sprintf(`{"pathParam": "%s"}`,
								pathParam)); err != nil {
							t.Fatalf("write response json error(error: %v)", err)
						}
						return nil
					}
				}(),
				m: func() []MiddlewareFunc {
					return []MiddlewareFunc{
						func(next HandlerFunc) HandlerFunc {
							return func(c handlerctx.APIContext) error {
								c.SetPathParam("a")
								next(c)
								return nil
							}
						},
						func(next HandlerFunc) HandlerFunc {
							return func(c handlerctx.APIContext) error {
								next(c)
								return nil
							}
						},
					}
				}(),
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "異常ケース",
			args: args{
				ctxFactory: func(w http.ResponseWriter, r *http.Request) handlerctx.APIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().SetPathParam(gomock.Any())
					return mock
				},
				h: func() HandlerFunc {
					return func(c handlerctx.APIContext) error {
						return nil
					}
				}(),
				m: func() []MiddlewareFunc {
					return []MiddlewareFunc{
						func(next HandlerFunc) HandlerFunc {
							return func(c handlerctx.APIContext) error {
								c.SetPathParam("a")
								return errors.New("ng")
							}
						},
						func(next HandlerFunc) HandlerFunc {
							return func(c handlerctx.APIContext) error {
								next(c)
								return nil
							}
						},
					}
				}(),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Apply(tt.args.ctxFactory, tt.args.h, tt.args.m...)

			w := &httptest.ResponseRecorder{Code: http.StatusOK}
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))

			got(w, r)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Apply function statusCode = %v, want %v", w.Code, tt.wantStatusCode)
			}
		})
	}
}
