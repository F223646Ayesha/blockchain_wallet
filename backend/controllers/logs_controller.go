package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func GetLogs(c *gin.Context) {

	logs, err := services.GetLogsService()
	if err != nil {
		utils.Error(c, "failed to fetch logs")
		return
	}

	utils.Success(c, "logs fetched successfully", logs)
}
