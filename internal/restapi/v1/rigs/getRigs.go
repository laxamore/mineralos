package rigs

import (
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
)

func GetRigs(c *gin.Context) {
	ctrl := RigController{
		DB: db.DB,
	}

	ctrl.GetRigs(c)
}

func (ctrl RigController) GetRigs(c *gin.Context) {
	rigs := []models.Rig{}
	err := ctrl.DB.Find(&rigs).Error

	if err != nil {
		logger.Errorf("Error while getting rigs: %v", err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if len(rigs) > 0 {
		c.JSON(http.StatusOK, rigs)
		return
	}

	c.JSON(http.StatusNoContent, "No record found")
}
