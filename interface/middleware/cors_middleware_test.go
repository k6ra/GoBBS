package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"GoBBS/mock/mock_handler/mock_handlerctx"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewCORS(t *testing.T) {
	type args struct {
		allowOrigin  string
		allowMethods []string
		allowHeaders []string
		maxAge       int
	}
	tests := []struct {
		name string
		args args
		want *cors
	}{
		{
			name: "正常ケース",
			args: args{
				allowOrigin:  "a",
				allowMethods: []string{"b", "c"},
				allowHeaders: []string{"d", "e"},
				maxAge:       10,
			},
			want: &cors{
				allowOrigin:  "a",
				allowMethods: []string{"b", "c"},
				allowHeaders: []string{"d", "e"},
				maxAge:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCORS(tt.args.allowOrigin, tt.args.allowMethods, tt.args.allowHeaders, tt.args.maxAge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCORS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cors_AddResponseHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		next middlewarehelper.HandlerFunc
	}
	tests := []struct {
		name       string
		cors       *cors
		args       args
		apiContext handlerctx.APIContext
	}{
		{
			name: "指定値",
			cors: &cors{
				allowOrigin:  "a",
				allowMethods: []string{"b", "c"},
				allowHeaders: []string{"d", "e"},
				maxAge:       10,
			},
			args: args{
				next: func(c handlerctx.APIContext) error {
					return nil
				},
			},
			apiContext: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				gomock.InOrder(
					mock.EXPECT().AddResponseHeader("Access-Control-Allow-Origin", "a"),
					mock.EXPECT().AddResponseHeader("Access-Control-Allow-Methods", "b"),
					mock.EXPECT().AddResponseHeader("Access-Control-Allow-Methods", "c"),
					mock.EXPECT().AddResponseHeader("Access-Control-Allow-Headers", "d"),
					mock.EXPECT().AddResponseHeader("Access-Control-Allow-Headers", "e"),
					mock.EXPECT().AddResponseHeader("Access-Control-Max-Age", "10"),
				)
				return mock
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cors.AddResponseHeader(tt.args.next)
			got(tt.apiContext)
		})
	}
}
