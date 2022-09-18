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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteRigRepositoryMock struct {
	mock.Mock
}

func (a DeleteRigRepositoryMock) DeleteOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) (*mongo.DeleteResult, error) {
	dummyRigData := map[string]interface{}{
		"rig_id": "1234567",
	}

	var input map[string]interface{}
	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if dummyRigData["rig_id"] == input["rig_id"] {
		return nil, nil
	}

	return nil, fmt.Errorf("error rig_id not found")
}

func TestDeleteRig(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		privilege    string
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "SuccessDeleteRig",
			expectedCode: http.StatusOK,
			privilege:    "admin",
			bodyData: gin.H{
				"rig_id": "1234567",
			},
		},
		{
			testName:     "FailedDeleteRigNotFound",
			expectedCode: http.StatusNotFound,
			privilege:    "admin",
			bodyData: gin.H{
				"rig_id": "7654321",
			},
		},
		{
			testName:     "FailedDeleteRigForbidden",
			expectedCode: http.StatusForbidden,
			privilege:    "readOnly",
			bodyData: gin.H{
				"rig_id": "1234567",
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

			repo := DeleteRigRepositoryMock{}
			cntrl := DeleteRigController{}

			cntrl.TryDeleteRig(c, nil, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
