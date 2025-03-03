package transactions

import (
	"github.com/gin-gonic/gin"

	"concurrent_money_transfer_system/utils"
)

type TransactionController interface {
	CreateTransfer(c *gin.Context)
	GetTransaction(c *gin.Context)
	GetTransactionsByUserID(c *gin.Context)
	GetAllTransactions(c *gin.Context)
}

type transactionController struct {
	service TransactionService
}

func NewTransactionController(service TransactionService) TransactionController {
	return &transactionController{service: service}
}

func (tc *transactionController) CreateTransfer(c *gin.Context) {
	transferRequest := TransferRequest{}
	err := utils.BindAndValidateRequest(c, &transferRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	transaction, err := tc.service.CreateTransaction(c.Request.Context(), &transferRequest)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, transaction)
}

func (tc *transactionController) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResponseError(c, utils.NewErrorWithMessage(utils.ErrInvalidRequest, "Transaction ID is required"))
		return
	}
	transaction, err := tc.service.GetTransaction(c.Request.Context(), id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, transaction)
}

func (tc *transactionController) GetTransactionsByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		utils.ResponseError(c, utils.NewErrorWithMessage(utils.ErrInvalidRequest, "User ID is required"))
		return
	}

	transactions, err := tc.service.GetTransactionsByUserID(c.Request.Context(), userID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, transactions)
}

func (tc *transactionController) GetAllTransactions(c *gin.Context) {
	transactions, err := tc.service.GetAllTransactions(c.Request.Context())
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, transactions)
}