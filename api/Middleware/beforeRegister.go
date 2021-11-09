package Middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
)

type BeforeRegisterRepositoryInterface interface {
	Find(string, string, interface{}) ([]map[string]interface{}, error)
	FindOne(string, string, interface{}) map[string]interface{}
}

type BeforeRegisterController struct{}

func (a BeforeRegisterController) TryBeforeRegister(c *gin.Context, repositoryInterface BeforeRegisterRepositoryInterface) {
	c.Set("token", nil)
	regToken := c.GetHeader("regToken")

	results, err := repositoryInterface.Find(os.Getenv("PROJECT_NAME"), "users", bson.D{{}})
	// Log.Printf("%v", len(results))

	if err != nil {
		Log.Panicf("BeforeRegister List All Users Failed:\n%v", err)
	}

	if len(results) > 0 {
		result := repositoryInterface.FindOne(os.Getenv("PROJECT_NAME"), "registerToken", bson.D{{Key: "token", Value: regToken}})
		// Log.Printf("%v", result)

		if len(result) > 0 {
			c.Set("token", result)
			c.Writer.WriteHeader(http.StatusOK)
			return
		}
	} else {
		c.Set("registerAdmin", true)
		c.Writer.WriteHeader(http.StatusOK)
		return
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
