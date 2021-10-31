package ApiRigs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mineralos/api"
	"mineralos/db"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteRig(c *gin.Context) {
	result := api.Result{
		Code: 400,
		Response: map[string]interface{}{
			"msg": "delete failed",
		},
	}

	bodyByte, err := c.GetRawData()

	var bodyData map[string]interface{}
	json.Unmarshal(bodyByte, &bodyData)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.MongoClient(ctx)
	defer cancel()

	if err != nil {
		log.Fatalf("DB Connection Error:\n%v", err)
	}

	collection := client.Database(fmt.Sprintf("%s", os.Getenv("PROJECT_NAME"))).Collection("rigs")

	_, err = collection.DeleteOne(ctx, bson.D{{
		"rig_id", bodyData["rig_id"],
	}})

	if err != nil {
		log.Fatalf("Delete RIG Error:\n%v", err)
	}

	result.Code = 200
	result.Response = map[string]interface{}{
		"msg": "delete success",
	}

	c.JSON(result.Code, result.Response)
}
