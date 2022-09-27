package rigs

import (
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"net/http"
)

type NewWalletRequest struct {
	WalletCoin    string `json:"wallet_coin" binding:"required"`
	WalletName    string `json:"wallet_name" binding:"required"`
	WalletAddress string `json:"wallet_address" binding:"required"`
}

func NewWallet(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.NewWallet(c)
}

func (ctrl RigController) NewWallet(c *gin.Context) {
	var request NewWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	newWallet := models.Wallet{}
	err := ctrl.DB.Create(&models.Wallet{
		WalletCoin:    request.WalletCoin,
		WalletName:    request.WalletName,
		WalletAddress: request.WalletAddress,
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"wallet": newWallet,
	})
}
