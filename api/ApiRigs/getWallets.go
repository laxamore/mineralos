package ApiRigs

import (
	"context"
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

type GetWalletsRepositoryInterface interface {
	Find(*mongo.Client, string, string, interface{}) ([]map[string]interface{}, error)
}

type GetWalletsController struct{}

func (a GetWalletsController) TryGetWallets(c *gin.Context, client *mongo.Client, repositoryInterface GetWalletsRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
	}

	walletsList, err := repositoryInterface.Find(client, "mineralos", "wallets", bson.D{{}})

	if err != nil {
		Log.Printf("error find wallets %v", err)
	} else {
		response.Code = http.StatusOK
		response.Response = gin.H{
			"wallets": walletsList,
		}
	}

	c.JSON(response.Code, response.Response)
}

func GetWallets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	repo := db.MongoDB{}
	cntrl := GetWalletsController{}

	cntrl.TryGetWallets(c, client, repo)
}
