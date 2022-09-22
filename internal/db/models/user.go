package models

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"uniqueIndex" json:"username"`
	Password  string `json:"password"`
	Email     string `gorm:"uniqueIndex" json:"email"`
	RoleID    string `json:"role_id"`
	Role      Role   `gorm:"references:RoleName" json:"role"`
}
