package security

import (
	"reflect"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestNewJWTToken(t *testing.T) {
	type args struct {
		secretKey string
	}
	tests := []struct {
		name string
		args args
		want *jwtToken
	}{
		{
			name: "正常ケース",
			args: args{
				secretKey: "key",
			},
			want: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJWTToken(tt.args.secretKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWTToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtToken_Generate(t *testing.T) {
	type args struct {
		userID string
		now    time.Time
	}
	tests := []struct {
		name    string
		j       *jwtToken
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常ケース",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
			args: args{
				userID: "uid",
				now:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI1MzQ4MDAsImlkIjoidWlkIn0.hzK_r1gcKXyNeKi-lp3Kk0TKCP5EKeCiDUNgMpptZ3g",
			wantErr: false,
		},
		{
			name: "異常ケース(署名エラー)",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: &jwt.SigningMethodHMAC{},
			},
			args: args{
				userID: "uid",
				now:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.Generate(tt.args.userID, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtToken.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jwtToken.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtToken_Verify(t *testing.T) {
	type args struct {
		tokenString string
		userID      string
	}
	tests := []struct {
		name string
		j    *jwtToken
		args args
		want bool
	}{
		{
			name: "検証成功",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImlkIjoidWlkIn0.KtfT-8QDcz5zzOGxYdfB5cHorDrtWuY1EOmJFlXRiOo",
				userID:      "uid",
			},
			want: true,
		},
		{
			name: "検証失敗(期限切れ)",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI1MzQ4MDAsImlkIjoidWlkIn0.hzK_r1gcKXyNeKi-lp3Kk0TKCP5EKeCiDUNgMpptZ3g",
				userID:      "uid",
			},
			want: false,
		},
		{
			name: "検証失敗(アルゴリズム不一致)",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI1MzQ4MDAsImlkIjoidWlkIn0.JXwNXypso6L53JDS0BH8RjY99QzaIMlNlnW9uCVbrMQ5wkpNqwT1-84pD9YEwWASlfWyCSJrnGINmNQfU07GbA",
				userID:      "uid",
			},
			want: false,
		},
		{
			name: "検証失敗(id不一致)",
			j: &jwtToken{
				secretKey:     "key",
				signingMethod: jwt.SigningMethodHS256,
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImp0aSI6InVpZCJ9.txa8RkeQiq_BJRBWlI2WdKrGih6HXe5narIKQwm1-no",
				userID:      "ng",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.Verify(tt.args.tokenString, tt.args.userID); got != tt.want {
				t.Errorf("jwtToken.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
