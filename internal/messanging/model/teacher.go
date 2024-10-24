package event

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TeacherCreateRequest struct {
	FirstName   string `json:"first_name,omitempty" validate:"required,min=4,max=30"`
	LastName    string `json:"last_name,omitempty" validate:"required,min=4,max=30"`
	MiddleName  string `json:"middle_name,omitempty" validate:"max=30"`
	ReportEmail string `json:"report_email,omitempty" validate:"required,email"`
	Username    string `json:"username,omitempty" validate:"required,min=4,max=30"`
}

func (request *TeacherCreateRequest) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		_, err := uuid.Parse(fl.Field().String())
		return err == nil
	})

	if err := validate.Struct(*request); err != nil {
		return err
	}

	return nil
}

type TeacherCreateResponse struct {
	CreatedTeacherID string `json:"created_teacher_id,omitempty"`
}
