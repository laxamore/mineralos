package ApiUsers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/utils/JWT"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenRepositoryMock struct {
	mock.Mock
}

func (a RefreshTokenRepositoryMock) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	var input map[string]interface{}
	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if input["_id"].(primitive.ObjectID).Hex() == "61866b3920b512a8608788ad" {
		return map[string]interface{}{
			"username":  "test",
			"email":     "test@test.com",
			"privilege": "readOnly",
		}
	}

	return map[string]interface{}{}
}

func TestRefreshToken(t *testing.T) {
	type TestData struct {
		testName string
		expected map[string]interface{}
		id       string
	}

	testData := []TestData{
		{
			testName: "SuccessRefreshToken",
			expected: map[string]interface{}{
				"Code": http.StatusOK,
			},
			id: "61866b3920b512a8608788ad",
		},
		{
			testName: "FailRefreshToken",
			expected: map[string]interface{}{
				"Code": http.StatusBadRequest,
				"Body": "Bad Request",
			},
			id: "61866b5bdb2cbb26a1e06c7a",
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			// Create test jwt
			claims := jwt.MapClaims{
				"id": td.id,
			}
			rtoken, err := JWT.SignJWT(claims, time.Now().Unix()+5)
			if err != nil {
				t.Fatalf("Login Token Sign Failed:\n%v", err)
			}
			//

			c.Request.Header.Set("rt", rtoken)

			repo := RefreshTokenRepositoryMock{}
			cntrl := RefreshTokenController{}

			cntrl.TryRefreshToken(c, repo)

			assert.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expected["Code"]), fmt.Sprintf("HTTP Status Code: %d", w.Code))

			if w.Code != 200 {
				var responseBody string
				json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.EqualValues(t, td.expected["Body"], responseBody)
			} else {
				var responseBody map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NotNil(t, responseBody["jwt_token"])
				require.NotNil(t, responseBody["jwt_token_expiry"])
			}
		})
	}
}
