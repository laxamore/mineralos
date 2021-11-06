package ApiUsers

import (
	"fmt"
	"net/http"
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

type RefreshTokenRepositoryInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
}

type RefreshTokenController struct{}

func (a RefreshTokenController) TryRefreshToken(c *gin.Context, repositoryInterface RefreshTokenRepositoryInterface) {
	response := api.Result{
		Code:     http.StatusBadRequest,
		Response: "Bad Request",
	}

	rtoken := c.GetHeader("rt")

	if rtoken != "" {
		claims, tokenParsed, err := JWT.VerifyJWT(rtoken)

		if tokenParsed == nil || err != nil {
			Log.Printf("Couldn't handle this token: %v", err)
		} else if tokenParsed.Valid {
			objectID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", claims["id"]))
			if err != nil {
				Log.Print("Invalid id")
			}

			result := repositoryInterface.FindOne(os.Getenv("PROJECT_NAME"), "users", bson.D{
				{
					Key: "_id", Value: objectID,
				},
			})

			if len(result) > 0 {
				exp := time.Now().Unix() + 300
				newClaims := jwt.MapClaims{
					"username":  result["username"],
					"email":     result["email"],
					"privilege": result["privilege"],
				}
				newToken, err := JWT.SignJWT(newClaims, exp)

				if err != nil {
					Log.Panicf("Login Token Sign Failed:\n%v", err)
				}

				response.Code = http.StatusOK
				response.Response = gin.H{
					"jwt_token":        newToken,
					"jwt_token_expiry": exp,
				}
				c.JSON(response.Code, response.Response)
				return
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				Log.Print("Couldn't handle this token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				Log.Print("Token Expired.")
			} else {
				Log.Printf("Couldn't handle this token: %v", err)
			}
		} else {
			Log.Printf("Couldn't handle this token: %v", err)
		}
	}

	c.JSON(response.Code, response.Response)
}

func RefreshToken(c *gin.Context) {
	repo := db.MongoDB{}
	cntrl := RefreshTokenController{}

	cntrl.TryRefreshToken(c, repo)
}
