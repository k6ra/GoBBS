package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"GoBBS/mock/mock_handler/mock_handlerctx"
	"GoBBS/mock/mock_usecase"
	"GoBBS/usecase"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		usecase usecase.User
	}
	tests := []struct {
		name string
		args args
		want *auth
	}{
		{
			name: "正常ケース",
			args: args{
				usecase: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					return mock
				}(),
			},
			want: &auth{
				uc: mock_usecase.NewMockUser(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuth(tt.args.usecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_VerifyAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		next middlewarehelper.HandlerFunc
	}
	tests := []struct {
		name        string
		m           *auth
		args        args
		ctx         handlerctx.APIContext
		wantFuncErr bool
	}{
		{
			name: "正常ケース",
			m: &auth{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().VerifyAuthorization(gomock.Any()).Return(true)
					return mock
				}(),
			},
			args: args{
				next: func() middlewarehelper.HandlerFunc {
					return func(c handlerctx.APIContext) error {
						return nil
					}
				}(),
			},
			ctx: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				header := http.Header{"Authorization": {"Bearer abc"}}
				mock.EXPECT().RequestHeader().Return(header)
				return mock
			}(),
			wantFuncErr: false,
		},
		{
			name: "異常ケース(認証失敗)",
			m: &auth{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().VerifyAuthorization(gomock.Any()).Return(false)
					return mock
				}(),
			},
			args: args{
				next: func() middlewarehelper.HandlerFunc {
					return func(c handlerctx.APIContext) error {
						return nil
					}
				}(),
			},
			ctx: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				header := http.Header{"Authorization": {"Bearer abc"}}
				gomock.InOrder(
					mock.EXPECT().RequestHeader().Return(header),
					mock.EXPECT().WriteStatusCode(http.StatusUnauthorized),
				)
				return mock
			}(),
			wantFuncErr: false,
		},
		{
			name: "異常ケース(AuthorizationヘッダーがBearerで始まらない)",
			m:    &auth{},
			args: args{
				next: func() middlewarehelper.HandlerFunc {
					return func(c handlerctx.APIContext) error {
						return nil
					}
				}(),
			},
			ctx: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				header := http.Header{"Authorization": {"abc"}}
				gomock.InOrder(
					mock.EXPECT().RequestHeader().Return(header),
					mock.EXPECT().WriteStatusCode(http.StatusUnauthorized),
				)
				return mock
			}(),
			wantFuncErr: false,
		},
		{
			name: "異常ケース(Authorizationヘッダーがない)",
			m:    &auth{},
			args: args{
				next: func() middlewarehelper.HandlerFunc {
					return func(c handlerctx.APIContext) error {
						return nil
					}
				}(),
			},
			ctx: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				header := http.Header{}
				gomock.InOrder(
					mock.EXPECT().RequestHeader().Return(header),
					mock.EXPECT().WriteStatusCode(http.StatusUnauthorized),
				)
				return mock
			}(),
			wantFuncErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.VerifyAuth(tt.args.next)
			err := got(tt.ctx)
			if (err != nil) != tt.wantFuncErr {
				t.Errorf("auth.VerifyAuth() error = %v, want %v", err, tt.wantFuncErr)
			}
		})
	}
}
