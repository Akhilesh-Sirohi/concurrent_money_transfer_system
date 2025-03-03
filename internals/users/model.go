package users

import (
	"concurrent_money_transfer_system/internals/wallet"
	"time"
)

type User struct {
	wallet.Wallet
	ID          string     `json:"id"`
	FirstName   string     `json:"first_name" validate:"required"`
	LastName    string     `json:"last_name,omitempty"`
	PhoneNumber string     `json:"phone_number" validate:"required,e164"`
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password,omitempty" validate:"required,min=4"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}
