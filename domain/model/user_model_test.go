package model

import (
	"reflect"
	"testing"
)

func TestUser_ID(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "正常ケース",
			u:    &User{id: "id"},
			want: "id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ID(); got != tt.want {
				t.Errorf("User.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Name(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "正常ケース",
			u:    &User{name: "name"},
			want: "name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Name(); got != tt.want {
				t.Errorf("User.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Email(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "正常ケース",
			u:    &User{email: "email"},
			want: "email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Email(); got != tt.want {
				t.Errorf("User.Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Password(t *testing.T) {
	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "正常ケース",
			u:    &User{password: "password"},
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Password(); got != tt.want {
				t.Errorf("User.Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		u    *User
		args args
		want bool
	}{
		{
			name: "一致ケース",
			u:    &User{password: "password"},
			args: args{password: "password"},
			want: true,
		},
		{
			name: "不一致ケース",
			u:    &User{password: "password"},
			args: args{password: "ng"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.VerifyPassword(tt.args.password); got != tt.want {
				t.Errorf("User.VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		id       string
		name     string
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "正常ケース",
			args: args{
				id:       "id",
				name:     "name",
				email:    "email",
				password: "password",
			},
			want: &User{
				id:       "id",
				name:     "name",
				email:    "email",
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.id, tt.args.name, tt.args.email, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
