package ApiRigs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteRigRepositoryInterface interface {
	DeleteOne(*mongo.Client, string, string, interface{}) (*mongo.DeleteResult, error)
}

type DeleteRigController struct{}

func (a DeleteRigController) TryDeleteRig(c *gin.Context, client *mongo.Client, repositoryInterface DeleteRigRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
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
				_, err = repositoryInterface.DeleteOne(client, os.Getenv("PROJECT_NAME"), "rigs", bson.D{
					{
						Key: "rig_id", Value: fmt.Sprintf("%s", bodyData["rig_id"]),
					},
				})

				if err != nil {
					response.Code = http.StatusNotFound
					response.Response = "rig id not found"
					Log.Printf("error creating new rig %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = "delete success"
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func DeleteRig(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := DeleteRigController{}

	cntrl.TryDeleteRig(c, client, repo)
}
