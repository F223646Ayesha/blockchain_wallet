package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func GetSystemAnalytics(c *gin.Context) {

	data, err := services.GetSystemAnalyticsService()
	if err != nil {
		utils.Error(c, "Failed to load analytics")
		return
	}

	utils.Success(c, "analytics loaded", data)
}
