// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/teacher/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

// TeacherService is a mock of Service interface.
type TeacherService struct {
	ctrl     *gomock.Controller
	recorder *TeacherServiceMockRecorder
}

// TeacherServiceMockRecorder is the mock recorder for TeacherService.
type TeacherServiceMockRecorder struct {
	mock *TeacherService
}

// NewTeacherService creates a new mock instance.
func NewTeacherService(ctrl *gomock.Controller) *TeacherService {
	mock := &TeacherService{ctrl: ctrl}
	mock.recorder = &TeacherServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *TeacherService) EXPECT() *TeacherServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *TeacherService) Create(ctx context.Context, teacher *business.Teacher) (*business.TeacherCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, teacher)
	ret0, _ := ret[0].(*business.TeacherCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *TeacherServiceMockRecorder) Create(ctx, teacher interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*TeacherService)(nil).Create), ctx, teacher)
}

// FindByUsername mocks base method.
func (m *TeacherService) FindByUsername(ctx context.Context, teacherUsername string) (*business.Teacher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", ctx, teacherUsername)
	ret0, _ := ret[0].(*business.Teacher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *TeacherServiceMockRecorder) FindByUsername(ctx, teacherUsername interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*TeacherService)(nil).FindByUsername), ctx, teacherUsername)
}

// unusedTeacherRepo1 is a mock of repository interface.
type unusedTeacherRepo1 struct {
	ctrl     *gomock.Controller
	recorder *unusedTeacherRepo1MockRecorder
}

// unusedTeacherRepo1MockRecorder is the mock recorder for unusedTeacherRepo1.
type unusedTeacherRepo1MockRecorder struct {
	mock *unusedTeacherRepo1
}

// NewunusedTeacherRepo1 creates a new mock instance.
func NewunusedTeacherRepo1(ctrl *gomock.Controller) *unusedTeacherRepo1 {
	mock := &unusedTeacherRepo1{ctrl: ctrl}
	mock.recorder = &unusedTeacherRepo1MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *unusedTeacherRepo1) EXPECT() *unusedTeacherRepo1MockRecorder {
	return m.recorder
}

// CheckDuplicateExists mocks base method.
func (m *unusedTeacherRepo1) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicateExists", ctx, reportEmail, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicateExists indicates an expected call of CheckDuplicateExists.
func (mr *unusedTeacherRepo1MockRecorder) CheckDuplicateExists(ctx, reportEmail, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicateExists", reflect.TypeOf((*unusedTeacherRepo1)(nil).CheckDuplicateExists), ctx, reportEmail, username)
}

// FindByUsername mocks base method.
func (m *unusedTeacherRepo1) FindByUsername(arg0 context.Context, arg1 string) (*domain.Teacher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*domain.Teacher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *unusedTeacherRepo1MockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*unusedTeacherRepo1)(nil).FindByUsername), arg0, arg1)
}

// Save mocks base method.
func (m *unusedTeacherRepo1) Save(arg0 context.Context, arg1 *domain.Teacher) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *unusedTeacherRepo1MockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*unusedTeacherRepo1)(nil).Save), arg0, arg1)
}