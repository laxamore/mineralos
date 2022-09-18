package rigs

import (
	"encoding/json"
	"fmt"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/internal/restapi"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewWalletRepositoryInterface interface {
	InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
}

type NewWalletController struct{}

func (a NewWalletController) TryNewWallet(c *gin.Context, client *mongo.Client, repositoryInterface NewWalletRepositoryInterface) {
	response := restapi.Result{
		Code:     http.StatusForbidden,
		Response: "Forbidden",
	}

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

			if tokenClaims["privilege"] == "admin" || tokenClaims["privilege"] == "readAndWrite" {
				insertResultID, err := repositoryInterface.InsertOne(client, "mineralos", "wallets", bson.D{
					{
						Key: "wallet_name", Value: fmt.Sprintf("%s", bodyData["wallet_name"]),
					},
					{
						Key: "wallet_address", Value: fmt.Sprintf("%s", bodyData["wallet_address"]),
					},
					{
						Key: "coin", Value: fmt.Sprintf("%s", bodyData["coin"]),
					},
				})

				if err != nil {
					logger.Printf("error creating new rig %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = gin.H{
						"_id":            insertResultID.InsertedID,
						"wallet_name":    fmt.Sprintf("%s", bodyData["wallet_name"]),
						"wallet_address": fmt.Sprintf("%s", bodyData["wallet_address"]),
						"coin":           fmt.Sprintf("%s", bodyData["coin"]),
					}
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func NewWallet(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := databases.MongoDB{}
		cntrl := NewWalletController{}

		cntrl.TryNewWallet(c, client, repo)
	}
}
