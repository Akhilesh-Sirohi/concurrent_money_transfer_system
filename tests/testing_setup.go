package tests

import (
	"context"
	"log"

	"concurrent_money_transfer_system/internals/users"
	"concurrent_money_transfer_system/internals/wallet"
)

var userService users.UserService

func CreateTestUsers() {
	log.Println("Creating test users...")
	userRepo := users.NewUserRepo()
	walletRepo := wallet.NewWalletRepo()
	userService = users.NewUserService(userRepo, wallet.NewWalletService(walletRepo))

	ctx := context.Background()
	userService.CreateUser(ctx, users.User{
		ID:        "1",
		FirstName: "Mark",
		Email:     "mark@facebook.com",
		PhoneNumber: "+1234567890",
		Password:  "password",
		Wallet: wallet.Wallet{
			Balance: 100,
		},
	})

	userService.CreateUser(ctx, users.User{
		ID:          "2",
		FirstName:   "Jane",
		Email:       "jane@gmail.com",
		PhoneNumber: "+1234567890",
		Password:    "password",
		Wallet: wallet.Wallet{
			Balance: 50,
		},
	})

	userService.CreateUser(ctx, users.User{
		ID:          "3",
		FirstName:   "Adam",
		Email:       "adam@gmail.com",
		PhoneNumber: "+1234567890",
		Password:    "password",
		Wallet: wallet.Wallet{
			Balance: 0,
		},
	})

}
