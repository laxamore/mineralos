package rigs

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//type UpdateOCRepositoryMock struct {
//	mock.Mock
//}
//
//func (a UpdateOCRepositoryMock) UpdateOne(client *mongo.Client, db_name string, collection_name string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	dumyData := map[string]interface{}{
//		"rig_id": "testtesttest",
//	}
//
//	var input map[string]interface{}
//	filterBytes, _ := bson.Marshal(filter)
//	bson.Unmarshal(filterBytes, &input)
//
//	if input["rig_id"] == dumyData["rig_id"] {
//		return nil, nil
//	}
//
//	return nil, fmt.Errorf("error")
//}
//
//func TestUpdateOC(t *testing.T) {
//	type TestData struct {
//		testName     string
//		expectedCode int
//		privilege    string
//		bodyData     UpdateOCProps
//	}
//
//	testData := []TestData{
//		{
//			testName:     "SuccessUpdateOC",
//			expectedCode: http.StatusOK,
//			privilege:    "admin",
//			bodyData: UpdateOCProps{
//				RIG_ID: "testtesttest",
//				Vendor: "AMD",
//				ID:     0,
//			},
//		},
//		{
//			testName:     "FailUpdateOC",
//			expectedCode: http.StatusForbidden,
//			privilege:    "admin",
//			bodyData: UpdateOCProps{
//				RIG_ID: "failfailfail",
//				Vendor: "AMD",
//				ID:     0,
//			},
//		},
//		{
//			testName:     "FailReadOnly",
//			expectedCode: http.StatusForbidden,
//			privilege:    "readOnly",
//			bodyData: UpdateOCProps{
//				RIG_ID: "testtesttest",
//				Vendor: "AMD",
//				ID:     0,
//			},
//		},
//	}
//
//	for _, td := range testData {
//		t.Run(td.testName, func(t *testing.T) {
//			gin.SetMode(gin.TestMode)
//			w := httptest.NewRecorder()
//			c, _ := gin.CreateTestContext(w)
//
//			c.Request = &http.Request{
//				Header: make(http.Header),
//			}
//
//			jsonbytes, err := json.Marshal(td.bodyData)
//			if err != nil {
//				panic(err)
//			}
//			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
//
//			c.Set("tokenClaims", map[string]interface{}{
//				"privilege": td.privilege,
//			})
//
//			repo := UpdateOCRepositoryMock{}
//			cntrl := UpdateOCController{}
//
//			cntrl.TryUpdateOC(c, nil, repo)
//			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
//		})
//	}
//}
