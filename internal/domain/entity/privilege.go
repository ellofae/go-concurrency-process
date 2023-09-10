package entity

import "time"

type Privilege struct {
	ID             int       `json:"id" db:"id"`
	PrivilegeTitle string    `json:"privilege_title" db:"privilege_title"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
