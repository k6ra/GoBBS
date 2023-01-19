package dto

import (
	"GoBBS/domain/model"
	"GoBBS/mock/mock_model"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user model.User
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "正常ケース",
			args: args{
				user: func() model.User {
					mock := mock_model.NewMockUser(ctrl)
					gomock.InOrder(
						mock.EXPECT().ID().Return("id"),
						mock.EXPECT().Name().Return("name"),
						mock.EXPECT().Email().Return("email"),
						mock.EXPECT().Password().Return("password"),
						mock.EXPECT().Salt().Return("salt"),
					)
					return mock
				}(),
			},
			want: &User{
				ID:       "id",
				Name:     "name",
				Email:    "email",
				Password: "password",
				Salt:     "salt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_MapUserModel(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want model.User
	}{
		{
			name: "正常ケース",
			u: &User{
				ID:       "id",
				Name:     "name",
				Email:    "email",
				Password: "password",
				Salt:     "salt",
			},
			want: model.NewUser("id", "name", "email", "password", "salt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.MapUserModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.MapUserModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
