package entity

import "time"

type PrivilegedUser struct {
	UserID      int       `json:"user_id" db:"user_id"`
	PrivilegeID int       `json:"privilege_id" db:"privilege_id"`
	AssignedAt  time.Time `json:"assigned_at" db:"assigned_at"`
}
