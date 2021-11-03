package Middleware

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
)

type BeforeRegisterRepositoryMock struct {
	mock.Mock
	emptyUsers bool
}

func (a BeforeRegisterRepositoryMock) Find(db_name string, collection_name string, filter interface{}) ([]map[string]interface{}, error) {
	dumyUsers := []map[string]interface{}{
		{
			"username":  "test",
			"email":     "test@test.com",
			"password":  "937e8d5fbb48bd4949536cd65b8d35c426b80d2f830c5c308e2cdec422ae2244",
			"privilege": "admin",
		},
	}

	if a.emptyUsers {
		return []map[string]interface{}{{}}, nil
	}
	return dumyUsers, nil
}

func (a BeforeRegisterRepositoryMock) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	token := map[string]interface{}{
		"token": "testtesttesttest",
	}

	var input map[string]interface{}
	filterbytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterbytes, &input)

	if input["token"] == token["token"] {
		return token
	}
	return map[string]interface{}{}
}

func TestBeforeRegister(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		bodyData     map[string]interface{}
		emptyUsers   bool
	}

	testData := []TestData{
		{
			testName:     "SuccesBeforeRegister",
			expectedCode: http.StatusOK,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@test.com",
				"password": "test1234",
				"token":    "testtesttesttest",
			},
			emptyUsers: false,
		},
		{
			testName:     "SuccesBeforeRegisterAdmin",
			expectedCode: http.StatusOK,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@test.com",
				"password": "test1234",
				"token":    "testtesttesttest",
			},
			emptyUsers: true,
		},
		{
			testName:     "FailedBeforeRegister",
			expectedCode: http.StatusUnauthorized,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@test.com",
				"password": "test1234",
				"token":    "testtesttesttess",
			},
			emptyUsers: true,
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

			repo := BeforeRegisterRepositoryMock{}
			cntrl := BeforeRegisterController{}

			repo.emptyUsers = td.emptyUsers

			cntrl.TryBeforeRegister(c, repo)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
