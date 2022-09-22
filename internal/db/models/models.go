package models

func GetModels() []interface{} {
	return []interface{}{
		&Rig{},
		&User{},
		&Role{},
	}
}
