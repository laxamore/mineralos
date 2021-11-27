package ApiRigs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"

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
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
	}
	bodyByte, err := c.GetRawData()

	if err != nil {
		Log.Printf("newwallet get body request failed:\n%v", err)
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
					Log.Printf("error creating new wallet %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = "delete success"
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func DeleteWallet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := DeleteWalletController{}

	cntrl.TryDeleteWallet(c, client, repo)
}
