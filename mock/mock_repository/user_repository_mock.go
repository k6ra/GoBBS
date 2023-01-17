// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/user_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "GoBBS/domain/model"
	reflect "reflect"
	time "time"

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

// Delete mocks base method.
func (m *MockUser) Delete(user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), user)
}

// FindByEmail mocks base method.
func (m *MockUser) FindByEmail(email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserMockRecorder) FindByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUser)(nil).FindByEmail), email)
}

// Regist mocks base method.
func (m *MockUser) Regist(user *model.User, now time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Regist", user, now)
	ret0, _ := ret[0].(error)
	return ret0
}

// Regist indicates an expected call of Regist.
func (mr *MockUserMockRecorder) Regist(user, now interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Regist", reflect.TypeOf((*MockUser)(nil).Regist), user, now)
}

// Update mocks base method.
func (m *MockUser) Update(user *model.User, now time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user, now)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserMockRecorder) Update(user, now interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), user, now)
}
