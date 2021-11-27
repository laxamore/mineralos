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
	"github.com/laxamore/mineralos/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteWalletRepositoryMock struct {
	mock.Mock
}

func (a DeleteWalletRepositoryMock) DeleteOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex("61a1b7f9a51371bd796efcb0")
	utils.CheckErr(err)

	dummyWalletData := map[string]interface{}{
		"wallet_id": objectId,
	}

	var input map[string]interface{}
	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if dummyWalletData["wallet_id"] == input["_id"] {
		return nil, nil
	}

	return nil, fmt.Errorf("error wallet_id not found")
}

func TestDeleteWallet(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		privilege    string
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "SuccessDeleteWallet",
			expectedCode: http.StatusOK,
			privilege:    "admin",
			bodyData: gin.H{
				"wallet_id": "61a1b7f9a51371bd796efcb0",
			},
		},
		{
			testName:     "FailedDeleteWalletNotFound",
			expectedCode: http.StatusNotFound,
			privilege:    "admin",
			bodyData: gin.H{
				"wallet_id": "61a1b80a4627b335bbe5bf0a",
			},
		},
		{
			testName:     "FailedDeleteWalletForbidden",
			expectedCode: http.StatusForbidden,
			privilege:    "readOnly",
			bodyData: gin.H{
				"wallet_id": "61a1b7f9a51371bd796efcb0",
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

			repo := DeleteWalletRepositoryMock{}
			cntrl := DeleteWalletController{}

			cntrl.TryDeleteWallet(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
