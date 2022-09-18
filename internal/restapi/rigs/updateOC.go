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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateOCRepository interface {
	// InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(*mongo.Client, string, string, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type UpdateOCController struct{}
type UpdateOCProps struct {
	RIG_ID string `json:"rig_id"`
	Vendor string `json:"vendor"`
	ID     uint8  `json:"id"`
	FS     uint8  `json:"fs"`
	CC     int    `json:"cc"`
	CV     int    `json:"cv"`
	MC     int    `json:"mc"`
	MV     int    `json:"mv"`
	PL     uint8  `json:"pl"`
}

func (a *UpdateOCController) TryUpdateOC(c *gin.Context, client *mongo.Client, repositoryInterface UpdateOCRepository) {
	response := restapi.Result{
		Code: http.StatusForbidden,
		Response: map[string]interface{}{
			"status": nil,
		},
	}

	bodyByte, err := c.GetRawData()

	if err != nil {
		logger.Printf("newrig get body request failed:\n%v", err)
	} else {
		var bodyData UpdateOCProps
		json.Unmarshal(bodyByte, &bodyData)

		res, _ := c.Get("tokenClaims")
		tokenClaimsByte, err := json.Marshal(res)

		if err != nil {
			logger.Printf("error marshal tokenClaims %v", err)
		} else {
			var tokenClaims map[string]interface{}
			json.Unmarshal(tokenClaimsByte, &tokenClaims)

			if tokenClaims["privilege"] == "admin" || tokenClaims["privilege"] == "readAndWrite" {
				update := bson.M{
					"$set": bson.M{
						fmt.Sprintf("oc.%s.%d", bodyData.Vendor, bodyData.ID): bson.M{
							"fs": int(bodyData.FS),
							"cc": int(bodyData.CC),
							"cv": int(bodyData.CV),
							"mc": int(bodyData.MC),
							"mv": int(bodyData.MV),
							"pl": int(bodyData.PL),
						},
					},
				}

				_, err := repositoryInterface.UpdateOne(client, "mineralos", "rigs", bson.M{
					"rig_id": bodyData.RIG_ID,
				}, update)

				if err != nil {
					logger.Printf("error updating overclock %v", err)
				} else {
					response.Code = http.StatusOK
					response.Response = gin.H{
						"status": "update overclock OK",
						"oc": gin.H{
							"fs": bodyData.FS,
							"cc": bodyData.CC,
							"cv": bodyData.CV,
							"mc": bodyData.MC,
							"mv": bodyData.MV,
							"pl": bodyData.PL,
						},
					}
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
}

func UpdateOC(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := databases.MongoDB{}
		cntrl := UpdateOCController{}

		cntrl.TryUpdateOC(c, client, repo)
	}
}
