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
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterTokenRepositoryMock struct {
	mock.Mock
}

func (a RegisterTokenRepositoryMock) InsertOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (*mongo.InsertOneResult, error) {
	var InsertOneResult *mongo.InsertOneResult
	return InsertOneResult, nil
}

func (a RegisterTokenRepositoryMock) IndexesReplaceOne(client *mongo.Client, db_name string, collection_name string, indexModel mongo.IndexModel) (string, error) {
	return "", nil
}

func TestRegisterToken(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		bodyData     map[string]interface{}
		admin        bool
	}

	testData := []TestData{
		{
			testName:     "SuccessRegisterToken",
			expectedCode: http.StatusCreated,
			bodyData: gin.H{
				"privilege": "readOnly",
			},
			admin: true,
		},
		{
			testName:     "PrivilegeUndefine",
			expectedCode: http.StatusBadRequest,
			bodyData:     gin.H{},
			admin:        true,
		},
		{
			testName:     "NonAdmin",
			expectedCode: http.StatusUnauthorized,
			bodyData:     gin.H{},
			admin:        false,
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
			c.Set("admin", td.admin)

			jsonbytes, err := json.Marshal(td.bodyData)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			repo := RegisterTokenRepositoryMock{}
			cntrl := RegisterTokenController{}

			cntrl.TryRegisterToken(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
