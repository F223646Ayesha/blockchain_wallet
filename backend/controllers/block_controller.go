package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func GetBlockchain(c *gin.Context) {
	blocks, err := services.GetBlockchainService()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}
	utils.Success(c, "blockchain fetched successfully", blocks)
}

func MinePendingTransactions(c *gin.Context) {
	miner := c.Query("miner")

	if miner == "" {
		utils.Error(c, "miner wallet is required")
		return
	}

	block, err := services.MinePendingTransactionsService(miner)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "block mined successfully", block)
}
