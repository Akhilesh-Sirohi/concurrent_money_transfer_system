package wallet

import (
	"concurrent_money_transfer_system/utils"
	"context"
)

type WalletService interface {
	CreateWallet(ctx context.Context, userID string, initialBalance float64) (Wallet, error)
	DisableWallet(ctx context.Context, userID string) error
	GetWallet(ctx context.Context, userID string) (Wallet, error)
	GetWalletForUpdate(ctx context.Context, userID string) (Wallet, error)
	ReleaseGetWalletForUpdateLock(ctx context.Context, userID string)
	UpdateWalletBalance(ctx context.Context, walletID string, newBalance float64) error
}

type walletService struct {
	repo WalletRepo
}

func (s *walletService) CreateWallet(ctx context.Context, userID string, initialBalance float64) (Wallet, error) {
	wallet := Wallet{
		ID:       userID,
		UserID:   userID,
		Balance:  initialBalance,
		Currency: utils.USD,
		Status:   Active,
	}
	return s.repo.CreateWallet(ctx, wallet)
}

func (s *walletService) DisableWallet(ctx context.Context, userID string) error {
	return s.repo.DisableWallet(ctx, userID)
}

func (s *walletService) GetWallet(ctx context.Context, userID string) (Wallet, error) {
	return s.repo.GetWalletByUserID(ctx, userID)
}

func (s *walletService) GetWalletForUpdate(ctx context.Context, userID string) (Wallet, error) {
	// This method acquires an exclusive lock on the wallet, similar to
	// SELECT ... FOR UPDATE in MySQL. The lock prevents concurrent modifications
	// to the same wallet and must be released by calling UpdateWallet. // TODO: but what to do if no update?
	return s.repo.GetWalletForUpdateByUserID(ctx, userID)
}

func (s *walletService) ReleaseGetWalletForUpdateLock(ctx context.Context, userID string) {
	s.repo.ReleaseGetWalletForUpdateLock(ctx, userID)
}

func (s *walletService) UpdateWalletBalance(ctx context.Context, walletID string, newBalance float64) error {
	err := s.repo.UpdateWalletBalance(ctx, walletID, newBalance)
	if err != nil {
		return err
	}
	return nil
}

var walletServiceInstance *walletService

func NewWalletService(repo WalletRepo) WalletService {
	if walletServiceInstance == nil {
		walletServiceInstance = &walletService{repo: repo}
	}
	return walletServiceInstance
}
