package controllers

import (
	"crypto-wallet/database"
	"crypto-wallet/services"
	"crypto-wallet/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func SendTransaction(c *gin.Context) {

	var tx services.TransactionInput
	if err := c.BindJSON(&tx); err != nil {
		utils.Error(c, "invalid data")
		return
	}

	// Auto timestamp if missing
	if tx.Timestamp == 0 {
		tx.Timestamp = time.Now().Unix()
	}

	err := services.SendTransactionService(tx)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "transaction submitted for mining", nil)
}

func GetTransactionHistory(c *gin.Context) {

	walletID := c.Param("id")

	history, err := database.GetFullTransactionHistory(walletID)

	if err != nil {
		utils.Error(c, "failed to fetch transaction history")
		return
	}

	utils.Success(c, "transaction history fetched", history)
}
