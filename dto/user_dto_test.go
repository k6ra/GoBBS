package dto

import (
	"GoBBS/domain/model"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "正常ケース",
			args: args{
				user: model.NewUser("id", "name", "email", "password"),
			},
			want: &User{
				ID:       "id",
				Name:     "name",
				Email:    "email",
				Password: "password",
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
		want *model.User
	}{
		{
			name: "正常ケース",
			u: &User{
				ID:       "id",
				Name:     "name",
				Email:    "email",
				Password: "password",
			},
			want: model.NewUser("id", "name", "email", "password"),
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
