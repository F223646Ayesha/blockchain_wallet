package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func ValidateBlockchain(c *gin.Context) {

	blocks, err := services.GetBlockchainService()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	ok, msg := services.ValidateBlockchainService(blocks)
	if !ok {
		utils.Error(c, msg)
		return
	}

	utils.Success(c, msg, nil)
}
