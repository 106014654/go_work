// Code generated by MockGen. DO NOT EDIT.
// Source: user_webook/internal/repository/dao/user.go

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	dao "go_work/user_webook/internal/repository/dao"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserDAOInter is a mock of UserDAOInter interface.
type MockUserDAOInter struct {
	ctrl     *gomock.Controller
	recorder *MockUserDAOInterMockRecorder
}

// MockUserDAOInterMockRecorder is the mock recorder for MockUserDAOInter.
type MockUserDAOInterMockRecorder struct {
	mock *MockUserDAOInter
}

// NewMockUserDAOInter creates a new mock instance.
func NewMockUserDAOInter(ctrl *gomock.Controller) *MockUserDAOInter {
	mock := &MockUserDAOInter{ctrl: ctrl}
	mock.recorder = &MockUserDAOInterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDAOInter) EXPECT() *MockUserDAOInterMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserDAOInter) FindByEmail(ctx context.Context, email string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserDAOInterMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserDAOInter)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockUserDAOInter) FindById(ctx context.Context, id int64) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUserDAOInterMockRecorder) FindById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUserDAOInter)(nil).FindById), ctx, id)
}

// Insert mocks base method.
func (m *MockUserDAOInter) Insert(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUserDAOInterMockRecorder) Insert(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDAOInter)(nil).Insert), ctx, u)
}

// Update mocks base method.
func (m *MockUserDAOInter) Update(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserDAOInterMockRecorder) Update(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserDAOInter)(nil).Update), ctx, u)
}