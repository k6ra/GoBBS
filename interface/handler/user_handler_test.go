package handler

import (
	"GoBBS/domain/service"
	"GoBBS/mock/mock_usecase"
	"GoBBS/usecase"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				usecase: mock_usecase.NewMockUser(ctrl),
			},
			want: &userHandler{
				uc: mock_usecase.NewMockUser(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.usecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userHandler_RegistUserHandlerFunc(t *testing.T) {
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
			tt.h.RegistUserHandlerFunc()
		})
	}
}

func Test_userHandler_regist(t *testing.T) {
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
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
		{
			name: "異常ケース(Jsonエラー)",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPost,
					"/user",
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
			tt.h.regist(tt.args.w, tt.args.r)
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
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
		{
			name: "異常ケース(Jsonエラー)",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodPut,
					"/user",
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
			tt.h.update(tt.args.w, tt.args.r)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("regist() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_delete(t *testing.T) {
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
					mock.EXPECT().Delete(gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
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
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"/user",
					bytes.NewBufferString(`{"id": "a", "name": "b", "email": "c", "password": "d"}`),
				),
			},
			wantRes: func() http.ResponseWriter {
				res := &httptest.ResponseRecorder{}
				res.WriteHeader(http.StatusBadRequest)
				return res
			}(),
		},
		{
			name: "異常ケース(Jsonエラー)",
			h:    &userHandler{},
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(
					http.MethodDelete,
					"/user",
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
			tt.h.delete(tt.args.w, tt.args.r)
			if !reflect.DeepEqual(tt.args.w, tt.wantRes) {
				t.Errorf("regist() = %#v, want %#v", tt.args.w, tt.wantRes)
			}
		})
	}
}

func Test_userHandler_userHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		h    *userHandler
		args args
	}{
		{
			name: "postケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/user",
					bytes.NewBufferString(`{}`),
				),
			},
		},
		{
			name: "putケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"/user",
					bytes.NewBufferString(`{}`),
				),
			},
		},
		{
			name: "deleteケース",
			h: &userHandler{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().Delete(gomock.Any()).Return(nil)
					return mock
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"/user",
					bytes.NewBufferString(`{}`),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.userHandler(tt.args.w, tt.args.r)
		})
	}
}
