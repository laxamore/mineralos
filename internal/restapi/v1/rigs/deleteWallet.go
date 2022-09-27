package rigs

import (
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"net/http"
)

type DeleteWalletRequest struct {
	WalletId uint `json:"wallet_id" binding:"required"`
}

func DeleteWallet(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.DeleteWallet(c)
}

func (ctrl RigController) DeleteWallet(c *gin.Context) {
	var request DeleteWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	err := ctrl.DB.Delete(&models.Wallet{}, "id = ?", request.WalletId).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
