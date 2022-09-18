package rigs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetWalletsRepositoryMock struct {
	mock.Mock
}

func (a GetWalletsRepositoryMock) Find(client *mongo.Client, db_name string, collection_name string, filter interface{}) ([]map[string]interface{}, error) {
	dummyRigsData := []map[string]interface{}{
		{
			"wallet_name":    "test",
			"wallet_address": "0x9062bd26e6a84086634bb4322cd526f368e4e402",
			"coin":           "eth",
		},
	}

	return dummyRigsData, nil
}

func TestGetWallets(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
	}

	testData := []TestData{
		{
			testName:     "SuccessGetRigs",
			expectedCode: http.StatusOK,
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

			repo := GetWalletsRepositoryMock{}
			cntrl := GetWalletsController{}

			cntrl.TryGetWallets(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
