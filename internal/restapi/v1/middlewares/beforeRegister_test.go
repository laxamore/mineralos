package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/restapi/v1/users"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type beforeRegisterMock struct {
	mock.Mock
	db.IDB
	db.IRedis
}

func (m beforeRegisterMock) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
	//if isAdminRegister {
	//	return &gorm.DB{}
	//} else {
	//	dest.(*models.User).Username = "test"
	//	return &gorm.DB{}
	//}
}

func (m beforeRegisterMock) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
	//stringCmd := &redis.StringCmd{}
	//
	//if key == "testtesttesttest" {
	//	registerTokenTest := models.RegisterToken{
	//		Token: key,
	//		Role:  models.RoleUser,
	//	}
	//	registerTokenTestString, err := json.Marshal(registerTokenTest)
	//	if err != nil {
	//		panic(err)
	//	}
	//	stringCmd.SetErr(nil)
	//	stringCmd.SetVal(string(registerTokenTestString))
	//} else {
	//	stringCmd.SetErr(redis.Nil)
	//}
	//
	//return stringCmd
}

func TestBeforeRegister(t *testing.T) {
	type TestData struct {
		testName      string
		expectedCode  int
		regToken      string
		adminRegister bool
	}

	testData := []TestData{
		{
			testName:      "SuccesBeforeRegister",
			expectedCode:  http.StatusOK,
			regToken:      "testtesttesttest",
			adminRegister: false,
		},
		{
			testName:      "SuccesBeforeRegisterAdmin",
			expectedCode:  http.StatusOK,
			regToken:      "testtesttesttest",
			adminRegister: true,
		},
		{
			testName:      "TokenFailedBeforeRegister",
			expectedCode:  http.StatusUnauthorized,
			regToken:      "testfailed",
			adminRegister: false,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			registerRequest, err := json.Marshal(users.RegisterRequest{
				RegisterToken: td.regToken,
			})
			require.NoError(t, err)

			c.Request.Body = ioutil.NopCloser(bytes.NewReader(registerRequest))

			roleMock, err := json.Marshal(models.RoleUser)
			require.NoError(t, err)
			// End of setup

			// Mocking
			mockInterface := &beforeRegisterMock{}

			switch td.testName {
			case "SuccesBeforeRegister":
				redisResponse := redis.StringCmd{}
				redisResponse.SetErr(nil)
				redisResponse.SetVal(string(roleMock))
				mockInterface.On("First", mock.Anything).Return(&gorm.DB{}).Run(func(args mock.Arguments) {
					args.Get(0).(*models.User).Username = "test"
				})
				mockInterface.On("Get", mock.Anything, mock.Anything).Return(&redisResponse)
			case "SuccesBeforeRegisterAdmin":
				mockInterface.On("First", mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			case "TokenFailedBeforeRegister":
				redisResponse := redis.StringCmd{}
				redisResponse.SetErr(redis.Nil)
				redisResponse.SetVal(string(roleMock))
				mockInterface.On("First", mock.Anything).Return(&gorm.DB{}).Run(func(args mock.Arguments) {
					args.Get(0).(*models.User).Username = "test"
				})
				mockInterface.On("Get", mock.Anything, mock.Anything).Return(&redisResponse)
			}
			// End of mocking

			// Run Test
			ctrl := MiddlewareController{
				DB:  mockInterface,
				RDB: mockInterface,
			}
			ctrl.BeforeRegister(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
