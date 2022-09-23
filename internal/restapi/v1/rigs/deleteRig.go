package rigs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
)

type DeleteRigRequest struct {
	RigID string `json:"rig_id"`
}

func DeleteRig(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.DeleteRig(c)
}

func (ctrl RigController) DeleteRig(c *gin.Context) {
	bodyByte, err := c.GetRawData()
	deleteRigRequest := DeleteRigRequest{}
	err = json.Unmarshal(bodyByte, &deleteRigRequest)

	if err != nil {
		logger.Errorf("newrig get body request failed:\n%v", err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rig := models.Rig{}
	err = ctrl.DB.Delete(&rig, "rig_id = ?", deleteRigRequest.RigID).Error
	if err != nil {
		logger.Errorf("error deleting rig %v", err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, "OK")
}
