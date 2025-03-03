package transactions

import (
	"context"
	"sync"
	"time"

	"concurrent_money_transfer_system/utils"
)

type TransactionRepo interface {
	CreateTransaction(ctx context.Context, transaction Transaction) (Transaction, error)
	GetTransaction(ctx context.Context, id string) (Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status TransactionStatus) (Transaction, error)
	GetAllTransactions(ctx context.Context) ([]Transaction, error)
}

type transactionRepo struct {
	transactions sync.Map
}

func (r *transactionRepo) CreateTransaction(ctx context.Context, transaction Transaction) (Transaction, error) {
	if transaction.ID == "" {
		transaction.ID = utils.GenerateUniqueEntityId()
	}
	r.transactions.Store(transaction.ID, transaction)
	return transaction, nil
}

func (r *transactionRepo) GetTransaction(ctx context.Context, id string) (Transaction, error) {
	transaction, ok := r.transactions.Load(id)
	if !ok {
		return Transaction{}, utils.NewError(utils.ErrTransactionNotFound)
	}
	return transaction.(Transaction), nil
}

func (r *transactionRepo) UpdateTransactionStatus(ctx context.Context, id string, status TransactionStatus) (Transaction, error) {
	transaction, err := r.GetTransaction(ctx, id)
	if err != nil {
		return Transaction{}, err
	}
	transaction.Status = status
	transaction.UpdatedAt = time.Now()
	r.transactions.Store(id, transaction)
	return transaction, nil
}

func (r *transactionRepo) GetTransactionsByUserID(ctx context.Context, userID string) ([]Transaction, error) {
	transactions := make([]Transaction, 0)
	r.transactions.Range(func(key, value any) bool {
		if value.(Transaction).DebitUserID == userID || value.(Transaction).CreditUserID == userID {
			transactions = append(transactions, value.(Transaction))
		}
		return true
	})
	return transactions, nil
}

func (r *transactionRepo) GetAllTransactions(ctx context.Context) ([]Transaction, error) {
	transactions := make([]Transaction, 0)
	r.transactions.Range(func(key, value any) bool {
		transactions = append(transactions, value.(Transaction))
		return true
	})
	return transactions, nil
}

var transactionRepoInstance *transactionRepo

func NewTransactionRepo() TransactionRepo {
	if transactionRepoInstance == nil {
		transactionRepoInstance = &transactionRepo{
			transactions: sync.Map{},
		}	
	}
	return transactionRepoInstance
}

// Reset is just used for testing purposes
func Reset() {
	transactionRepoInstance.transactions.Range(func(key, value any) bool {
		transactionRepoInstance.transactions.Delete(key)
		return true
	})
}
