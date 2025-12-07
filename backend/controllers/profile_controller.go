package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	walletId := c.Param("walletId")

	data, err := services.GetUserProfileService(walletId)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "profile fetched", data)
}

func UpdateUserProfile(c *gin.Context) {

	var body struct {
		WalletID string `json:"wallet_id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		CNIC     string `json:"cnic"`
	}

	if err := c.BindJSON(&body); err != nil {
		utils.Error(c, "invalid request")
		return
	}

	err := services.UpdateUserProfileService(body.WalletID, body.Name, body.Email, body.CNIC)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "profile updated", nil)
}
