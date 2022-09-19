package rig

import (
	"encoding/json"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/databases/models/rig"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/internal/restapi"
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
		DB: databases.DB,
	}
	ctrl.NewRig(c)
}

func (ctrl RigController) NewRig(c *gin.Context) {
	response := restapi.Result{
		Code: http.StatusForbidden,
		Response: map[string]interface{}{
			"rig_id": nil,
		},
	}

	bodyByte, err := c.GetRawData()

	if err != nil {
		logger.Printf("newRig get body request failed:\n%v", err)
		c.JSON(response.Code, response.Response)
		return
	}

	var bodyData NewRigRequest
	err = json.Unmarshal(bodyByte, &bodyData)

	if err != nil {
		logger.Printf("newRig unmarshal body request failed:\n%v", err)
		c.JSON(response.Code, response.Response)
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

	err = ctrl.DB.Create(&rig.Rig{
		RigID:   newUUID.String(),
		RigName: bodyData.RigName,
	}).Error

	if err != nil {
		logger.Printf("error creating new rig %v", err)
	} else {
		response.Code = http.StatusOK
		response.Response = gin.H{
			"rig_name": bodyData.RigName,
			"rig_id":   newUUID,
		}
	}

	c.JSON(response.Code, response.Response)
}
