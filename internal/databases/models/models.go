package models

import "github.com/laxamore/mineralos/internal/databases/models/rig"

func GetModels() []interface{} {
	return []interface{}{
		&rig.Rig{},
	}
}
