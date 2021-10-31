package ApiUsers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *gin.Context) {
	bodyRaw, err := c.GetRawData()

	if err != nil {
		log.Panicf("Login Get Body Request Failed:\n%v", err)
	}

	var bodyData map[string]interface{}
	json.Unmarshal(bodyRaw, &bodyData)

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s", bodyData["password"])))
	sha256_hash := hex.EncodeToString(hasher.Sum(nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	collection := client.Database(fmt.Sprintf("%s", os.Getenv("PROJECT_NAME"))).Collection("users")

	var result map[string]interface{}
	collection.FindOne(ctx, bson.D{
		{
			"username", fmt.Sprintf("%s", bodyData["username"]),
		},
		{
			"password", sha256_hash,
		},
	}).Decode(&result)

	if len(result) != 0 {
		token := generateJWT(result)
		c.JSON(200, map[string]interface{}{
			"jwt": token,
		})
	} else {
		c.JSON(401, "Login Failed")
	}
}

func generateJWT(loginInfo map[string]interface{}) string {
	// Create the Claims
	claims := jwt.MapClaims{
		"username":  loginInfo["username"],
		"email":     loginInfo["email"],
		"privilege": loginInfo["privilege"],
		"exp":       time.Now().Unix() + 2592000, // Expired After 1 Month
		"iat":       time.Now().Unix(),
	}

	log.Print(time.Now().Unix())

	signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Panicf("Login Token Sign Failed:\n%v", err)
	}
	return token
}
