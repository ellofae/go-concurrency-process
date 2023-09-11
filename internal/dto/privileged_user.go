package dto

type PrivilegedUserCreateDTO struct {
	UserID    int    `json:"user_id"`
	Privilege string `json:"privilege_title"`
}
