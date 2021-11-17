package ApiRigs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetRigRepositoryMock struct {
	mock.Mock
}

func (a GetRigRepositoryMock) FindOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) map[string]interface{} {
	dumyData := map[string]interface{}{
		"rig_id": "testtesttest",
	}

	var input map[string]interface{}
	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if input["rig_id"] == dumyData["rig_id"] {
		return dumyData
	}

	return map[string]interface{}{}
}

func TestGetRig(t *testing.T) {
	type TestData struct {
		testName     string
		rig_id       string
		expectedCode int
	}

	testData := []TestData{
		{
			testName:     "SuccessGetRig",
			rig_id:       "testtesttest",
			expectedCode: http.StatusOK,
		},
		{
			testName:     "FailGetRig",
			rig_id:       "failfailfail",
			expectedCode: http.StatusNotFound,
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

			c.Params = []gin.Param{
				{
					Key:   "rig_id",
					Value: td.rig_id,
				},
			}

			repo := GetRigRepositoryMock{}
			cntrl := GetRigController{}

			cntrl.TryGetRig(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
