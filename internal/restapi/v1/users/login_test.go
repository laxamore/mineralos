package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type loginMock struct {
	mock.Mock
	db.IDB
	jwt.Cmd
}

func (m loginMock) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func (m loginMock) SignJWT(claims jwt.Claims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

func TestLogin(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		loginRequest LoginRequest
	}

	testData := []TestData{
		{
			testName:     "SuccessLogin",
			expectedCode: http.StatusOK,
			loginRequest: LoginRequest{
				Username: "test",
				Password: "test1234",
			},
		},
		{
			testName:     "FailedLogin",
			expectedCode: http.StatusUnauthorized,
			loginRequest: LoginRequest{
				Username: "test",
				Password: "test1234",
			},
		}, {
			testName:     "NoPassword",
			expectedCode: http.StatusBadRequest,
			loginRequest: LoginRequest{
				Username: "test",
			},
		}, {
			testName:     "NoUserName",
			expectedCode: http.StatusBadRequest,
			loginRequest: LoginRequest{
				Password: "test1234",
			},
		},
	}

	config.Config = &config.ConfigStruct{
		ACCESS_TOKEN_EXPIRED:  3600,
		REFRESH_TOKEN_EXPIRED: 3600,
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Mocking
			mockInterface := &loginMock{}

			switch td.testName {
			case "SuccessLogin":
				mockInterface.On("First", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{Error: nil}).Run(func(args mock.Arguments) {
					user := args.Get(0).(*models.User)
					user.Username = "test"
				})
				mockInterface.On("SignJWT", mock.Anything).Return("token", nil)
			case "FailedLogin":
				mockInterface.On("First", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			}
			// End of mocking

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			jsonbytes, err := json.Marshal(td.loginRequest)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			// End of setup

			// Run Test
			ctrl := UserController{
				DB:         mockInterface,
				JWTService: mockInterface,
			}
			ctrl.Login(c)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
