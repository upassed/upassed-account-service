package domain

import "github.com/google/uuid"

type Group struct {
	ID                 uuid.UUID
	SpecializationCode string
	GroupNumber        string
}

func (Group) TableName() string {
	return "group"
}

type GroupFilter struct {
	SpecializationCode string
	GroupNumber        string
}
