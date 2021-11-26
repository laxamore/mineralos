package ApiRigs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"
	"github.com/pebbe/zmq4"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewRigRepositoryInterface interface {
	InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
}

type NewRigController struct{}

func (a *NewRigController) TryNewRig(c *gin.Context, client *mongo.Client, repositoryInterface NewRigRepositoryInterface) {
	response := api.Result{
		Code: http.StatusForbidden,
		Response: map[string]interface{}{
			"rig_id": nil,
		},
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

				newUUID := uuid.New()
				pubKey, key, err := zmq4.NewCurveKeypair()
				if err != nil {
					Log.Panicf("NewCurveKeypair: %v", err)
				}

				insertResultID, err := repositoryInterface.InsertOne(client, os.Getenv("PROJECT_NAME"), "rigs", bson.D{
					{
						Key: "rig_id", Value: newUUID.String(),
					},
					{
						Key: "rig_name", Value: fmt.Sprintf("%s", bodyData["rig_name"]),
					},
					{
						Key: "key", Value: key,
					},
					{
						Key: "pubkey", Value: pubKey,
					},
				})

				if err != nil {
					Log.Printf("error creating new rig %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = gin.H{
						"_id":      insertResultID.InsertedID,
						"rig_name": fmt.Sprintf("%s", bodyData["rig_name"]),
						"rig_id":   newUUID,
					}
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func NewRig(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := NewRigController{}

	cntrl.TryNewRig(c, client, repo)
}
