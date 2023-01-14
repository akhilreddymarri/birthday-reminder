package models

type User struct {
	ID              int64      `json:"id"`
	FirstName       string     `json:"firstName" validate:"required,alpha,min=3,max=50"`
	LastName        string     `json:"lastName" validate:"required,alpha,min=1,max=50"`
	Email           string     `json:"email" validate:"required,email"`
	FriendsBirthday []Birthday `json:"friendsBirthday"`
}
