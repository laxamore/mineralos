package users

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
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
	db.IRedis
}

func (m *loginMock) First(out interface{}, where ...interface{}) *gorm.DB {
	correctPassword := "test1234"
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s", correctPassword)))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	if where[0].(string) == "username = ? AND password = ?" && where[1].(string) == "test" && where[2].(string) == hashedPassword {
		*out.(*models.User) = models.User{
			Username: "test",
		}
		return &gorm.DB{Error: nil}
	}
	return &gorm.DB{Error: gorm.ErrRecordNotFound}
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
				Password: "testfailed",
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

			ctrl := UserController{
				DB: &loginMock{},
			}
			ctrl.Login(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
