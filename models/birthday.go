package models

import "time"

type Birthday struct {
	UniqueIdentifier string    `json:"uniqueIdentifier"`
	ID               int64     `json:"id"`
	FirstName        string    `json:"firstName" validate:"required,alpha,min=3,max=50"`
	LastName         string    `json:"lastName" validate:"required,alpha,min=1,max=50"`
	DateOfBirth      time.Time `json:"dateOfBirth" validate:"isValidDOB"`
	Age              int       `json:"age"`
	Gender           string    `json:"gender" validate:"required,oneof=MALE FEMALE"`
	ProfilePic       string    `json:"profilePic"`
	Interests        []int64   `json:"interests"`
}
