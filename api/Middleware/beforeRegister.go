package Middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
)

func BeforeRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyRaw, err := c.GetRawData()

		if err != nil {
			log.Panicf("BeforeRegister Get Body Request Failed:\n%v", err)
		}

		var bodyData map[string]interface{}
		json.Unmarshal(bodyRaw, &bodyData)
		c.Set("bodyData", bodyData)

		if bodyData["username"] != nil && bodyData["email"] != nil && bodyData["password"] != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			client, err := db.MongoClient(ctx)
			defer cancel()

			if err != nil {
				log.Panicf("BeforeRegister DB Connection Error:\n%v", err)
			}

			collection := client.Database(os.Getenv("PROJECT_NAME")).Collection("users")
			cur, err := collection.Find(ctx, bson.D{{}})

			if err != nil {
				log.Panicf("BeforeRegister List All Users Failed:\n%v", err)
			}

			var results []bson.D

			for cur.Next(ctx) {
				//Create a value into which the single document can be decoded
				var elem bson.D
				err := cur.Decode(&elem)
				if err != nil {
					log.Fatal(err)
				}

				results = append(results, elem)
			}

			if len(results) > 0 {
				if bodyData["token"] != nil {
					var result map[string]interface{}
					collection = client.Database(os.Getenv("PROJECT_NAME")).Collection("registerToken")
					collection.FindOne(ctx, bson.D{{Key: "token", Value: bodyData["token"]}}).Decode(&result)

					log.Print(result)

					if len(result) > 0 {
						c.Set("token", result)
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

}
