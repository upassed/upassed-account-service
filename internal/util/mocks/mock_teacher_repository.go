// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/teacher/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
)

// TeacherRepository is a mock of Repository interface.
type TeacherRepository struct {
	ctrl     *gomock.Controller
	recorder *TeacherRepositoryMockRecorder
}

// TeacherRepositoryMockRecorder is the mock recorder for TeacherRepository.
type TeacherRepositoryMockRecorder struct {
	mock *TeacherRepository
}

// NewTeacherRepository creates a new mock instance.
func NewTeacherRepository(ctrl *gomock.Controller) *TeacherRepository {
	mock := &TeacherRepository{ctrl: ctrl}
	mock.recorder = &TeacherRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *TeacherRepository) EXPECT() *TeacherRepositoryMockRecorder {
	return m.recorder
}

// CheckDuplicateExists mocks base method.
func (m *TeacherRepository) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicateExists", ctx, reportEmail, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicateExists indicates an expected call of CheckDuplicateExists.
func (mr *TeacherRepositoryMockRecorder) CheckDuplicateExists(ctx, reportEmail, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicateExists", reflect.TypeOf((*TeacherRepository)(nil).CheckDuplicateExists), ctx, reportEmail, username)
}

// FindByUsername mocks base method.
func (m *TeacherRepository) FindByUsername(arg0 context.Context, arg1 string) (*domain.Teacher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*domain.Teacher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *TeacherRepositoryMockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*TeacherRepository)(nil).FindByUsername), arg0, arg1)
}

// Save mocks base method.
func (m *TeacherRepository) Save(arg0 context.Context, arg1 *domain.Teacher) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *TeacherRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*TeacherRepository)(nil).Save), arg0, arg1)
}