package controllers

import (
	"crypto-wallet/services"
	"crypto-wallet/utils"

	"github.com/gin-gonic/gin"
)

func CreateWallet(c *gin.Context) {

	resp, err := services.CreateWalletService()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "wallet created successfully", resp)
}

func GetWallet(c *gin.Context) {

	id := c.Param("id")
	data, err := services.GetWalletService(id)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "wallet fetched", data)
}

// Get wallet + beneficiaries
func GetWalletProfile(c *gin.Context) {
	id := c.Param("id")

	data, err := services.GetWalletService(id)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "wallet profile fetched", data)
}
func AddBeneficiary(c *gin.Context) {
	var body struct {
		WalletID      string `json:"wallet_id"`
		BeneficiaryID string `json:"beneficiary_id"`
	}

	if err := c.BindJSON(&body); err != nil {
		utils.Error(c, "invalid data")
		return
	}

	err := services.AddBeneficiaryService(body.WalletID, body.BeneficiaryID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "beneficiary added", nil)
}
func RemoveBeneficiary(c *gin.Context) {
	var body struct {
		WalletID      string `json:"wallet_id"`
		BeneficiaryID string `json:"beneficiary_id"`
	}

	if err := c.BindJSON(&body); err != nil {
		utils.Error(c, "invalid data")
		return
	}

	err := services.RemoveBeneficiaryService(body.WalletID, body.BeneficiaryID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "beneficiary removed", nil)
}
func UpdateWalletProfile(c *gin.Context) {

	var body struct {
		WalletID string `json:"wallet_id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		CNIC     string `json:"cnic"`
	}

	if err := c.BindJSON(&body); err != nil {
		utils.Error(c, "invalid request body")
		return
	}

	err := services.UpdateWalletProfileService(body.WalletID, body.Name, body.Email, body.CNIC)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, "profile updated successfully", nil)
}
