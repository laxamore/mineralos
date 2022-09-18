package middleware

import (
	"context"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BeforeRegisterRepositoryInterface interface {
	Find(*mongo.Client, string, string, interface{}) ([]map[string]interface{}, error)
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
}

type BeforeRegisterController struct{}

func (a BeforeRegisterController) TryBeforeRegister(c *gin.Context, client *mongo.Client, repositoryInterface BeforeRegisterRepositoryInterface) {
	c.Set("token", nil)
	regToken := c.GetHeader("regToken")

	results, err := repositoryInterface.Find(client, "mineralos", "users", bson.D{{}})
	// logger.Printf("%v", len(results))

	if err != nil {
		logger.Panicf("BeforeRegister List All Users Failed:\n%v", err)
	}

	if len(results) > 0 {
		result := repositoryInterface.FindOne(client, "mineralos", "registerToken", bson.D{{Key: "token", Value: regToken}})
		// logger.Printf("%v", result)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := databases.MongoClient(ctx)
	utils.CheckErr(err)

	repo := databases.MongoDB{}
	cntrl := BeforeRegisterController{}

	cntrl.TryBeforeRegister(c, client, repo)
}
