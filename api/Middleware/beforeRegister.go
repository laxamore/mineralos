package Middleware

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/db"
	Log "github.com/laxamore/mineralos/log"
	"go.mongodb.org/mongo-driver/bson"
)

type BeforeRegisterRepositoryInterface interface {
	Find(string, string, interface{}) ([]map[string]interface{}, error)
	FindOne(string, string, interface{}) map[string]interface{}
}

type BeforeRegisterController struct{}

func (a BeforeRegisterController) TryBeforeRegister(c *gin.Context, repositoryInterface BeforeRegisterRepositoryInterface) {
	c.Set("token", nil)
	bodyRaw, err := c.GetRawData()

	if err != nil {
		Log.Panicf("BeforeRegister Get Body Request Failed:\n%v", err)
	}

	var bodyData map[string]interface{}
	json.Unmarshal(bodyRaw, &bodyData)
	c.Set("bodyData", bodyData)

	if bodyData["username"] != nil && bodyData["email"] != nil && bodyData["password"] != nil {
		results, err := repositoryInterface.Find(os.Getenv("PROJECT_NAME"), "users", bson.D{{}})

		if err != nil {
			Log.Panicf("BeforeRegister List All Users Failed:\n%v", err)
		}

		if len(results) > 0 {
			if bodyData["token"] != nil {
				result := repositoryInterface.FindOne(os.Getenv("PROJECT_NAME"), "registerToken", bson.D{{Key: "token", Value: bodyData["token"]}})

				Log.Printf("%v", result)

				if len(result) > 0 {
					c.Set("token", result)
					c.Writer.WriteHeader(http.StatusOK)
					return
				}
			}
		} else {
			c.Set("registerAdmin", true)
		}
	}

	c.Abort()
	c.Writer.WriteHeader(http.StatusUnauthorized)
	c.Writer.Write([]byte("Unauthorized"))
}

func BeforeRegister(c *gin.Context) {
	repo := db.MongoDB{}
	cntrl := BeforeRegisterController{}

	cntrl.TryBeforeRegister(c, repo)
}
