package rigs

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

type NewWalletRepositoryMock struct {
	mock.Mock
}

func (a NewWalletRepositoryMock) InsertOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (*mongo.InsertOneResult, error) {
	InsertID := mongo.InsertOneResult{
		InsertedID: "testinsertid",
	}
	return &InsertID, nil
}

func TestNewWallet(t *testing.T) {
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
				"wallet_name":    "test",
				"wallet_address": "0x9062bd26e6a84086634bb4322cd526f368e4e402",
				"coin":           "eth",
			},
		},
		{
			testName:     "FailedNewRig",
			expectedCode: http.StatusForbidden,
			privilege:    "readOnly",
			bodyData: gin.H{
				"wallet_name":    "test",
				"wallet_address": "0x9062bd26e6a84086634bb4322cd526f368e4e402",
				"coin":           "eth",
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

			ctrl := RigController{}

			ctrl.NewWallet(c)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
