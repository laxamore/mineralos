package models

import "time"

// Default Role
var (
	RoleAdmin = Role{
		RoleName: "admin",
		RoleDesc: "Admin",
		Level:    0,
	}
	RoleOperator = Role{
		RoleName: "operator",
		RoleDesc: "Operator",
		Level:    1,
	}
	RoleUser = Role{
		RoleName: "users",
		RoleDesc: "User",
		Level:    9,
	}
)

type Role struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoleName  string    `gorm:"uniqueIndex" json:"role_name"`
	RoleDesc  string    `json:"role_desc"`
	Level     uint8     `json:"level"`
}
