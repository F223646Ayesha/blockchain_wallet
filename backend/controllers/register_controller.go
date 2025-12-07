package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		CNIC  string `json:"cnic"`
	}

	if err := c.BindJSON(&body); err != nil {
		utils.Error(c, "invalid registration data")
		return
	}

	// Actual service
	resp, err := services.RegisterUserService(body.Name, body.Email, body.CNIC)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "registration complete", resp)
}
