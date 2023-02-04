package handler

import (
	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/mock/mock_handler/mock_handlerctx"
	"GoBBS/mock/mock_usecase"
	"GoBBS/usecase"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockUser(ctrl)

	type args struct {
		usecase usecase.User
	}
	tests := []struct {
		name string
		args args
		want *userHandler
	}{
		{
			name: "正常ケース",
			args: args{
				usecase: mockUC,
			},
			want: &userHandler{
				uc:          mockUC,
				jsonMarshal: json.Marshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.usecase); got.uc != tt.want.uc || reflect.ValueOf(got.jsonMarshal).Pointer() != reflect.ValueOf(tt.want.jsonMarshal).Pointer() {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHandler_RegistHandlerFunc(t *testing.T) {
	tests := []struct {
		name string
		h    *userHandler
	}{
		{
			name: "正常ケース",
			h:    &userHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.RegistHandlerFunc()
		})
	}
}

func Test_userHandler_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c handlerctx.APIContext
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().RequestMethod().Return(http.MethodPost),
						mock.EXPECT().WriteStatusCode(http.StatusOK),
					)
					return mock
				}(),
			},
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(メソッド不正)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().RequestMethod().Return(http.MethodPatch),
						mock.EXPECT().WriteStatusCode(http.StatusMethodNotAllowed),
					)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(リクエストボディ読み込みエラー)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{`))),
						mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
					)
					return mock
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.new(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("userHandler.new() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userHandler_edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c handlerctx.APIContext
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース(ユーザー更新)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().PathParam().Return("1"),
						mock.EXPECT().RequestMethod().Return(http.MethodPut),
						mock.EXPECT().WriteStatusCode(http.StatusOK),
					)
					return mock
				}(),
			},
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "正常ケース(ユーザー削除)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().PathParam().Return("1"),
						mock.EXPECT().RequestMethod().Return(http.MethodDelete),
						mock.EXPECT().WriteStatusCode(http.StatusOK),
					)
					return mock
				}(),
			},
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(nil)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(メソッド不正)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().PathParam().Return("1"),
						mock.EXPECT().RequestMethod().Return(http.MethodPatch),
						mock.EXPECT().WriteStatusCode(http.StatusMethodNotAllowed),
					)
					return mock
				}(),
			},
			h:       &userHandler{},
			wantErr: false,
		},
		{
			name: "異常ケース(パスパラメータ不正)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().PathParam().Return(""),
						mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
					)
					return mock
				}(),
			},
			h:       &userHandler{},
			wantErr: false,
		},
		{
			name: "異常ケース(リクエストボディ読み込みエラー)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{`))),
						mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
					)
					return mock
				}(),
			},
			h:       &userHandler{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.edit(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("userHandler.edit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userHandler_auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c handlerctx.APIContext
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().RequestMethod().Return(http.MethodPost),
						mock.EXPECT().WriteResponseJSON(http.StatusOK, []byte(`{"token": "abc"}`)),
					)
					return mock
				}(),
			},
			h: &userHandler{
				jsonMarshal: func(a any) ([]byte, error) {
					return []byte(`{"token": "abc"}`), nil
				},
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("abc", nil)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(メソッド不正)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{"id":"1"}`))),
						mock.EXPECT().RequestMethod().Return(http.MethodPatch),
						mock.EXPECT().WriteStatusCode(http.StatusMethodNotAllowed),
					)
					return mock
				}(),
			},
			h:       &userHandler{},
			wantErr: false,
		},
		{
			name: "異常ケース(リクエストボディ読み込みエラー)",
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					gomock.InOrder(
						mock.EXPECT().RequestBody().Return(io.NopCloser(bytes.NewBufferString(`{`))),
						mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
					)
					return mock
				}(),
			},
			h:       &userHandler{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.auth(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("userHandler.auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userHandler_regist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c    handlerctx.APIContext
		user dto.User
	}
	tests := []struct {
		name string
		h    *userHandler
		args args
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusOK)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(ユーザー登録済み)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(service.ErrUserAlreadyRegistered)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusBadRequest)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(想定外のエラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusInternalServerError)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.regist(tt.args.c, tt.args.user)
		})
	}
}

func Test_userHandler_update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c    handlerctx.APIContext
		user dto.User
	}
	tests := []struct {
		name string
		h    *userHandler
		args args
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusOK)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(ユーザー未登録)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(service.ErrUserNotFound)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusBadRequest)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(想定外のエラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusInternalServerError)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.update(tt.args.c, tt.args.user)
		})
	}
}

func Test_userHandler_delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c    handlerctx.APIContext
		user dto.User
	}
	tests := []struct {
		name string
		h    *userHandler
		args args
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusOK)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(ユーザー未登録)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(service.ErrUserNotFound)
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusBadRequest)
					return mock
				}(),
			},
		},
		{
			name: "異常ケース(想定外のエラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusInternalServerError)
					return mock
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.delete(tt.args.c, tt.args.user)
		})
	}
}

func Test_userHandler_login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c    handlerctx.APIContext
		user dto.User
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("abc", nil)
					return mock
				}(),
				jsonMarshal: func(a any) ([]byte, error) {
					return []byte(`{"token": "abc"}`), nil
				},
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteResponseJSON(http.StatusOK, []byte(`{"token": "abc"}`)).Return(nil)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(JSONマーシャルエラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("abc", nil)
					return mock
				}(),
				jsonMarshal: func(a any) ([]byte, error) {
					return []byte(""), errors.New("ng")
				},
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusInternalServerError)
					return mock
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(認証エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("", errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				c: func() *mock_handlerctx.MockAPIContext {
					mock := mock_handlerctx.NewMockAPIContext(ctrl)
					mock.EXPECT().WriteStatusCode(http.StatusUnauthorized)
					return mock
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.login(tt.args.c, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userHandler.login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
