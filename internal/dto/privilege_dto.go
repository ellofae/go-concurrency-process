package dto

type PrivilegeCreateDTO struct {
	PrivilegeTitle string `json:"privilege_title", validate:"required,lte=20"`
}
