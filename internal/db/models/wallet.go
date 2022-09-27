package models

import "time"

type Wallet struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	WalletCoin    string `json:"wallet_coin"`
	WalletName    string `json:"wallet_name"`
	WalletAddress string `json:"wallet_address"`
}
