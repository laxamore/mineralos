package ApiUsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRepositoryMock struct {
	mock.Mock
}

func (r LoginRepositoryMock) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	ObjectID, _ := primitive.ObjectIDFromHex("61866b3920b512a8608788ad")
	users := map[string]interface{}{
		"_id":       ObjectID,
		"username":  "test",
		"email":     "test@test.com",
		"password":  "937e8d5fbb48bd4949536cd65b8d35c426b80d2f830c5c308e2cdec422ae2244",
		"privilege": "admin",
	}

	var input map[string]interface{}

	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if (input["username"] == users["username"] && input["password"] == users["password"]) ||
		(input["email"] == users["email"] && input["password"] == users["password"]) {
		return users
	}

	return nil
}

func TestLogin(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "SuccessLogin",
			expectedCode: http.StatusOK,
			bodyData: gin.H{
				"username": "test",
				"password": "test1234",
			},
		},
		{
			testName:     "FailedLogin",
			expectedCode: http.StatusUnauthorized,
			bodyData: gin.H{
				"username": "test",
				"password": "test4321",
			},
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

			jsonbytes, err := json.Marshal(td.bodyData)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			repo := LoginRepositoryMock{}
			cntrl := LoginController{}

			cntrl.TryLogin(c, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
