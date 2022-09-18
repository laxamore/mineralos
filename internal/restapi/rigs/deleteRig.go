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

type DeleteRigRepositoryInterface interface {
	DeleteOne(*mongo.Client, string, string, interface{}) (*mongo.DeleteResult, error)
}

type DeleteRigController struct{}

func (a DeleteRigController) TryDeleteRig(c *gin.Context, client *mongo.Client, repositoryInterface DeleteRigRepositoryInterface) {
	response := restapi.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
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
				_, err = repositoryInterface.DeleteOne(client, "mineralos", "rigs", bson.D{
					{
						Key: "rig_id", Value: fmt.Sprintf("%s", bodyData["rig_id"]),
					},
				})

				if err != nil {
					response.Code = http.StatusNotFound
					response.Response = "rig id not found"
					logger.Printf("error creating new rig %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = "delete success"
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func DeleteRig(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := databases.MongoDB{}
		cntrl := DeleteRigController{}

		cntrl.TryDeleteRig(c, client, repo)
	}
}
