package dto

import (
	"reflect"
	"testing"
)

func TestNewToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "正常ケース",
			args: args{
				token: "token",
			},
			want: &Token{
				Token: "token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewToken(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
