package group

import "github.com/google/uuid"

type Group struct {
	ID                 uuid.UUID
	SpecializationCode string
	GroupNumber        string
}

func (Group) TableName() string {
	return "group"
}
