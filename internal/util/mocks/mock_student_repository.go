// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/student/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
)

// StudentRepository is a mock of Repository interface.
type StudentRepository struct {
	ctrl     *gomock.Controller
	recorder *StudentRepositoryMockRecorder
}

// StudentRepositoryMockRecorder is the mock recorder for StudentRepository.
type StudentRepositoryMockRecorder struct {
	mock *StudentRepository
}

// NewStudentRepository creates a new mock instance.
func NewStudentRepository(ctrl *gomock.Controller) *StudentRepository {
	mock := &StudentRepository{ctrl: ctrl}
	mock.recorder = &StudentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *StudentRepository) EXPECT() *StudentRepositoryMockRecorder {
	return m.recorder
}

// CheckDuplicateExists mocks base method.
func (m *StudentRepository) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicateExists", ctx, educationalEmail, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicateExists indicates an expected call of CheckDuplicateExists.
func (mr *StudentRepositoryMockRecorder) CheckDuplicateExists(ctx, educationalEmail, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicateExists", reflect.TypeOf((*StudentRepository)(nil).CheckDuplicateExists), ctx, educationalEmail, username)
}

// FindByUsername mocks base method.
func (m *StudentRepository) FindByUsername(arg0 context.Context, arg1 string) (*domain.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", arg0, arg1)
	ret0, _ := ret[0].(*domain.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *StudentRepositoryMockRecorder) FindByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*StudentRepository)(nil).FindByUsername), arg0, arg1)
}

// Save mocks base method.
func (m *StudentRepository) Save(arg0 context.Context, arg1 *domain.Student) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *StudentRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*StudentRepository)(nil).Save), arg0, arg1)
}
