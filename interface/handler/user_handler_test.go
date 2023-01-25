package handler

import (
	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/mock/mock_middleware"
	"GoBBS/mock/mock_usecase"
	"GoBBS/usecase"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func Test_userHandler_regist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w    http.ResponseWriter
		user dto.User
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
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
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "異常ケース(登録エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusInternalServerError)
				return res
			}(),
		},
		{
			name: "異常ケース(登録済みエラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(service.ErrUserAlreadyRegistered)
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.regist(tt.args.w, tt.args.user)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("regist() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w    http.ResponseWriter
		user dto.User
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
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
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "異常ケース(更新エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusInternalServerError)
				return res
			}(),
		},
		{
			name: "異常ケース(ユーザー未登録エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(service.ErrUserNotFound)
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.update(tt.args.w, tt.args.user)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("update() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w    http.ResponseWriter
		user dto.User
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
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
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "異常ケース(削除エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusInternalServerError)
				return res
			}(),
		},
		{
			name: "異常ケース(ユーザー未登録エラー)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(service.ErrUserNotFound)
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.delete(tt.args.w, tt.args.user)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("delete() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w    http.ResponseWriter
		user dto.User
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("token", nil)
					return mock
				}(),
				jsonMarshal: json.Marshal,
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				res.Header().Add("Content-Type", "application/json")
				res.Write([]byte(`{Token: "token"}`))
				return res
			}(),
		},
		{
			name: "jsonマーシャルエラーケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("token", nil)
					return mock
				}(),
				jsonMarshal: func(any) ([]byte, error) {
					return nil, errors.New("ng")
				},
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusInternalServerError)
				return res
			}(),
		},
		{
			name: "認証エラーケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("", errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				w:    &httptest.ResponseRecorder{},
				user: dto.User{},
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusUnauthorized)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.login(tt.args.w, tt.args.user)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("login() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_getUserFromReqBody(t *testing.T) {
	type args struct {
		body io.ReadCloser
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		want    dto.User
		wantErr bool
	}{
		{
			name: "正常ケース",
			h:    &userHandler{},
			args: args{
				body: io.NopCloser(strings.NewReader(`{"id": "1"}`)),
			},
			want:    dto.User{ID: "1"},
			wantErr: false,
		},
		{
			name: "json不正ケース",
			h:    &userHandler{},
			args: args{
				body: io.NopCloser(strings.NewReader(`{`)),
			},
			want:    dto.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.getUserFromReqBody(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("userHandler.getUserFromReqBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userHandler.getUserFromReqBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHandler_getUserIDFromPathParam(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		h    *userHandler
		args args
		want string
	}{
		{
			name: "正常ケース",
			h:    &userHandler{},
			args: args{
				path: "/user/123",
			},
			want: "123",
		},
		{
			name: "path不正ケース",
			h:    &userHandler{},
			args: args{
				path: "/user/123/",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.getUserIDFromPathParam(tt.args.path); got != tt.want {
				t.Errorf("userHandler.getUserIDFromPathParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHandler_new(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost/user",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "想定外メソッドケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost/user",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusMethodNotAllowed)
				return res
			}(),
		},
		{
			name: "Json読み込みエラーケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost/user",
					bytes.NewBufferString(`{`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.new(tt.args.w, tt.args.r)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("new() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
	}{
		{
			name: "正常ケース(update)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
				authMiddleware: func() *mock_middleware.MockAuth {
					mock := mock_middleware.NewMockAuth(ctrl)
					mock.EXPECT().VerifyAuth(gomock.Any(), gomock.Any()).Return(true)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost/users/1",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "正常ケース(delete)",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(nil)
					return mock
				}(),
				authMiddleware: func() *mock_middleware.MockAuth {
					mock := mock_middleware.NewMockAuth(ctrl)
					mock.EXPECT().VerifyAuth(gomock.Any(), gomock.Any()).Return(true)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost/users/1",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusOK)
				return res
			}(),
		},
		{
			name: "想定外メソッドケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost/users/1",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusMethodNotAllowed)
				return res
			}(),
		},
		{
			name: "認証エラーケース(update)",
			h: &userHandler{
				authMiddleware: func() *mock_middleware.MockAuth {
					mock := mock_middleware.NewMockAuth(ctrl)
					mock.EXPECT().VerifyAuth(gomock.Any(), gomock.Any()).Return(false)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost/users/1",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusUnauthorized)
				return res
			}(),
		},
		{
			name: "認証エラーケース(delete)",
			h: &userHandler{
				authMiddleware: func() *mock_middleware.MockAuth {
					mock := mock_middleware.NewMockAuth(ctrl)
					mock.EXPECT().VerifyAuth(gomock.Any(), gomock.Any()).Return(false)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost/users/1",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusUnauthorized)
				return res
			}(),
		},
		{
			name: "パスパラメータ取得エラーケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost/users/",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
		{
			name: "リクエストBody取得エラーケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost/users/",
					bytes.NewBufferString(`{`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.edit(tt.args.w, tt.args.r)
		})
	}
}

func Test_userHandler_auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		h       *userHandler
		args    args
		wantRes http.ResponseWriter
	}{
		{
			name: "正常ケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return("token", nil)
					return mock
				}(),
				jsonMarshal: json.Marshal,
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost/auth",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.Header().Add("Content-Type", "application/json")
				res.WriteHeader(http.StatusOK)
				if _, err := res.Write([]byte(`{"token": "token"}`)); err != nil {
					t.Fatalf("response write fail %v", err)
				}
				return res
			}(),
		},
		{
			name: "想定外メソッドケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost/auth",
					bytes.NewBufferString(`{"id": "1"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusMethodNotAllowed)
				return res
			}(),
		},
		{
			name: "リクエストBodyエラーケース",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost/auth",
					bytes.NewBufferString(`{`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.auth(tt.args.w, tt.args.r)
		})
	}
}
