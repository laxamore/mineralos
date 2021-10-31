package ApiRigs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func NewRig(c *gin.Context) {
	result := api.Result{
		Code: 400,
		Response: map[string]interface{}{
			"rig_id": nil,
		},
	}

	newUUID := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.MongoClient(ctx)
	defer cancel()

	if err != nil {
		log.Panicf("DB Connection Error:\n%v", err)
	}

	collection := client.Database(fmt.Sprintf("%s", os.Getenv("PROJECT_NAME"))).Collection("rigs")

	_, err = collection.InsertOne(ctx, bson.D{{
		"rig_id", fmt.Sprintf("%s", newUUID),
	}})

	if err != nil {
		log.Panicf("Create New RIG Error:\n%v", err)
	}

	result.Code = 200
	result.Response = map[string]interface{}{
		"rig_id": fmt.Sprintf("%s", newUUID),
	}

	c.JSON(result.Code, result.Response)
}
