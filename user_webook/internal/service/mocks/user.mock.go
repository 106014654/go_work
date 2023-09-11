// Code generated by MockGen. DO NOT EDIT.
// Source: user_webook/internal/service/user.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "go_work/user_webook/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserServiceInter is a mock of UserServiceInter interface.
type MockUserServiceInter struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceInterMockRecorder
}

// MockUserServiceInterMockRecorder is the mock recorder for MockUserServiceInter.
type MockUserServiceInterMockRecorder struct {
	mock *MockUserServiceInter
}

// NewMockUserServiceInter creates a new mock instance.
func NewMockUserServiceInter(ctrl *gomock.Controller) *MockUserServiceInter {
	mock := &MockUserServiceInter{ctrl: ctrl}
	mock.recorder = &MockUserServiceInterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceInter) EXPECT() *MockUserServiceInterMockRecorder {
	return m.recorder
}

// EditUserDetail mocks base method.
func (m *MockUserServiceInter) EditUserDetail(ctx context.Context, id int64, name, birth, intro string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditUserDetail", ctx, id, name, birth, intro)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditUserDetail indicates an expected call of EditUserDetail.
func (mr *MockUserServiceInterMockRecorder) EditUserDetail(ctx, id, name, birth, intro interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditUserDetail", reflect.TypeOf((*MockUserServiceInter)(nil).EditUserDetail), ctx, id, name, birth, intro)
}

// Login mocks base method.
func (m *MockUserServiceInter) Login(ctx context.Context, email, password string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceInterMockRecorder) Login(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserServiceInter)(nil).Login), ctx, email, password)
}

// Profile mocks base method.
func (m *MockUserServiceInter) Profile(ctx context.Context, id int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockUserServiceInterMockRecorder) Profile(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockUserServiceInter)(nil).Profile), ctx, id)
}

// SignUp mocks base method.
func (m *MockUserServiceInter) SignUp(ctx context.Context, u domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceInterMockRecorder) SignUp(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserServiceInter)(nil).SignUp), ctx, u)
}
