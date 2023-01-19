package model

import (
	"reflect"
	"testing"
)

func TestUser_ID(t *testing.T) {
	tests := []struct {
		name string
		u    *user
		want string
	}{
		{
			name: "正常ケース",
			u:    &user{id: "id"},
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
		u    *user
		want string
	}{
		{
			name: "正常ケース",
			u:    &user{name: "name"},
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
		u    *user
		want string
	}{
		{
			name: "正常ケース",
			u:    &user{email: "email"},
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
		u    *user
		want string
	}{
		{
			name: "正常ケース",
			u:    &user{password: "password"},
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
		name    string
		u       *user
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "一致ケース",
			u: &user{
				password: "e+e8yncmvsKaaEGdcrvWMhjEvMH/3eyLGJ/mBlhFxQA=",
				salt:     "b",
			},
			args: args{
				password: "a",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "不一致ケース",
			u: &user{
				password: "a",
				salt:     "b",
			},
			args: args{
				password: "a",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.VerifyPassword(tt.args.password)
			if got != tt.want {
				t.Errorf("User.VerifyPassword() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("User.VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
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
		salt     string
	}
	tests := []struct {
		name string
		args args
		want *user
	}{
		{
			name: "正常ケース",
			args: args{
				id:       "id",
				name:     "name",
				email:    "email",
				password: "password",
				salt:     "salt",
			},
			want: &user{
				id:       "id",
				name:     "name",
				email:    "email",
				password: "password",
				salt:     "salt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.id, tt.args.name, tt.args.email, tt.args.password, tt.args.salt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_generateSalt(t *testing.T) {
	tests := []struct {
		name    string
		u       *user
		wantLen int
		wantErr bool
	}{
		{
			name:    "正常ケース",
			u:       &user{},
			wantLen: saltLen,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.generateSalt()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.generateSalt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("User.generateSalt() = %v, wantLen %v", got, tt.wantLen)
			}
			for _, r := range got {
				if !(r >= 33 && r <= 126) {
					t.Errorf("User.generateSalt() = %v, unexpect character %v", got, r)
				}
			}
		})
	}
}

func TestUser_encryptPassword(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		u       *user
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常ケース",
			u:    &user{},
			args: args{
				password: "a",
				salt:     "b",
			},
			want:    "e+e8yncmvsKaaEGdcrvWMhjEvMH/3eyLGJ/mBlhFxQA=",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.encryptPassword(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.encryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.encryptPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_EncryptPassword(t *testing.T) {
	tests := []struct {
		name     string
		u        *user
		notWant  string
		notWant2 string
		wantErr  bool
	}{
		{
			name: "正常ケース",
			u: &user{
				password: "a",
			},
			notWant:  "a",
			notWant2: "a",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, err := tt.u.EncryptPassword()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.EncryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.notWant {
				t.Errorf("User.EncryptPassword() = %v, notWant %v", got, tt.notWant)
			}
			if got2 == tt.notWant2 {
				t.Errorf("User.EncryptPassword() = %v, notWant2 %v", got2, tt.notWant)
			}
		})
	}
}

func Test_user_Salt(t *testing.T) {
	tests := []struct {
		name string
		u    *user
		want string
	}{
		{
			name: "正常ケース",
			u: &user{
				salt: "a",
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Salt(); got != tt.want {
				t.Errorf("user.Salt() = %v, want %v", got, tt.want)
			}
		})
	}
}
