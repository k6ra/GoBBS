// Code generated by MockGen. DO NOT EDIT.
// Source: domain/service/user_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "GoBBS/domain/model"
	repository "GoBBS/domain/repository"
	service "GoBBS/domain/service"
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

// Authorize mocks base method.
func (m *MockUser) Authorize(email, password string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", email, password)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockUserMockRecorder) Authorize(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockUser)(nil).Authorize), email, password)
}

// Delete mocks base method.
func (m *MockUser) Delete(user model.User) error {
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

// IsDuplicate mocks base method.
func (m *MockUser) IsDuplicate(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDuplicate", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsDuplicate indicates an expected call of IsDuplicate.
func (mr *MockUserMockRecorder) IsDuplicate(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDuplicate", reflect.TypeOf((*MockUser)(nil).IsDuplicate), email)
}

// Regist mocks base method.
func (m *MockUser) Regist(user model.User, now time.Time) error {
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
func (m *MockUser) Update(user model.User, now time.Time) error {
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

// MockUserFactory is a mock of UserFactory interface.
type MockUserFactory struct {
	ctrl     *gomock.Controller
	recorder *MockUserFactoryMockRecorder
}

// MockUserFactoryMockRecorder is the mock recorder for MockUserFactory.
type MockUserFactoryMockRecorder struct {
	mock *MockUserFactory
}

// NewMockUserFactory creates a new mock instance.
func NewMockUserFactory(ctrl *gomock.Controller) *MockUserFactory {
	mock := &MockUserFactory{ctrl: ctrl}
	mock.recorder = &MockUserFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserFactory) EXPECT() *MockUserFactoryMockRecorder {
	return m.recorder
}

// NewUserService mocks base method.
func (m *MockUserFactory) NewUserService(repo repository.User) service.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUserService", repo)
	ret0, _ := ret[0].(service.User)
	return ret0
}

// NewUserService indicates an expected call of NewUserService.
func (mr *MockUserFactoryMockRecorder) NewUserService(repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUserService", reflect.TypeOf((*MockUserFactory)(nil).NewUserService), repo)
}
