package dto

type PrivilegeDTO struct {
	PrivilegeTitle string `json:"privilege_title" validate:"required,uppercase,max=20"`
}

type PrivilegeResponseDTO struct {
	ID             int    `json:"id"`
	PrivilegeTitle string `json:"privilege_title"`
}

type PrivilegeUpdateDTO struct {
	PrivilegeTitle string `json:"privilege_title" validate:"required,uppercase,max=20"`
}
