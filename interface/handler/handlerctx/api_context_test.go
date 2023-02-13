package handlerctx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewAPIContext(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *apiContext
	}{
		{
			name: "正常ケース",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("1")),
			},
			want: &apiContext{
				jsonMarshal: json.Marshal,
				request:     httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("1")),
				response:    httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAPIContext(tt.args.w, tt.args.r)

			c, ok := got.(*apiContext)
			if !ok {
				t.Errorf("NewAPIContext() = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(
				reflect.ValueOf(c.jsonMarshal).Pointer(),
				reflect.ValueOf(tt.want.jsonMarshal).Pointer()) {
				t.Errorf("NewAPIContext() = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(c.request, tt.want.request) {
				t.Errorf("NewAPIContext() = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(c.response, tt.want.response) {
				t.Errorf("NewAPIContext() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_PathParam(t *testing.T) {
	tests := []struct {
		name string
		c    *apiContext
		want string
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				pathParam: "a",
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.PathParam(); got != tt.want {
				t.Errorf("apiContext.PathParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_SetPathParam(t *testing.T) {
	type args struct {
		param string
	}
	tests := []struct {
		name           string
		c              *apiContext
		args           args
		wantAPIContext *apiContext
	}{
		{
			name: "正常ケース",
			c:    &apiContext{},
			args: args{
				param: "a",
			},
			wantAPIContext: &apiContext{
				pathParam: "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetPathParam(tt.args.param)
		})

		if !reflect.DeepEqual(tt.c, tt.wantAPIContext) {
			t.Errorf("apiContext = %#v, want %#v", tt.c, tt.wantAPIContext)
		}
	}
}

func Test_apiContext_WriteStatusCode(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name           string
		c              *apiContext
		args           args
		wantAPIContext *apiContext
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				response: &httptest.ResponseRecorder{},
			},
			args: args{
				statusCode: http.StatusOK,
			},
			wantAPIContext: &apiContext{
				response: &httptest.ResponseRecorder{
					Code: http.StatusOK,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.WriteStatusCode(tt.args.statusCode)

			if !reflect.DeepEqual(tt.c.response.Header(), tt.wantAPIContext.response.Header()) {
				t.Errorf("apiContext.WriteStatusCode() = %#v, want %#v", tt.c.response.Header(), tt.wantAPIContext.response.Header())
			}
		})
	}
}

func Test_apiContext_WriteResponseJSON(t *testing.T) {
	type args struct {
		statusCode int
		json       any
	}
	tests := []struct {
		name           string
		c              *apiContext
		args           args
		wantAPIContext *apiContext
		wantErr        bool
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				jsonMarshal: json.Marshal,
				response:    &httptest.ResponseRecorder{},
			},
			args: args{
				statusCode: http.StatusOK,
				json:       `{"test", "a"}`,
			},
			wantAPIContext: &apiContext{
				response: &httptest.ResponseRecorder{
					Code: http.StatusOK,
					HeaderMap: http.Header{
						"Content-Type": {"application/json"},
					},
					Body: bytes.NewBuffer([]byte(`{"test", "a"}`)),
				},
			},
			wantErr: false,
		},
		{
			name: "異常ケース",
			c: &apiContext{
				jsonMarshal: func(any) ([]byte, error) {
					return nil, errors.New("ng")
				},
				response: &httptest.ResponseRecorder{},
			},
			args: args{
				statusCode: http.StatusOK,
				json:       `{"test", "a"}`,
			},
			wantAPIContext: &apiContext{
				response: &httptest.ResponseRecorder{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.WriteResponseJSON(tt.args.statusCode, tt.args.json); (err != nil) != tt.wantErr {
				t.Errorf("apiContext.WriteResponseJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.c.response.Header(), tt.wantAPIContext.response.Header()) {
				t.Errorf("apiContext.WriteResponseJSON()= %#v, want %#v", tt.c.response, tt.wantAPIContext)
			}
		})
	}
}

func Test_apiContext_RequestHeader(t *testing.T) {
	tests := []struct {
		name string
		c    *apiContext
		want http.Header
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				request: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(""))
					req.Header = http.Header{
						"header": {"test"},
					}
					return req
				}(),
			},
			want: http.Header{
				"header": {"test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.RequestHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiContext.RequestHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_URL(t *testing.T) {
	tests := []struct {
		name string
		c    *apiContext
		want *url.URL
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				request: httptest.NewRequest(http.MethodPost, "http://host/parent/sub/1", bytes.NewBufferString("")),
			},
			want: &url.URL{
				Scheme:      "http",
				Opaque:      "",
				User:        nil,
				Host:        "host",
				Path:        "/parent/sub/1",
				RawPath:     "",
				OmitHost:    false,
				ForceQuery:  false,
				RawQuery:    "",
				Fragment:    "",
				RawFragment: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.URL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiContext.URL() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_RequestBody(t *testing.T) {
	tests := []struct {
		name string
		c    *apiContext
		want io.ReadCloser
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				request: &http.Request{
					Body: io.NopCloser(bytes.NewReader([]byte("abc"))),
				},
			},
			want: io.NopCloser(bytes.NewReader([]byte("abc"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.RequestBody()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apiContext.RequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_RequestMethod(t *testing.T) {
	tests := []struct {
		name string
		c    *apiContext
		want string
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				request: &http.Request{Method: http.MethodPost},
			},
			want: http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.RequestMethod(); got != tt.want {
				t.Errorf("apiContext.RequestMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiContext_AddResponseHeader(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name               string
		c                  *apiContext
		args               args
		wantResponseHeader http.Header
	}{
		{
			name: "正常ケース",
			c: &apiContext{
				response: httptest.NewRecorder(),
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantResponseHeader: http.Header{
				"Key": {"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AddResponseHeader(tt.args.key, tt.args.value)

			if !reflect.DeepEqual(tt.c.response.Header(), tt.wantResponseHeader) {
				t.Errorf("apiContext.AddResponseHeader() = %v, want %v", tt.c.response.Header(), tt.wantResponseHeader)

			}
		})
	}
}
