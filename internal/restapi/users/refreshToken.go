package users

//
//import (
//	"fmt"
//	"github.com/laxamore/mineralos/internal/db"
//	JWT "github.com/laxamore/mineralos/internal/jwt"
//	"github.com/laxamore/mineralos/internal/logger"
//	"github.com/laxamore/mineralos/internal/restapi"
//	"net/http"
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
//type RefreshTokenRepositoryInterface interface {
//	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
//}
//
//type RefreshTokenController struct{}
//
//func (a RefreshTokenController) TryRefreshToken(c *gin.Context, client *mongo.Client, repositoryInterface RefreshTokenRepositoryInterface) {
//	response := restapi.Result{
//		Code:     http.StatusBadRequest,
//		Response: "Bad Request",
//	}
//
//	rtoken, _ := c.Cookie("rtoken")
//
//	if rtoken != "" {
//		claims, tokenParsed, err := JWT.VerifyJWT(rtoken)
//
//		if tokenParsed == nil || err != nil {
//			logger.Printf("Couldn't handle this token: %v", err)
//		} else if tokenParsed.Valid {
//			objectID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", claims["id"]))
//			if err != nil {
//				logger.Print("Invalid id")
//			}
//
//			result := repositoryInterface.FindOne(client, "mineralos", "users", bson.D{
//				{
//					Key: "_id", Value: objectID,
//				},
//			})
//
//			if len(result) > 0 {
//				exp, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRED"), 10, 64)
//				if err != nil {
//					logger.Printf("ACCESS_TOKEN_EXPIRED env is not number %v", err)
//				}
//				exp = exp + time.Now().Unix()
//
//				// exp = exp + time.Now().Unix()exp := time.Now().Unix() + 60
//
//				newClaims := jwt.MapClaims{
//					"username":  result["username"],
//					"email":     result["email"],
//					"privilege": result["privilege"],
//				}
//				newToken, err := JWT.SignJWT(newClaims, exp)
//
//				if err != nil {
//					logger.Panicf("Login Token Sign Failed:\n%v", err)
//				}
//
//				response.Code = http.StatusOK
//				response.Response = gin.H{
//					"jwt_token":        newToken,
//					"jwt_token_expiry": exp,
//				}
//				c.JSON(response.Code, response.Response)
//				return
//			}
//		} else {
//			logger.Printf("Couldn't handle this token: %v", err)
//		}
//	}
//
//	c.JSON(response.Code, response.Response)
//}
//
//func RefreshToken(client *mongo.Client) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		repo := db.MongoDB{}
//		cntrl := RefreshTokenController{}
//
//		cntrl.TryRefreshToken(c, client, repo)
//	}
//}
