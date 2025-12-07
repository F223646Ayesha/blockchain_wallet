package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

// ===============================
// SEND OTP (SMTP Email OTP)
// ===============================
func SendOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.BindJSON(&body) != nil {
		utils.Error(c, "invalid request")
		return
	}

	err := services.SendOTPService(body.Email)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "OTP sent to email!", nil)
}

// ===============================
// LOGIN USING EMAIL + OTP
// ===============================
