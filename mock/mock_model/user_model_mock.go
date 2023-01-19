// Code generated by MockGen. DO NOT EDIT.
// Source: domain/model/user_model.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Email mocks base method.
func (m *MockUser) Email() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Email")
	ret0, _ := ret[0].(string)
	return ret0
}

// Email indicates an expected call of Email.
func (mr *MockUserMockRecorder) Email() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Email", reflect.TypeOf((*MockUser)(nil).Email))
}

// EncryptPassword mocks base method.
func (m *MockUser) EncryptPassword() (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptPassword")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EncryptPassword indicates an expected call of EncryptPassword.
func (mr *MockUserMockRecorder) EncryptPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptPassword", reflect.TypeOf((*MockUser)(nil).EncryptPassword))
}

// ID mocks base method.
func (m *MockUser) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockUserMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockUser)(nil).ID))
}

// Name mocks base method.
func (m *MockUser) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockUserMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockUser)(nil).Name))
}

// Password mocks base method.
func (m *MockUser) Password() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Password")
	ret0, _ := ret[0].(string)
	return ret0
}

// Password indicates an expected call of Password.
func (mr *MockUserMockRecorder) Password() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Password", reflect.TypeOf((*MockUser)(nil).Password))
}

// Salt mocks base method.
func (m *MockUser) Salt() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Salt")
	ret0, _ := ret[0].(string)
	return ret0
}

// Salt indicates an expected call of Salt.
func (mr *MockUserMockRecorder) Salt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Salt", reflect.TypeOf((*MockUser)(nil).Salt))
}

// VerifyPassword mocks base method.
func (m *MockUser) VerifyPassword(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPassword", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyPassword indicates an expected call of VerifyPassword.
func (mr *MockUserMockRecorder) VerifyPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPassword", reflect.TypeOf((*MockUser)(nil).VerifyPassword), arg0)
}