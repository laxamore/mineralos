package ApiRigs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetRigRepositoryInterface interface {
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
}

type GetRigController struct{}

func (a GetRigController) TryGetRig(c *gin.Context, client *mongo.Client, repositoryInterface GetRigRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusNotFound,
		Response: nil,
	}
	rig_id := c.Param("rig_id")

	res := repositoryInterface.FindOne(client, "mineralos", "rigs", bson.D{
		{
			Key: "rig_id", Value: rig_id,
		},
	})

	if len(res) > 0 {
		response.Code = http.StatusOK
		response.Response = res

		c.JSON(response.Code, response.Response)
		return
	}

	c.JSON(response.Code, response.Response)
}

func GetRig(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := db.MongoDB{}
		cntrl := GetRigController{}

		cntrl.TryGetRig(c, client, repo)
	}
}
