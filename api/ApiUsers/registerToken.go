package ApiUsers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterToken(c *gin.Context) {
	_, isAdmin := c.Get("admin")

	if !isAdmin {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	}

	bodyRaw, err := c.GetRawData()

	if err != nil {
		log.Panicf("RegisterToken Get Body Request Failed:\n%v", err)
	}

	var bodyData map[string]interface{}
	json.Unmarshal(bodyRaw, &bodyData)

	if bodyData["privilege"] == nil {
		c.JSON(400, "400 Bad Request")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.MongoClient(ctx)
	defer cancel()

	if err != nil {
		c.JSON(500, "500 Internal Server Error")
		log.Panicf("BeforeRegister DB Connection Error:\n%v", err)
	}

	collection := client.Database(fmt.Sprintf("%s", os.Getenv("PROJECT_NAME"))).Collection("registerToken")
	var expireAfterSeconds int32 = 43200
	createIndexRes, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"createdAt", 1}}, Options: &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds}})

	if err != nil {
		dropIndexRes, err := collection.Indexes().DropOne(ctx, "createdAt_1")
		if err != nil {
			c.JSON(500, "500 Internal Server Error")
			log.Panicf("RegisterToken DropIndex Failed:\n%v", err)
		}
		log.Printf("RegisterToken DropIndex Response: ", dropIndexRes)

		createIndexRes, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"createdAt", 1}}, Options: &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds}})
		if err != nil {
			c.JSON(500, "500 Internal Server Error")
			log.Panicf("RegisterToken CreateIndex Failed:\n%v", err)
		}
	}
	log.Printf("RegisterToken CreateIndex Response: ", createIndexRes)

	res, err := collection.InsertOne(ctx, bson.D{
		{
			"createdAt", time.Now(),
		},
		{
			"token", tokenGenerator(),
		},
		{
			"privilege", bodyData["privilege"],
		},
	})

	if err != nil {
		c.JSON(500, "500 Internal Server Error")
		log.Panicf("RegisterToken InsertOne Failed:\n%v", err)
	}
	log.Printf("RegisterToken InsertOne Respone: %v", res)

	c.JSON(201, gin.H{
		"token": res,
	})
}

func tokenGenerator() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
