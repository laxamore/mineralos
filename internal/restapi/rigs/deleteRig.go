package rigs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
)

func DeleteRig(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.DeleteRig(c)
}

func (ctrl RigController) DeleteRig(c *gin.Context) {
	bodyByte, err := c.GetRawData()

	if err != nil {
		logger.Printf("newrig get body request failed:\n%v", err)
	} else {
		var bodyData map[string]interface{}
		json.Unmarshal(bodyByte, &bodyData)

		res, _ := c.Get("tokenClaims")
		tokenClaimsByte, err := json.Marshal(res)

		if err != nil {
			logger.Printf("error marshal tokenClaims %v", err)
		} else {
			var tokenClaims map[string]interface{}
			json.Unmarshal(tokenClaimsByte, &tokenClaims)

			//if tokenClaims["privilege"] == "admin" || tokenClaims["privilege"] == "readAndWrite" {
			//	_, err = repositoryInterface.DeleteOne(client, "mineralos", "rigs", bson.D{
			//		{
			//			Key: "rig_id", Value: fmt.Sprintf("%s", bodyData["rig_id"]),
			//		},
			//	})
			//
			//	if err != nil {
			//		response.Code = http.StatusNotFound
			//		response.Response = "rigs id not found"
			//		logger.Printf("error creating new rigs %v", err)
			//	} else {
			//		response.Code = http.StatusOK
			//		response.Response = "delete success"
			//	}
			//}
		}
	}

	c.JSON(http.StatusBadRequest, "Bad Request")
}
