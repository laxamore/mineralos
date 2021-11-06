package ApiUsers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/JWT"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRepositoryInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
}

type LoginController struct{}

func (a LoginController) TryLogin(c *gin.Context, repositoryInterface LoginRepositoryInterface) {
	var response api.Result
	bodyRaw, err := c.GetRawData()

	if err != nil {
		Log.Panicf("Login Get Body Request Failed:\n%v", err)
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
		exp := time.Now().Unix() + 300
		expRT := time.Now().Unix() + 2592000 // 1 Month Expired

		// Create the Claims for token
		claims := jwt.MapClaims{
			"username":  result["username"],
			"email":     result["email"],
			"privilege": result["privilege"],
		}
		token, err := JWT.SignJWT(claims, exp)
		if err != nil {
			Log.Panicf("Login Token Sign Failed:\n%v", err)
		}
		//

		// Create the claims for refresh token
		claims = jwt.MapClaims{
			"id": result["_id"].(primitive.ObjectID).Hex(),
		}
		rtoken, err := JWT.SignJWT(claims, expRT)

		if err != nil {
			Log.Panicf("Login Token Sign Failed:\n%v", err)
		}

		response.Code = 200
		response.Response = gin.H{
			"jwt_token":        token,
			"jwt_token_expiry": exp,
			"r_token":          rtoken,
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
