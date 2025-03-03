package server

import (
	"concurrent_money_transfer_system/internals/transactions"
	"concurrent_money_transfer_system/internals/users"
	"concurrent_money_transfer_system/internals/wallet"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	setupUserRoutes(router)
	setupWalletRoutes(router)
	setupTransactionRoutes(router)

	return router
}

func setupUserRoutes(router *gin.Engine) {
	userRepo := users.NewUserRepo()
	walletRepo := wallet.NewWalletRepo()
	walletService := wallet.NewWalletService(walletRepo)
	userService := users.NewUserService(userRepo, walletService)
	userController := users.NewUserController(userService)

	userRouter := router.Group("api/user")
	{
		userRouter.POST("/signup", userController.CreateUser)
		userRouter.GET("/:id", userController.GetUser)
		userRouter.GET("/", userController.GetAllUsers)
		userRouter.PUT("/:id", userController.UpdateUser)
		userRouter.DELETE("/:id", userController.DeleteUser)
	}
}

func setupWalletRoutes(router *gin.Engine) {
	walletRepo := wallet.NewWalletRepo()
	walletService := wallet.NewWalletService(walletRepo)
	walletController := wallet.NewWalletController(walletService)

	walletRouter := router.Group("/wallets")
	{
		walletRouter.GET("", walletController.GetWallet)
		walletRouter.PUT("/disable", walletController.DisableWallet)
	}
}

func setupTransactionRoutes(router *gin.Engine) {
	transactionRepo := transactions.NewTransactionRepo()
	walletRepo := wallet.NewWalletRepo()
	walletService := wallet.NewWalletService(walletRepo)
	transactionService := transactions.NewTransactionService(transactionRepo, walletService)
	transactionController := transactions.NewTransactionController(transactionService)
	transactionRouter := router.Group("api/transaction")
	{
		transactionRouter.POST("/transfer", transactionController.CreateTransfer)
		transactionRouter.GET("/:id", transactionController.GetTransaction)
		transactionRouter.GET("/user/:user_id", transactionController.GetTransactionsByUserID)
		transactionRouter.GET("/", transactionController.GetAllTransactions)
	}
}
