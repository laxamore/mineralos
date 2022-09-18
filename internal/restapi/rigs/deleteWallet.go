package rigs

import (
	"encoding/json"
	"fmt"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/internal/restapi"
	"github.com/laxamore/mineralos/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteWalletRepositoryInterface interface {
	DeleteOne(*mongo.Client, string, string, interface{}) (*mongo.DeleteResult, error)
}

type DeleteWalletController struct{}

func (a DeleteWalletController) TryDeleteWallet(c *gin.Context, client *mongo.Client, repositoryInterface DeleteWalletRepositoryInterface) {
	response := restapi.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
	}
	bodyByte, err := c.GetRawData()

	if err != nil {
		logger.Printf("newwallet get body request failed:\n%v", err)
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

			if tokenClaims["privilege"] == "admin" || tokenClaims["privilege"] == "readAndWrite" {
				objectId, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", bodyData["wallet_id"]))
				utils.CheckErr(err)

				_, err = repositoryInterface.DeleteOne(client, "mineralos", "wallets", bson.D{
					{
						Key: "_id", Value: objectId,
					},
				})

				if err != nil {
					response.Code = http.StatusNotFound
					response.Response = "wallet id not found"
					logger.Printf("error creating new wallet %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = "delete success"
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func DeleteWallet(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := databases.MongoDB{}
		cntrl := DeleteWalletController{}

		cntrl.TryDeleteWallet(c, client, repo)
	}
}
