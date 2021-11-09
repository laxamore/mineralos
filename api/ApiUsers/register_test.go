package ApiUsers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRepositoryMock struct {
	mock.Mock
}

func (a RegisterRepositoryMock) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	users := []map[string]interface{}{{
		"username":  "testexist",
		"email":     "test@testexist.com",
		"password":  "937e8d5fbb48bd4949536cd65b8d35c426b80d2f830c5c308e2cdec422ae2244",
		"privilege": "admin",
	}}

	var input map[string]interface{}

	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	for i := range users {
		if input["username"] == users[i]["username"] || input["email"] == users[i]["email"] {
			return users[i]
		}
	}

	return map[string]interface{}{}
}

func (a RegisterRepositoryMock) InsertOne(db_name string, collection_name string, filter interface{}) (*mongo.InsertOneResult, error) {
	var insertOneResult *mongo.InsertOneResult
	return insertOneResult, nil
}

func (a RegisterRepositoryMock) DeleteOne(db_name string, collection_name string, filter interface{}) (*mongo.DeleteResult, error) {
	registerToken := []map[string]interface{}{{
		"createdAt": time.Now(),
		"token":     "testtesttesttest",
		"privilege": "readOnly",
	}}

	var DeleteResult *mongo.DeleteResult
	var input map[string]interface{}
	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	for i := range registerToken {
		if registerToken[i]["token"] == input["token"] {
			return DeleteResult, nil
		}
	}

	return DeleteResult, errors.New("DeleteOne Failed:\n")
}

func (a RegisterRepositoryMock) IndexesReplaceMany(db_name string, collection_name string, indexModel []mongo.IndexModel) ([]string, error) {
	return []string{}, nil
}

func TestRegister(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		bodyData     map[string]interface{}
		token        map[string]interface{}
	}

	testData := []TestData{
		{
			testName:     "SuccessRegister",
			expectedCode: http.StatusOK,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@test.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "TokenFailed",
			expectedCode: http.StatusNotFound,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@test.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "WithExistingUsername",
			expectedCode: http.StatusConflict,
			bodyData: map[string]interface{}{
				"username": "testexist",
				"email":    "test@test.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "WithExistingEmail",
			expectedCode: http.StatusConflict,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "test@testexist.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "InvalidEmail",
			expectedCode: http.StatusBadRequest,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "testemail.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "InvalidUsername",
			expectedCode: http.StatusBadRequest,
			bodyData: map[string]interface{}{
				"username": "!@test",
				"email":    "testemail.com",
				"password": "test1234",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
			},
		},
		{
			testName:     "InvalidPassword",
			expectedCode: http.StatusBadRequest,
			bodyData: map[string]interface{}{
				"username": "test",
				"email":    "testemail.com",
				"password": "test123",
			},
			token: map[string]interface{}{
				"createdAt": time.Now(),
				"token":     "testtesttesttest",
				"privilege": "readOnly",
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

			c.Set("token", td.token)

			jsonbytes, err := json.Marshal(td.bodyData)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			repo := RegisterRepositoryMock{}
			cntrl := RegisterController{}

			cntrl.TryRegister(c, repo)

			t.Logf("Response Body: %s", w.Body.String())
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
