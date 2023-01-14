package models

type EmailParam struct {
	Email string `uri:"email" validate:"required,email"`
}

type UpdateParam struct {
	ID    int64  `uri:"id" validate:"required"`
	Email string `uri:"email" validate:"required,email"`
}
