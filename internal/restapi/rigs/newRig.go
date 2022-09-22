package rigs

import (
	"encoding/json"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NewRigRequest struct {
	RigName string `json:"rig_name"`
}

func NewRig(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.NewRig(c)
}

func (ctrl RigController) NewRig(c *gin.Context) {
	bodyByte, err := c.GetRawData()

	if err != nil {
		logger.Printf("newRig get body request failed:\n%v", err)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var bodyData NewRigRequest
	err = json.Unmarshal(bodyByte, &bodyData)

	if err != nil {
		logger.Printf("newRig unmarshal body request failed:\n%v", err)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	reflectNewRigRequest := reflect.ValueOf(&bodyData).Elem()
	typeOfNewRigRequest := reflectNewRigRequest.Type()
	for i := 0; i < reflectNewRigRequest.NumField(); i++ {
		if reflectNewRigRequest.Field(i).Interface() == "" {
			logger.Printf("newRig request failed: %s is empty", typeOfNewRigRequest.Field(i).Name)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": typeOfNewRigRequest.Field(i).Tag.Get("json") + " is undefined",
			})
			return
		}
	}

	newUUID := uuid.New()

	err = ctrl.DB.Create(&models.Rig{
		RigID:   newUUID.String(),
		RigName: bodyData.RigName,
	}).Error

	if err != nil {
		logger.Printf("error creating new rigs %v", err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"rig_name": bodyData.RigName,
			"rig_id":   newUUID,
		})
	}
}
