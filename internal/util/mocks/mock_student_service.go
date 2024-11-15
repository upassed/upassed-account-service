// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/student/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

// StudentService is a mock of Service interface.
type StudentService struct {
	ctrl     *gomock.Controller
	recorder *StudentServiceMockRecorder
}

// StudentServiceMockRecorder is the mock recorder for StudentService.
type StudentServiceMockRecorder struct {
	mock *StudentService
}

// NewStudentService creates a new mock instance.
func NewStudentService(ctrl *gomock.Controller) *StudentService {
	mock := &StudentService{ctrl: ctrl}
	mock.recorder = &StudentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *StudentService) EXPECT() *StudentServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *StudentService) Create(ctx context.Context, student *business.Student) (*business.StudentCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, student)
	ret0, _ := ret[0].(*business.StudentCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *StudentServiceMockRecorder) Create(ctx, student interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*StudentService)(nil).Create), ctx, student)
}

// FindByUsername mocks base method.
func (m *StudentService) FindByUsername(ctx context.Context, studentUsername string) (*business.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", ctx, studentUsername)
	ret0, _ := ret[0].(*business.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *StudentServiceMockRecorder) FindByUsername(ctx, studentUsername interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*StudentService)(nil).FindByUsername), ctx, studentUsername)
}

// unusedStudentRepo1 is a mock of repository interface.
type unusedStudentRepo1 struct {
	ctrl     *gomock.Controller
	recorder *unusedStudentRepo1MockRecorder
}

// unusedStudentRepo1MockRecorder is the mock recorder for unusedStudentRepo1.
type unusedStudentRepo1MockRecorder struct {
	mock *unusedStudentRepo1
}

// NewunusedStudentRepo1 creates a new mock instance.
func NewunusedStudentRepo1(ctrl *gomock.Controller) *unusedStudentRepo1 {
	mock := &unusedStudentRepo1{ctrl: ctrl}
	mock.recorder = &unusedStudentRepo1MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *unusedStudentRepo1) EXPECT() *unusedStudentRepo1MockRecorder {
	return m.recorder
}

// CheckDuplicateExists mocks base method.
func (m *unusedStudentRepo1) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicateExists", ctx, educationalEmail, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicateExists indicates an expected call of CheckDuplicateExists.
func (mr *unusedStudentRepo1MockRecorder) CheckDuplicateExists(ctx, educationalEmail, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicateExists", reflect.TypeOf((*unusedStudentRepo1)(nil).CheckDuplicateExists), ctx, educationalEmail, username)
}

// FindByUsername mocks base method.
func (m *unusedStudentRepo1) FindByUsername(arg0 context.Context, arg1 string) (*domain.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*domain.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *unusedStudentRepo1MockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*unusedStudentRepo1)(nil).FindByUsername), arg0, arg1)
}

// Save mocks base method.
func (m *unusedStudentRepo1) Save(arg0 context.Context, arg1 *domain.Student) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *unusedStudentRepo1MockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*unusedStudentRepo1)(nil).Save), arg0, arg1)
}

// unusedGroupRepo2 is a mock of groupRepository interface.
type unusedGroupRepo2 struct {
	ctrl     *gomock.Controller
	recorder *unusedGroupRepo2MockRecorder
}

// unusedGroupRepo2MockRecorder is the mock recorder for unusedGroupRepo2.
type unusedGroupRepo2MockRecorder struct {
	mock *unusedGroupRepo2
}

// NewunusedGroupRepo2 creates a new mock instance.
func NewunusedGroupRepo2(ctrl *gomock.Controller) *unusedGroupRepo2 {
	mock := &unusedGroupRepo2{ctrl: ctrl}
	mock.recorder = &unusedGroupRepo2MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *unusedGroupRepo2) EXPECT() *unusedGroupRepo2MockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *unusedGroupRepo2) Exists(arg0 context.Context, arg1 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *unusedGroupRepo2MockRecorder) Exists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*unusedGroupRepo2)(nil).Exists), arg0, arg1)
}

// FindByID mocks base method.
func (m *unusedGroupRepo2) FindByID(arg0 context.Context, arg1 uuid.UUID) (*domain.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*domain.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *unusedGroupRepo2MockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*unusedGroupRepo2)(nil).FindByID), arg0, arg1)
}
