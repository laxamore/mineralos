package ApiRigs

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
)

type GetRigsRepositoryInterface interface {
	Find(string, string, interface{}) ([]map[string]interface{}, error)
}

type GetRigsController struct{}

func (a GetRigsController) TryGetRigs(c *gin.Context, repositoryInterface GetRigsRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusForbidden,
		Response: "forbidden",
	}

	rigsList, err := repositoryInterface.Find(os.Getenv("PROJECT_NAME"), "rigs", bson.D{{}})

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
	repo := db.MongoDB{}
	cntrl := GetRigsController{}

	cntrl.TryGetRigs(c, repo)
}
