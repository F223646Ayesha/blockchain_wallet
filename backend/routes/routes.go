package routes

import (
	"crypto-wallet/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api")

	// ============================
	// AUTH
	// ============================
	api.POST("/register", controllers.RegisterUser)
	api.POST("/auth/send-otp", controllers.SendOTP)
	api.POST("/login", controllers.Login)

	// ============================
	// WALLET
	// ============================
	api.POST("/wallet/create", controllers.CreateWallet)
	api.GET("/wallet/:id", controllers.GetWallet)
	api.GET("/wallet/profile/:id", controllers.GetWalletProfile)
	api.POST("/wallet/profile/update", controllers.UpdateWalletProfile)
	api.POST("/wallet/beneficiary/add", controllers.AddBeneficiary)
	api.POST("/wallet/beneficiary/remove", controllers.RemoveBeneficiary)
	api.GET("/user/profile/:walletId", controllers.GetUserProfile)
	api.POST("/user/profile/update", controllers.UpdateUserProfile)

	// ============================
	// TRANSACTION
	// ============================
	api.POST("/transaction/send", controllers.SendTransaction)
	api.GET("/transaction/history/:id", controllers.GetTransactionHistory)

	// ============================
	// BLOCKCHAIN + MINING
	// ============================
	api.GET("/blockchain", controllers.GetBlockchain)
	api.POST("/mine", controllers.MinePendingTransactions)

	// ============================
	// LOGS
	// ============================
	api.GET("/logs", controllers.GetLogs)

	// ============================
	// ZAKAT
	// ============================
	api.POST("/zakat/run", controllers.RunZakat)
	api.GET("/analytics/system", controllers.GetSystemAnalytics)
	api.GET("/reports/wallet/:id", controllers.GetWalletReport)
	api.GET("/blockchain/validate", controllers.ValidateBlockchain)

}
