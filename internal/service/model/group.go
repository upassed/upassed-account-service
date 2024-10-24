package business

import "github.com/google/uuid"

type Group struct {
	ID                 uuid.UUID
	SpecializationCode string
	GroupNumber        string
}

type GroupFilter struct {
	SpecializationCode string
	GroupNumber        string
}
