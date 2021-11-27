package ApiRigs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewWalletRepositoryInterface interface {
	InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
}

type NewWalletController struct{}

func (a NewWalletController) TryNewWallet(c *gin.Context, client *mongo.Client, repositoryInterface NewWalletRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "Forbidden",
	}

	bodyByte, err := c.GetRawData()

	if err != nil {
		Log.Printf("newrig get body request failed:\n%v", err)
	} else {
		var bodyData map[string]interface{}
		json.Unmarshal(bodyByte, &bodyData)

		res, _ := c.Get("tokenClaims")
		tokenClaimsByte, err := json.Marshal(res)

		if err != nil {
			Log.Printf("error marshal tokenClaims %v", err)
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
					Log.Printf("error creating new rig %v", err)
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

func NewWallet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := NewWalletController{}

	cntrl.TryNewWallet(c, client, repo)
}
