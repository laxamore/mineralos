package ApiRigs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewRigRepositoryInterface interface {
	InsertOne(string, string, interface{}) (*mongo.InsertOneResult, error)
}

type NewRigController struct{}

func (a *NewRigController) TryNewRig(c *gin.Context, repositoryInterface NewRigRepositoryInterface) {
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

				_, err = repositoryInterface.InsertOne(os.Getenv("PROJECT_NAME"), "rigs", bson.D{
					{
						Key: "rig_id", Value: newUUID.String(),
					},
					{
						Key: "rig_name", Value: fmt.Sprintf("%s", bodyData["rig_name"]),
					},
				})

				if err != nil {
					Log.Printf("error creating new rig %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = gin.H{
						"rig_id": newUUID,
					}
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func NewRig(c *gin.Context) {
	repo := db.MongoDB{}
	cntrl := NewRigController{}

	cntrl.TryNewRig(c, repo)
}
