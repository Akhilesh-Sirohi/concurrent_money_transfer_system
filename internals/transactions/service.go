package transactions

import (
	"concurrent_money_transfer_system/internals/wallet"
	"concurrent_money_transfer_system/utils"
	"context"
	"time"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transferRequest *TransferRequest) (Transaction, error)
	GetTransaction(ctx context.Context, id string) (Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error)
	GetAllTransactions(ctx context.Context) ([]Transaction, error)
}

type transactionService struct {
	repo          TransactionRepo
	walletService wallet.WalletService
}

func (s *transactionService) getSenderAndReceiverWalletsWithLockingOrder(ctx context.Context, transferRequest *TransferRequest) (senderWallet wallet.Wallet, receiverWallet wallet.Wallet, err error) {
	// GetWalletForUpdate takes a lock on wallet_id, so always fetching the wallet with
	// minimum ID first, we are making sure of not falling into deadlock.
	// This prevents a situation where two concurrent transactions could lock
	// wallets in opposite order, causing a deadlock.
	if transferRequest.SenderID < transferRequest.ReceiverID {
		senderWallet, err = s.walletService.GetWalletForUpdate(ctx, transferRequest.SenderID)
		if err != nil {
			return
		}
		receiverWallet, err = s.walletService.GetWalletForUpdate(ctx, transferRequest.ReceiverID)
	} else {
		receiverWallet, err = s.walletService.GetWalletForUpdate(ctx, transferRequest.ReceiverID)
		if err != nil {
			return
		}
		senderWallet, err = s.walletService.GetWalletForUpdate(ctx, transferRequest.SenderID)
	}

	return senderWallet, receiverWallet, err
}

func (s *transactionService) CreateTransaction(ctx context.Context, transferRequest *TransferRequest) (Transaction, error) {
	if transferRequest.SenderID == transferRequest.ReceiverID {
		return Transaction{}, utils.NewError(utils.ErrTransactionSameUser)
	}

	defer s.walletService.ReleaseGetWalletForUpdateLock(ctx, transferRequest.SenderID)
	defer s.walletService.ReleaseGetWalletForUpdateLock(ctx, transferRequest.ReceiverID)

	senderWallet, receiverWallet, err := s.getSenderAndReceiverWalletsWithLockingOrder(ctx, transferRequest)

	if err != nil {
		return Transaction{}, err
	}

	if senderWallet.Balance < transferRequest.Amount {
		return Transaction{}, utils.NewError(utils.ErrInsufficientBalance)
	}

	if senderWallet.Status == wallet.Inactive || receiverWallet.Status == wallet.Inactive {
		return Transaction{}, utils.NewError(utils.ErrWalletInactive)
	}

	transaction := Transaction{
		ID:              utils.GenerateUniqueEntityId(),
		DebitUserID:     transferRequest.SenderID,
		CreditUserID:    transferRequest.ReceiverID,
		Amount:          transferRequest.Amount,
		Currency:        transferRequest.Currency,
		Status:          Pending,
		TransactionType: Transfer,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Description:     transferRequest.Description,
		PaymentDetails:  transferRequest.PaymentDetails,
	}

	newSenderBalance := senderWallet.Balance - transferRequest.Amount
	newReceiverBalance := receiverWallet.Balance + transferRequest.Amount
	s.repo.CreateTransaction(ctx, transaction)

	err = s.walletService.UpdateWalletBalance(ctx, senderWallet.ID, newSenderBalance)
	if err != nil {
		transaction, err = s.repo.UpdateTransactionStatus(ctx, transaction.ID, Failed)
		s.walletService.UpdateWalletBalance(ctx, senderWallet.ID, senderWallet.Balance) // make the balance as it is
		return transaction, err
	}

	err = s.walletService.UpdateWalletBalance(ctx, receiverWallet.ID, newReceiverBalance)
	if err != nil {
		transaction, err = s.repo.UpdateTransactionStatus(ctx, transaction.ID, Failed)
		s.walletService.UpdateWalletBalance(ctx, senderWallet.ID, senderWallet.Balance) // make the balance as it is
		return Transaction{}, err
	}

	transaction, err = s.repo.UpdateTransactionStatus(ctx, transaction.ID, Completed)

	return transaction, err
}

func (s *transactionService) GetTransaction(ctx context.Context, id string) (Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *transactionService) GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
	return s.repo.GetTransactionsByUserID(ctx, userID)
}

func (s *transactionService) GetAllTransactions(ctx context.Context) ([]Transaction, error) {
	return s.repo.GetAllTransactions(ctx)
}

var transactionServiceInstance *transactionService

func NewTransactionService(repo TransactionRepo, walletService wallet.WalletService) TransactionService {
	if transactionServiceInstance == nil {
		transactionServiceInstance = &transactionService{repo: repo, walletService: walletService}
	}
	return transactionServiceInstance
}
