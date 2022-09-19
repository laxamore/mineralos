package user

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
		RoleName: "user",
		RoleDesc: "User",
		Level:    9,
	}
)

type Role struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RoleName  string `gorm:"unique" json:"role_name"`
	RoleDesc  string `gorm:"not null" json:"role_desc"`
	Level     uint8  `gorm:"not null" json:"level"`
}
