package group

import (
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
)

func ConvertToFindStudentsInGroupResponse(studentsInGroup []business.Student) *client.FindStudentsInGroupResponse {
	response := client.FindStudentsInGroupResponse{}
	convertedStudents := make([]*client.StudentDTO, 0, len(studentsInGroup))
	for _, studentToConvert := range studentsInGroup {
		convertedStudents = append(convertedStudents, convertStudent(studentToConvert))
	}

	response.StudentsInGroup = convertedStudents
	return &response
}

func convertStudent(studentToConvert business.Student) *client.StudentDTO {
	return &client.StudentDTO{
		Id:               studentToConvert.ID.String(),
		FirstName:        studentToConvert.FirstName,
		LastName:         studentToConvert.LastName,
		MiddleName:       studentToConvert.MiddleName,
		EducationalEmail: studentToConvert.EducationalEmail,
		Username:         studentToConvert.Username,
		Group: &client.GroupDTO{
			Id:                 studentToConvert.Group.ID.String(),
			SpecializationCode: studentToConvert.Group.SpecializationCode,
			GroupNumber:        studentToConvert.Group.GroupNumber,
		},
	}
}
