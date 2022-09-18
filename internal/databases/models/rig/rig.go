package rig

import "gorm.io/gorm"

type Rig struct {
	gorm.Model
	RigID   string `gorm:"uniqueIndex" json:"rig_id"`
	RigName string `json:"rig_name"`
}
