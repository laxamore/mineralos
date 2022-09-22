package models

import (
	"time"
)

type Rig struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RigID     string `gorm:"uniqueIndex" json:"rig_id"`
	RigName   string `json:"rig_name"`
}
