package ApiUsers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context) {
	_, registerAdmin := c.Get("registerAdmin")
	tokenInfo := c.GetStringMap("token")
	bodyData := c.GetStringMap("bodyData")

	privilege := "readOnly"

	if tokenInfo != nil {
		privilege = fmt.Sprintf("%s", tokenInfo["privilege"])
	}

	if registerAdmin {
		privilege = "admin"
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s", bodyData["password"])))
	sha256_hash := hex.EncodeToString(hasher.Sum(nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.MongoClient(ctx)
	defer cancel()

	if err != nil {
		log.Panicf("Register DB Connection Error:\n%v", err)
	}

	var result map[string]interface{}
	collection := client.Database(os.Getenv("PROJECT_NAME")).Collection("users")
	collection.FindOne(ctx, bson.D{
		{
			Key: "email", Value: bodyData["email"],
		},
	}).Decode(&result)

	if len(result) == 0 {
		collection = client.Database(os.Getenv("PROJECT_NAME")).Collection("users")
		res, err := collection.InsertOne(ctx, bson.D{
			{
				Key: "username", Value: fmt.Sprintf("%s", bodyData["username"]),
			}, {
				Key: "email", Value: fmt.Sprintf("%s", bodyData["email"]),
			}, {
				Key: "password", Value: sha256_hash,
			},
			{
				Key: "privilege", Value: privilege,
			},
		})

		if err != nil {
			log.Panicf("Register InsertOne Failed:\n%v", err)
		}
		log.Printf("Register InsertOne Respone: %v", res)

		c.JSON(200, gin.H{
			"msg":  bodyData,
			"hash": sha256_hash,
		})

		if tokenInfo != nil {
			collection = client.Database(os.Getenv("PROJECT_NAME")).Collection("registerToken")
			_, err := collection.DeleteOne(ctx, bson.D{
				{
					Key: "token", Value: fmt.Sprintf("%s", tokenInfo["token"]),
				},
			})

			if err != nil {
				log.Panicf("Register DeleteOne Error:\n%v", err)
			}
		}
	}

	c.JSON(409, gin.H{
		"msg": "User Already Exists.",
	})
}
