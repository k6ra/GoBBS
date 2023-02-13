// Code generated by MockGen. DO NOT EDIT.
// Source: interface/handler/handlerctx/api_context.go

// Package mock_handlerctx is a generated GoMock package.
package mock_handlerctx

import (
	io "io"
	http "net/http"
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAPIContext is a mock of APIContext interface.
type MockAPIContext struct {
	ctrl     *gomock.Controller
	recorder *MockAPIContextMockRecorder
}

// MockAPIContextMockRecorder is the mock recorder for MockAPIContext.
type MockAPIContextMockRecorder struct {
	mock *MockAPIContext
}

// NewMockAPIContext creates a new mock instance.
func NewMockAPIContext(ctrl *gomock.Controller) *MockAPIContext {
	mock := &MockAPIContext{ctrl: ctrl}
	mock.recorder = &MockAPIContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPIContext) EXPECT() *MockAPIContextMockRecorder {
	return m.recorder
}

// AddResponseHeader mocks base method.
func (m *MockAPIContext) AddResponseHeader(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddResponseHeader", arg0, arg1)
}

// AddResponseHeader indicates an expected call of AddResponseHeader.
func (mr *MockAPIContextMockRecorder) AddResponseHeader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddResponseHeader", reflect.TypeOf((*MockAPIContext)(nil).AddResponseHeader), arg0, arg1)
}

// PathParam mocks base method.
func (m *MockAPIContext) PathParam() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PathParam")
	ret0, _ := ret[0].(string)
	return ret0
}

// PathParam indicates an expected call of PathParam.
func (mr *MockAPIContextMockRecorder) PathParam() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PathParam", reflect.TypeOf((*MockAPIContext)(nil).PathParam))
}

// RequestBody mocks base method.
func (m *MockAPIContext) RequestBody() io.ReadCloser {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestBody")
	ret0, _ := ret[0].(io.ReadCloser)
	return ret0
}

// RequestBody indicates an expected call of RequestBody.
func (mr *MockAPIContextMockRecorder) RequestBody() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestBody", reflect.TypeOf((*MockAPIContext)(nil).RequestBody))
}

// RequestHeader mocks base method.
func (m *MockAPIContext) RequestHeader() http.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestHeader")
	ret0, _ := ret[0].(http.Header)
	return ret0
}

// RequestHeader indicates an expected call of RequestHeader.
func (mr *MockAPIContextMockRecorder) RequestHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestHeader", reflect.TypeOf((*MockAPIContext)(nil).RequestHeader))
}

// RequestMethod mocks base method.
func (m *MockAPIContext) RequestMethod() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestMethod")
	ret0, _ := ret[0].(string)
	return ret0
}

// RequestMethod indicates an expected call of RequestMethod.
func (mr *MockAPIContextMockRecorder) RequestMethod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestMethod", reflect.TypeOf((*MockAPIContext)(nil).RequestMethod))
}

// SetPathParam mocks base method.
func (m *MockAPIContext) SetPathParam(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPathParam", arg0)
}

// SetPathParam indicates an expected call of SetPathParam.
func (mr *MockAPIContextMockRecorder) SetPathParam(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPathParam", reflect.TypeOf((*MockAPIContext)(nil).SetPathParam), arg0)
}

// URL mocks base method.
func (m *MockAPIContext) URL() *url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(*url.URL)
	return ret0
}

// URL indicates an expected call of URL.
func (mr *MockAPIContextMockRecorder) URL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockAPIContext)(nil).URL))
}

// WriteResponseJSON mocks base method.
func (m *MockAPIContext) WriteResponseJSON(arg0 int, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteResponseJSON", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteResponseJSON indicates an expected call of WriteResponseJSON.
func (mr *MockAPIContextMockRecorder) WriteResponseJSON(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteResponseJSON", reflect.TypeOf((*MockAPIContext)(nil).WriteResponseJSON), arg0, arg1)
}

// WriteStatusCode mocks base method.
func (m *MockAPIContext) WriteStatusCode(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteStatusCode", arg0)
}

// WriteStatusCode indicates an expected call of WriteStatusCode.
func (mr *MockAPIContextMockRecorder) WriteStatusCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteStatusCode", reflect.TypeOf((*MockAPIContext)(nil).WriteStatusCode), arg0)
}
