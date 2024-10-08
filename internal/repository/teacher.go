package repository

import (
	"context"

	domain "github.com/upassed/upassed-account-service/internal/repository/model"
)

type TeacherRepositoryImpl struct {
}

func NewTeacherRepository() *TeacherRepositoryImpl {
	return &TeacherRepositoryImpl{}
}

func (repository *TeacherRepositoryImpl) Save(context.Context, domain.Teacher) error {
	panic("not implemented!")
}
