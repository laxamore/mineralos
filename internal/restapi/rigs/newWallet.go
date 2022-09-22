package rigs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
)

func NewWallet(c *gin.Context) {
	ctrl := &RigController{
		DB: db.DB,
	}
	ctrl.NewWallet(c)
}

func (ctrl RigController) NewWallet(c *gin.Context) {
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
			//	insertResultID, err := repositoryInterface.InsertOne(client, "mineralos", "wallets", bson.D{
			//		{
			//			Key: "wallet_name", Value: fmt.Sprintf("%s", bodyData["wallet_name"]),
			//		},
			//		{
			//			Key: "wallet_address", Value: fmt.Sprintf("%s", bodyData["wallet_address"]),
			//		},
			//		{
			//			Key: "coin", Value: fmt.Sprintf("%s", bodyData["coin"]),
			//		},
			//	})
			//
			//	if err != nil {
			//		logger.Printf("error creating new rigs %v", err)
			//	} else {
			//		response.Code = http.StatusOK
			//		response.Response = gin.H{
			//			"_id":            insertResultID.InsertedID,
			//			"wallet_name":    fmt.Sprintf("%s", bodyData["wallet_name"]),
			//			"wallet_address": fmt.Sprintf("%s", bodyData["wallet_address"]),
			//			"coin":           fmt.Sprintf("%s", bodyData["coin"]),
			//		}
			//	}
			//}
		}
	}

	c.JSON(http.StatusBadRequest, "Bad Request")
}
