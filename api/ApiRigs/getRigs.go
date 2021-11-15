package ApiRigs

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
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

	rigsList, err := repositoryInterface.Find(client, os.Getenv("PROJECT_NAME"), "rigs", bson.D{{}})

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

func GetRigs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := GetRigsController{}

	cntrl.TryGetRigs(c, client, repo)
}
