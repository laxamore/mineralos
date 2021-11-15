package ApiRigs

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

type NewRigRepositoryMock struct {
	mock.Mock
}

func (a NewRigRepositoryMock) InsertOne(client *mongo.Client, db_name string, collection_name string, fitler interface{}) (*mongo.InsertOneResult, error) {
	InsertID := mongo.InsertOneResult{
		InsertedID: "testinsertid",
	}
	return &InsertID, nil
}

func TestNewRig(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		privilege    string
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "SuccessNewRIG",
			expectedCode: http.StatusOK,
			privilege:    "admin",
			bodyData: gin.H{
				"rig_name": "test",
			},
		},
		{
			testName:     "FailedNewRig",
			expectedCode: http.StatusForbidden,
			privilege:    "readOnly",
			bodyData: gin.H{
				"rig_name": "test",
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

			c.Set("tokenClaims", map[string]interface{}{
				"privilege": td.privilege,
			})

			repo := NewRigRepositoryMock{}
			cntrl := NewRigController{}

			cntrl.TryNewRig(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
