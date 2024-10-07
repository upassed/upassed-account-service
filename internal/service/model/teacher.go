package business

import "github.com/google/uuid"

type TeacherCreateRequest struct {
}

type TeacherCreateResponse struct {
	CreatedTeacherID uuid.UUID
}
