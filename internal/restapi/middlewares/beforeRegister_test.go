package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type dbMock struct {
	mock.Mock
	db.IDB
	db.IRedis
}

var isAdminRegister = false

func (m dbMock) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	if isAdminRegister {
		return &gorm.DB{}
	} else {
		dest.(*models.User).Username = "test"
		return &gorm.DB{}
	}
}

func (m dbMock) Get(ctx context.Context, key string) *redis.StringCmd {
	stringCmd := &redis.StringCmd{}

	if key == "testtesttesttest" {
		registerTokenTest := models.RegisterToken{
			Token: key,
			Role:  models.RoleUser,
		}
		registerTokenTestString, err := json.Marshal(registerTokenTest)
		if err != nil {
			panic(err)
		}
		stringCmd.SetErr(nil)
		stringCmd.SetVal(string(registerTokenTestString))
	} else {
		stringCmd.SetErr(redis.Nil)
	}

	return stringCmd
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
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Request.Header.Add("regToken", td.regToken)

			isAdminRegister = td.adminRegister
			ctrl := MiddlewareController{
				DB:  dbMock{},
				RDB: dbMock{},
			}
			ctrl.BeforeRegister(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
