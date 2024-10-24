package domain

import "github.com/google/uuid"

type Group struct {
	ID                 uuid.UUID `json:"id"`
	SpecializationCode string    `json:"specialization_code"`
	GroupNumber        string    `json:"group_number"`
}

func (Group) TableName() string {
	return "group"
}

type GroupFilter struct {
	SpecializationCode string
	GroupNumber        string
}
