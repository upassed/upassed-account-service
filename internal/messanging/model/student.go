package event

import (
	"github.com/go-playground/validator/v10"
)

type StudentCreateRequest struct {
	FirstName        string `json:"first_name,omitempty" validate:"required,min=4,max=30"`
	LastName         string `json:"last_name,omitempty" validate:"required,min=4,max=30"`
	MiddleName       string `json:"middle_name,omitempty" validate:"max=30"`
	EducationalEmail string `json:"educational_email,omitempty" validate:"required,email"`
	Username         string `json:"username,omitempty" validate:"required,min=4,max=30,username"`
	GroupId          string `json:"group_id,omitempty" validate:"required,uuid"`
}

func (request *StudentCreateRequest) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("uuid", ValidateUUID())
	_ = validate.RegisterValidation("username", ValidateUsername())

	if err := validate.Struct(*request); err != nil {
		return err
	}

	return nil
}

type StudentCreateResponse struct {
	CreatedStudentID string `json:"created_student_id,omitempty"`
}
