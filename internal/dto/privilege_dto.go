package dto

type PrivilegeCreateDTO struct {
	PrivilegeTitle string `json:"privilege_title" validate:"required,uppercase,max=20"`
}

type PrivilegeUpdateDTO struct {
	PrivilegeTitle string `json:"privilege_title" validate:"required,uppercase,max=20"`
}
