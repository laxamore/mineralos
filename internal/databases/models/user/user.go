package user

import "time"

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"uniqueIndex" json:"username"`
	Password  string `gorm:"not null" json:"password"`
	Email     string `gorm:"uniqueIndex" json:"email"`
	Role      Role   `gorm:"not null" json:"role"`
}
