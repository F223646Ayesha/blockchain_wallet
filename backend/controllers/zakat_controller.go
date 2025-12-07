package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func RunZakat(c *gin.Context) {

	err := services.RunZakatService()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "Zakat processed successfully", nil)
}
