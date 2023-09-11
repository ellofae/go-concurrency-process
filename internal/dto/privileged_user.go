package dto

type PrivilegedUserDTO struct {
	UserID    int    `json:"user_id"`
	Privilege string `json:"privilege_title" validate:"required,uppercase,max=20"`
}

type PrivilegedUserUdateDTO struct {
	Privilege string `json:"privilege_title" validate:"required,uppercase,max=20"`
}
