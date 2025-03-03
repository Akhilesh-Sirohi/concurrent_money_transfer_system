package wallet

import (
	"time"

	"concurrent_money_transfer_system/utils"
)

type WalletStatus string

const (
	Active   WalletStatus = "active"
	Inactive WalletStatus = "inactive"
)

type Wallet struct {
	ID        string         `json:"id"`
	UserID    string         `json:"-"`
	Balance   float64        `json:"balance"`
	Currency  utils.Currency `json:"currency"`
	Status    WalletStatus   `json:"wallet_status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
