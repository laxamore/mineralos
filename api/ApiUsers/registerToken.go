package ApiUsers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegisterTokenRepositoryInterface interface {
	InsertOne(*mongo.Client, string, string, interface{}) (*mongo.InsertOneResult, error)
	IndexesReplaceOne(*mongo.Client, string, string, mongo.IndexModel) (string, error)
}

type RegisterTokenController struct{}

func (a RegisterTokenController) TryRegisterToken(c *gin.Context, client *mongo.Client, repositoryInterface RegisterTokenRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusInternalServerError,
		Response: "Internal Server Error",
	}
	errMsg := ""

	isAdmin, _ := c.Get("admin")

	if !isAdmin.(bool) {
		response.Code = http.StatusUnauthorized
		response.Response = "Unauthorized"
	} else {
		bodyRaw, err := c.GetRawData()

		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Response = "Internal Server Error"
			errMsg = fmt.Sprintf("RegisterToken Get Body Request Failed:\n%v", err)
		}

		var bodyData map[string]interface{}
		json.Unmarshal(bodyRaw, &bodyData)

		if bodyData["privilege"] == nil || (bodyData["privilege"] != "readOnly" && bodyData["privilege"] != "readAndWrite") {
			response.Code = http.StatusBadRequest
			response.Response = "Bad Request"
		} else {
			var expireAfterSeconds int32 = 43200
			createIndexRes, err := repositoryInterface.IndexesReplaceOne(client, "mineralos", "registerToken", mongo.IndexModel{Keys: bson.D{{Key: "createdAt", Value: 1}}, Options: &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds}})
			if err != nil {
				response.Code = http.StatusInternalServerError
				response.Response = "Internal Server Error"
				errMsg = fmt.Sprintf("RegisterToken CreateIndex Failed:\n%v", err)
			} else {
				Log.Printf("RegisterToken CreateIndex Response: %v", createIndexRes)

				registerToken := func() string {
					b := make([]byte, 8)
					rand.Read(b)
					return fmt.Sprintf("%x", b)
				}()

				res, err := repositoryInterface.InsertOne(client, "mineralos", "registerToken", bson.D{
					{
						Key: "createdAt", Value: time.Now(),
					},
					{
						Key: "token", Value: registerToken,
					},
					{
						Key: "privilege", Value: bodyData["privilege"],
					},
				})

				if err != nil {
					response.Code = http.StatusInternalServerError
					response.Response = "Internal Server Error"
					errMsg = fmt.Sprintf("RegisterToken InsertOne Failed:\n%v", err)
				}
				Log.Printf("RegisterToken InsertOne Respone: %v", res)

				response.Code = http.StatusCreated
				response.Response = gin.H{
					"token": registerToken,
				}
			}
		}
	}

	c.JSON(response.Code, response.Response)
	if errMsg != "" {
		Log.Panic(errMsg)
	}
}

func RegisterToken(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := db.MongoDB{}
		cntrl := RegisterTokenController{}

		cntrl.TryRegisterToken(c, client, repo)
	}
}
