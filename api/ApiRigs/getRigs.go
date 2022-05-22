package ApiRigs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetRigsRepositoryInterface interface {
	Find(*mongo.Client, string, string, interface{}) ([]map[string]interface{}, error)
}

type GetRigsController struct{}

func (a GetRigsController) TryGetRigs(c *gin.Context, client *mongo.Client, repositoryInterface GetRigsRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
	}

	rigsList, err := repositoryInterface.Find(client, "mineralos", "rigs", bson.D{{}})

	if err != nil {
		Log.Printf("error find rigs %v", err)
	} else {
		response.Code = http.StatusOK
		response.Response = gin.H{
			"rigs": rigsList,
		}
	}

	c.JSON(response.Code, response.Response)
}

func GetRigs(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := db.MongoDB{}
		cntrl := GetRigsController{}

		cntrl.TryGetRigs(c, client, repo)
	}
}
