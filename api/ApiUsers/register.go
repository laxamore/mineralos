package ApiUsers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegisterRepositoryInterface interface {
	FindOne(string, string, interface{}) map[string]interface{}
	InsertOne(string, string, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(string, string, interface{}) (*mongo.DeleteResult, error)
	IndexesReplaceMany(string, string, []mongo.IndexModel) ([]string, error)
}

type RegisterController struct{}

func (a RegisterController) TryRegister(c *gin.Context, repositoryInterface RegisterRepositoryInterface) {
	var response api.Result
	response.Code = http.StatusBadRequest
	response.Response = "Bad Request"

	errMsg := ""
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

	_, err := mail.ParseAddress(fmt.Sprintf("%s", bodyData["email"]))
	isEmailValid := (err == nil)

	regexUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)
	correctUsername := regexUsername.Match([]byte(fmt.Sprintf("%s", bodyData["username"])))
	corectPasswordLength := (len(fmt.Sprintf("%s", bodyData["password"])) >= 8)

	if isEmailValid && corectPasswordLength && correctUsername {
		hasher := sha256.New()
		hasher.Write([]byte(fmt.Sprintf("%s", bodyData["password"])))
		sha256_hash := hex.EncodeToString(hasher.Sum(nil))

		// Create/Replace MongoDB Indexes For Unique Username & Email
		isUnique := true
		createIndexRes, err := repositoryInterface.IndexesReplaceMany(os.Getenv("PROJECT_NAME"), "users", []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "username", Value: 1},
				},
				Options: &options.IndexOptions{
					Unique: &isUnique,
				},
			},
			{
				Keys: bson.D{
					{Key: "email", Value: 1},
				},
				Options: &options.IndexOptions{
					Unique: &isUnique,
				},
			},
		})

		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Response = "InternalServerError"
			errMsg = "failed to create/replace username indexes."
		} else {
			Log.Printf("Username CreateIndex/ReplaceIndex Response: %v", createIndexRes)
		}
		//

		// Check Users & Email Already Exists
		result := map[string]interface{}{}
		checkIndex := []string{
			"username",
			"email",
		}
		for _, check := range checkIndex {
			if len(result) == 0 {
				result = repositoryInterface.FindOne(os.Getenv("PROJECT_NAME"), "users", bson.D{
					{
						Key: check, Value: fmt.Sprintf("%s", bodyData[check]),
					},
				})
			} else {
				break
			}
		}
		//

		if len(result) == 0 {
			if tokenInfo != nil {
				// Delete Register Token
				_, err := repositoryInterface.DeleteOne(os.Getenv("PROJECT_NAME"), "registerToken", bson.D{
					{
						Key: "token", Value: fmt.Sprintf("%s", tokenInfo["token"]),
					},
				})
				Log.Printf("%v", err)
				//

				if err != nil {
					response.Code = http.StatusNotFound
					response.Response = "RegisterToken Not Found"
					errMsg = fmt.Sprintf("Register DeleteOne Error:\n%v", err)
				} else {
					// Register New User
					res, err := repositoryInterface.InsertOne(os.Getenv("PROJECT_NAME"), "users", bson.D{
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
						response.Code = http.StatusBadRequest
						response.Response = "Register Failed Bad Request"
						errMsg = fmt.Sprintf("Register InsertOne Failed:\n%v", err)
					} else {
						response.Code = http.StatusOK
						response.Response = gin.H{
							"msg":  bodyData,
							"hash": sha256_hash,
						}
						log.Printf("Register InsertOne Respone: %v", res)
					}
					//
				}
			}
		} else {
			response.Code = http.StatusConflict
			response.Response = "User Already Exists."
			errMsg = "User Already Exists."
		}
	}

	c.JSON(response.Code, response.Response)
	if errMsg != "" {
		Log.Panic(errMsg)
	}
}

func Register(c *gin.Context) {
	repo := db.MongoDB{}
	cntrl := RegisterController{}

	cntrl.TryRegister(c, repo)
}
