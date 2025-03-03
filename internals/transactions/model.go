package transactions

import (
	"concurrent_money_transfer_system/utils"
	"time"
)

type TransactionStatus string

const (
	Pending   TransactionStatus = "pending"
	Completed TransactionStatus = "completed"
	Failed    TransactionStatus = "failed"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)

type Transaction struct {
	ID              string            `json:"id"`
	DebitUserID     string            `json:"debit_user_id"`
	CreditUserID    string            `json:"credit_user_id"`
	Amount          float64           `json:"amount"`
	Currency        utils.Currency    `json:"currency"`
	Status          TransactionStatus `json:"status"`
	TransactionType TransactionType   `json:"transaction_type"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Description     string            `json:"description,omitempty"`
	PaymentDetails  string            `json:"payment_details,omitempty"`
}

type TransferRequest struct {
	SenderID       string         `json:"sender_id" validate:"required"`
	ReceiverID     string         `json:"receiver_id" validate:"required,nefield=SenderID"`
	Amount         float64        `json:"amount" validate:"required,min=0"`
	Currency       utils.Currency `json:"currency" validate:"required"`
	Description    string         `json:"description"`
	PaymentDetails string         `json:"payment_details"`
}
