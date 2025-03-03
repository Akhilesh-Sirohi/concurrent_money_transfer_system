package wallet

import (
	"concurrent_money_transfer_system/utils"
	"context"
	"sync"
	"time"
)

type WalletRepo interface {
	CreateWallet(ctx context.Context, wallet Wallet) (Wallet, error)
	DisableWallet(ctx context.Context, userID string) error
	GetWalletByUserID(ctx context.Context, userID string) (Wallet, error)
	GetWalletForUpdateByUserID(ctx context.Context, userID string) (Wallet, error)
	ReleaseGetWalletForUpdateLock(ctx context.Context, userID string)
	UpdateWalletBalance(ctx context.Context, walletID string, newBalance float64) error
}

type walletRepo struct {
	wallets       sync.Map // Using sync.Map for concurrent safe map[string]Wallet operations
	walletMutexes sync.Map // Mutex for each wallet to avoid race conditions
}

func (r *walletRepo) CreateWallet(ctx context.Context, wallet Wallet) (Wallet, error) {
	wallet.ID = wallet.UserID // Kept same as userID for simplicity
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	r.wallets.Store(wallet.ID, wallet)
	r.walletMutexes.Store(wallet.ID, &sync.Mutex{})
	return wallet, nil
}

func (r *walletRepo) DisableWallet(ctx context.Context, userID string) error {
	wallet, err := r.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}
	walletMutex, ok := r.walletMutexes.Load(userID)
	if !ok {
		return utils.NewErrorWithMessage(utils.ErrWalletNotFound, "Wallet Mutex not found for userID: "+userID)
	}
	walletMutex.(*sync.Mutex).Lock()
	defer walletMutex.(*sync.Mutex).Unlock()
	wallet.Status = Inactive
	wallet.UpdatedAt = time.Now()
	r.wallets.Store(userID, wallet)
	return nil
}

func (r *walletRepo) GetWalletByUserID(ctx context.Context, userID string) (Wallet, error) {
	wallet, ok := r.wallets.Load(userID)
	if !ok {
		return Wallet{}, utils.NewErrorWithMessage(utils.ErrWalletNotFound, "Wallet not found for userID: "+userID)
	}
	return wallet.(Wallet), nil
}

func (r *walletRepo) GetWalletForUpdateByUserID(ctx context.Context, userID string) (Wallet, error) {
	// This method acquires an exclusive lock on the wallet, similar to
	// SELECT ... FOR UPDATE in MySQL. The lock prevents concurrent modifications
	// to the same wallet and must be released by calling UpdateWallet.
	walletMutex, ok := r.walletMutexes.Load(userID) // Mutex is locked first, otherwise 2 transactions can get the same walletBalance
	if !ok {
		return Wallet{}, utils.NewErrorWithMessage(utils.ErrWalletNotFound, "Wallet not found for userID: "+userID)
	}
	walletMutex.(*sync.Mutex).Lock()

	wallet, err := r.GetWalletByUserID(ctx, userID)
	if err != nil {
		return Wallet{}, utils.NewErrorWithMessage(utils.ErrInternalServerError, "Mutex found but wallet not found for userID: "+userID)
	}

	return wallet, nil
}

func (r *walletRepo) ReleaseGetWalletForUpdateLock(ctx context.Context, userID string) {
	// Check if the mutex exists before unlocking to avoid "unlock of unlocked mutex" panic
	walletMutex, ok := r.walletMutexes.Load(userID)
	if ok {
		walletMutex.(*sync.Mutex).Unlock()
	}
}

func (r *walletRepo) UpdateWalletBalance(ctx context.Context, walletID string, newBalance float64) error {
	wallet, ok := r.wallets.Load(walletID)
	if !ok {
		return utils.NewErrorWithMessage(utils.ErrInternalServerError, "Wallet not found for walletID: "+walletID)
	}
	wal := wallet.(Wallet)
	wal.Balance = newBalance
	wal.UpdatedAt = time.Now()
	r.wallets.Store(walletID, wal)
	return nil
}

var walletRepoInstance *walletRepo

func NewWalletRepo() WalletRepo {
	if walletRepoInstance == nil {
		walletRepoInstance = &walletRepo{
			wallets:       sync.Map{},
			walletMutexes: sync.Map{},
		}
	}
	return walletRepoInstance
}
