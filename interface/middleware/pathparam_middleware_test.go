package middleware

import (
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware/middlewarehelper"
	"GoBBS/mock/mock_handler/mock_handlerctx"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewPathParam(t *testing.T) {
	type args struct {
		pathPattern string
	}
	tests := []struct {
		name string
		args args
		want *pathParam
	}{
		{
			name: "正常ケース",
			args: args{
				pathPattern: "/abc",
			},
			want: &pathParam{
				pathPattern: "/abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPathParam(tt.args.pathPattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPathParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathParam_Parse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		next middlewarehelper.HandlerFunc
	}
	tests := []struct {
		name    string
		m       *pathParam
		args    args
		context handlerctx.APIContext
		wantErr bool
	}{
		{
			name: "正常ケース",
			m: &pathParam{
				pathPattern: "/a/b/:id",
			},
			args: args{
				next: func(c handlerctx.APIContext) error {
					return nil
				},
			},
			context: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				gomock.InOrder(
					mock.EXPECT().URL().Return(&url.URL{Path: "/a/b/123"}),
					mock.EXPECT().SetPathParam("123"),
				)
				return mock
			}(),
			wantErr: false,
		},
		{
			name: "異常ケース(パス不一致)",
			m: &pathParam{
				pathPattern: "/a/b/:id",
			},
			args: args{},
			context: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				gomock.InOrder(
					mock.EXPECT().URL().Return(&url.URL{Path: "/a/c/123"}),
					mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
				)
				return mock
			}(),
			wantErr: false,
		},
		{
			name: "異常ケース(パス不一致2)",
			m: &pathParam{
				pathPattern: "/a/b/:id",
			},
			args: args{},
			context: func() *mock_handlerctx.MockAPIContext {
				mock := mock_handlerctx.NewMockAPIContext(ctrl)
				gomock.InOrder(
					mock.EXPECT().URL().Return(&url.URL{Path: "/a"}),
					mock.EXPECT().WriteStatusCode(http.StatusBadRequest),
				)
				return mock
			}(),
			wantErr: false,
		},
		{
			name: "異常ケース(パスパラメータなし)",
			m: &pathParam{
				pathPattern: "/a/b/:id/:no",
			},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Parse(tt.args.next)

			err := got(tt.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathParam.Parse() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
