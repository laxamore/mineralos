package models

type RegisterToken struct {
	Token string `json:"token"`
	Role  Role   `json:"role"`
}
