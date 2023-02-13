package handlerctx

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type (
	// APIContext APIコンテキストインターフェース
	// mockgen -source interface/handler/handlerctx/api_context.go -destination mock/mock_handler/mock_handlerctx/api_context_mock.go
	APIContext interface {
		WriteStatusCode(int)
		WriteResponseJSON(int, any) error
		RequestHeader() http.Header
		PathParam() string
		SetPathParam(string)
		URL() *url.URL
		RequestBody() io.ReadCloser
		RequestMethod() string
		AddResponseHeader(string, string)
	}

	// APIContextFactory APIコンテキストファクトリー
	APIContextFactory func(http.ResponseWriter, *http.Request) APIContext

	// apiContext APIコンテキスト
	apiContext struct {
		jsonMarshal func(any) ([]byte, error)
		pathParam   string
		request     *http.Request
		response    http.ResponseWriter
	}
)

var _ APIContext = (*apiContext)(nil)

// NewAPIContext APIコンテキストを生成する
func NewAPIContext(w http.ResponseWriter, r *http.Request) APIContext {
	return &apiContext{
		jsonMarshal: json.Marshal,
		request:     r,
		response:    w,
	}
}

// PathParam パスパラメーターを返す
func (c *apiContext) PathParam() string {
	return c.pathParam
}

// SetPathParam パスパラメータをセットする
func (c *apiContext) SetPathParam(param string) {
	c.pathParam = param
}

// WriteStatusCode レスポンスのステータスコードをセットする
func (c *apiContext) WriteStatusCode(statusCode int) {
	c.response.WriteHeader(statusCode)
}

// WriteResponseJSON レスポンスのステータスコード、JSONをセットする
func (c *apiContext) WriteResponseJSON(statusCode int, json any) error {
	jsonByte, err := c.jsonMarshal(json)
	if err != nil {
		return err
	}

	c.response.WriteHeader(statusCode)
	c.response.Header().Add("Content-Type", "application/json")
	c.response.Write(jsonByte)

	return nil
}

// RequestHeader リクエストヘッダーを返す
func (c *apiContext) RequestHeader() http.Header {
	return c.request.Header
}

// URL URLを返す
func (c *apiContext) URL() *url.URL {
	return c.request.URL
}

// RequestBody リクエストボディを返す
func (c *apiContext) RequestBody() io.ReadCloser {
	return c.request.Body
}

// RequestMethod リクエストメソッドを変えs
func (c *apiContext) RequestMethod() string {
	return c.request.Method
}

// AddResponseHeader レスポンスヘッダを追加する
func (c *apiContext) AddResponseHeader(key string, value string) {
	c.response.Header().Add(key, value)
}
