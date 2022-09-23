package rigs

import (
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"net/http"
)

func GetRig(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.GetRig(c)
}

func (ctrl RigController) GetRig(c *gin.Context) {
	rig_id := c.Param("rig_id")

	rig := models.Rig{}
	err := ctrl.DB.First(&rig, "rig_id = ?", rig_id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, "Rig not found")
		return
	}

	c.JSON(http.StatusOK, rig)
}
