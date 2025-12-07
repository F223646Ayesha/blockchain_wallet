package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func Login(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if c.BindJSON(&body) != nil {
		utils.Error(c, "invalid request body")
		return
	}

	token, user, wallet, err := services.LoginService(body.Email, body.OTP)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "Login successful", gin.H{
		"token":  token,
		"user":   user,
		"wallet": wallet,
	})
}
