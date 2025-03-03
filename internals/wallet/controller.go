package wallet

import (
	"github.com/gin-gonic/gin"

	"concurrent_money_transfer_system/utils"
)

type WalletController interface {
	DisableWallet(c *gin.Context)
	GetWallet(c *gin.Context)
}

type walletController struct {
	service WalletService
}

func NewWalletController(service WalletService) WalletController {
	return &walletController{service: service}
}

func (wc *walletController) DisableWallet(c *gin.Context) {
	userID := c.Query("user_id")
	err := wc.service.DisableWallet(c.Request.Context(), userID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, gin.H{"message": "Wallet disabled successfully"})
}

func (wc *walletController) GetWallet(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		utils.ResponseError(c, utils.NewErrorWithMessage(utils.ErrInvalidRequest, "User ID is required"))
		return
	}
	wallet, err := wc.service.GetWallet(c.Request.Context(), userID)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, wallet)
}
