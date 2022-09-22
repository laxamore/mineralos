package users

//
//import (
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"github.com/laxamore/mineralos/internal/db"
//	JWT "github.com/laxamore/mineralos/internal/jwt"
//	"github.com/laxamore/mineralos/internal/logger"
//	"github.com/laxamore/mineralos/internal/restapi"
//	"net/http"
//	"net/url"
//	"os"
//	"strconv"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//type LoginRepositoryInterface interface {
//	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
//}
//
//type LoginController struct{}
//
//func (a LoginController) TryLogin(c *gin.Context, client *mongo.Client, repositoryInterface LoginRepositoryInterface) {
//	var response restapi.Result
//	bodyRaw, err := c.GetRawData()
//
//	if err != nil {
//		logger.Panicf("Login Get Body Request Failed:\n%v", err)
//	}
//
//	var bodyData map[string]interface{}
//	json.Unmarshal(bodyRaw, &bodyData)
//
//	hasher := sha256.New()
//	hasher.Write([]byte(fmt.Sprintf("%s", bodyData["password"])))
//	sha256_hash := hex.EncodeToString(hasher.Sum(nil))
//
//	result := repositoryInterface.FindOne(client, "mineralos", "users", bson.D{
//		{
//			Key: "username", Value: fmt.Sprintf("%s", bodyData["username"]),
//		},
//		{
//			Key: "password", Value: sha256_hash,
//		},
//	})
//
//	if len(result) != 0 {
//		exp, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRED"), 10, 64)
//		if err != nil {
//			logger.Printf("ACCESS_TOKEN_EXPIRED env is not number %v", err)
//		}
//		exp = exp + time.Now().Unix()
//
//		expRT, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRED"), 10, 64)
//		if err != nil {
//			logger.Printf("REFRESH_TOKEN_EXPIRED env is not number %v", err)
//		}
//		expRT = expRT + time.Now().Unix()
//
//		// exp := time.Now().Unix() + os.Getenv("ACCESS_TOKEN_EXPIRED")
//		// expRT := time.Now().Unix() + strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRED")) // 1 Month Expired
//
//		// Create the Claims for token
//		claims := jwt.MapClaims{
//			"username":  result["username"],
//			"email":     result["email"],
//			"privilege": result["privilege"],
//		}
//		token, err := JWT.SignJWT(claims, exp)
//		if err != nil {
//			logger.Panicf("Login Token Sign Failed:\n%v", err)
//		}
//		//
//
//		// Create the claims for refresh token
//		claims = jwt.MapClaims{
//			"id": result["_id"].(primitive.ObjectID).Hex(),
//		}
//		rtoken, err := JWT.SignJWT(claims, expRT)
//
//		if err != nil {
//			logger.Panicf("Login Token Sign Failed:\n%v", err)
//		}
//
//		http.SetCookie(c.Writer, &http.Cookie{
//			Name:     "rtoken",
//			Value:    url.QueryEscape(rtoken),
//			MaxAge:   2592000,
//			Path:     "/",
//			Domain:   os.Getenv("DOMAIN"),
//			SameSite: http.SameSiteStrictMode,
//			Secure:   true,
//			HttpOnly: true,
//		})
//
//		response.Code = 200
//		response.Response = gin.H{
//			"jwt_token":        token,
//			"jwt_token_expiry": exp,
//			// "r_token":          rtoken,
//		}
//	} else {
//		response.Code = 401
//		response.Response = "Login Failed"
//	}
//
//	c.JSON(response.Code, response.Response)
//}
//
//func Login(client *mongo.Client) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		repo := db.MongoDB{}
//		cntrl := LoginController{}
//
//		cntrl.TryLogin(c, client, repo)
//	}
//}
