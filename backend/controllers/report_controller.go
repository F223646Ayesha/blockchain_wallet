package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func GetWalletReport(c *gin.Context) {

	walletId := c.Param("id")

	data, err := services.GenerateWalletReport(walletId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "wallet report generated", data)
}
