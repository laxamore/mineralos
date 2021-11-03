package ApiUsers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginRepositoryInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
}

type LoginController struct{}

func (a LoginController) TryLogin(c *gin.Context, repositoryInterface LoginRepositoryInterface) {
	var response api.Result
	bodyRaw, err := c.GetRawData()

	if err != nil {
		log.Panicf("Login Get Body Request Failed:\n%v", err)
	}

	var bodyData map[string]interface{}
	json.Unmarshal(bodyRaw, &bodyData)

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s", bodyData["password"])))
	sha256_hash := hex.EncodeToString(hasher.Sum(nil))

	result := repositoryInterface.FindOne(os.Getenv("PROJECT_NAME"), "users", bson.D{
		{
			Key: "username", Value: fmt.Sprintf("%s", bodyData["username"]),
		},
		{
			Key: "password", Value: sha256_hash,
		},
	})

	if len(result) != 0 {
		// Create the Claims
		claims := jwt.MapClaims{
			"username":  result["username"],
			"email":     result["email"],
			"privilege": result["privilege"],
			"exp":       time.Now().Unix() + 2592000, // Expired After 1 Month
			"iat":       time.Now().Unix(),
		}

		log.Print(time.Now().Unix())

		signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := signToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			log.Panicf("Login Token Sign Failed:\n%v", err)
		}

		response.Code = 200
		response.Response = gin.H{
			"jwt": token,
		}
	} else {
		response.Code = 401
		response.Response = "Login Failed"
	}

	c.JSON(response.Code, response.Response)
}

func Login(c *gin.Context) {
	repo := db.MongoDB{}
	cntrl := LoginController{}

	cntrl.TryLogin(c, repo)
}
