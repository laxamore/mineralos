package rigs

//import (
//	"github.com/laxamore/mineralos/internal/db"
//	"github.com/laxamore/mineralos/internal/logger"
//	"github.com/laxamore/mineralos/internal/restapi"
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//type GetWalletsRepositoryInterface interface {
//	Find(*mongo.Client, string, string, interface{}) ([]map[string]interface{}, error)
//}
//
//type GetWalletsController struct{}
//
//func (a GetWalletsController) TryGetWallets(c *gin.Context, client *mongo.Client, repositoryInterface GetWalletsRepositoryInterface) {
//	response := restapi.Result{
//		Code:     http.StatusForbidden,
//		Response: "forbidden",
//	}
//
//	walletsList, err := repositoryInterface.Find(client, "mineralos", "wallets", bson.D{{}})
//
//	if err != nil {
//		logger.Printf("error find wallets %v", err)
//	} else {
//		response.Code = http.StatusOK
//		response.Response = gin.H{
//			"wallets": walletsList,
//		}
//	}
//
//	c.JSON(response.Code, response.Response)
//}
//
//func GetWallets(client *mongo.Client) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		repo := db.MongoDB{}
//		cntrl := GetWalletsController{}
//
//		cntrl.TryGetWallets(c, client, repo)
//	}
//}
