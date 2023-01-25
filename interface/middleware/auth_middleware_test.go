package middleware

import (
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
		header http.Header
		userID string
	}
	tests := []struct {
		name string
		m    *auth
		args args
		want bool
	}{
		{
			name: "正常ケース",
			m: &auth{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().VerifyAuthorization(gomock.Any(), gomock.Any()).Return(true)
					return mock
				}(),
			},
			args: args{
				header: map[string][]string{
					"Authorization": {"Bearer abc"},
				},
			},
			want: true,
		},
		{
			name: "認証エラーケース",
			m: &auth{
				uc: func() *mock_usecase.MockUser {
					mock := mock_usecase.NewMockUser(ctrl)
					mock.EXPECT().VerifyAuthorization(gomock.Any(), gomock.Any()).Return(false)
					return mock
				}(),
			},
			args: args{
				header: map[string][]string{
					"Authorization": {"Bearer abc"},
				},
			},
			want: false,
		},
		{
			name: "AuthorizationヘッダBearerなしケース",
			m:    &auth{},
			args: args{
				header: map[string][]string{
					"Authorization": {""},
				},
			},
			want: false,
		},
		{
			name: "Authorizationヘッダなしケース",
			m:    &auth{},
			args: args{
				header: map[string][]string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.VerifyAuth(tt.args.header, tt.args.userID); got != tt.want {
				t.Errorf("auth.VerifyAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
