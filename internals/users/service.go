package users

import (
	"concurrent_money_transfer_system/internals/wallet"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]User, error)
}

type userService struct {
	userRepo      UserRepo
	walletService wallet.WalletService
}

func (s *userService) CreateUser(ctx context.Context, user User) (User, error) {
	user, err := s.userRepo.CreateUser(user)
	if err != nil {
		return User{}, err
	}
	wallet, err := s.walletService.CreateWallet(ctx, user.ID, user.Wallet.Balance)
	if err != nil {
		return User{}, err
	}
	user.Wallet = wallet
	user.Password = ""

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	for i, user := range users {
		user.Wallet, _ = s.walletService.GetWallet(ctx, user.ID)
		user.Password = ""
		users[i] = user
	}
	return users, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (User, error) {
	user, err := s.userRepo.GetUser(id)
	if err != nil {
		return User{}, err
	}

	wallet, err := s.walletService.GetWallet(ctx, id)
	if err != nil {
		return User{}, err
	}

	user.Wallet = wallet
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user User) (User, error) {
	user, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return User{}, err
	}
	user.Wallet = wallet.Wallet{}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	err = s.walletService.DisableWallet(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

var userServiceInstance *userService

func NewUserService(userRepo UserRepo, walletService wallet.WalletService) UserService {
	if userServiceInstance == nil {
		userServiceInstance = &userService{userRepo: userRepo, walletService: walletService}
	}
	return userServiceInstance
}
