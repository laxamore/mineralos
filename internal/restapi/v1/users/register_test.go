package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

type registerMock struct {
	mock.Mock
	db.IDB
	db.IRedis
}

func (m registerMock) Create(value interface{}) (tx *gorm.DB) {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestRegister(t *testing.T) {
	type TestData struct {
		testName        string
		expectedCode    int
		registerRequest RegisterRequest
	}

	testData := []TestData{
		{
			testName:     "SuccessRegister",
			expectedCode: http.StatusOK,
			registerRequest: RegisterRequest{
				Username: "test",
				Password: "test1234",
				Email:    "test@test.com",
			},
		},
		{
			testName:     "DuplicateUser",
			expectedCode: http.StatusConflict,
			registerRequest: RegisterRequest{
				Username: "test",
				Email:    "test@test.com",
				Password: "test1234",
			},
		},
		{
			testName:     "InvalidEmail",
			expectedCode: http.StatusBadRequest,
			registerRequest: RegisterRequest{
				Username: "test",
				Email:    "testemail.com",
				Password: "test1234",
			},
		},
		{
			testName:     "InvalidUsername",
			expectedCode: http.StatusBadRequest,
			registerRequest: RegisterRequest{
				Username: "!@test",
				Email:    "testemail.com",
				Password: "test1234",
			},
		},
		{
			testName:     "InvalidPassword",
			expectedCode: http.StatusBadRequest,
			registerRequest: RegisterRequest{
				Username: "test",
				Email:    "test@email.com",
				Password: "test123",
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Mocking
			mockInterface := &registerMock{}
			switch td.testName {
			case "SuccessRegister":
				mockInterface.On("Create", mock.Anything).Return(&gorm.DB{Error: nil})
			case "DuplicateUser":
				mockInterface.On("Create", mock.Anything).Return(&gorm.DB{Error: &mysql.MySQLError{Number: 1062}})
			}
			// End of mocking

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Set("role", models.RoleUser)

			jsonbytes, err := json.Marshal(td.registerRequest)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			// End of setup

			// Run Test
			ctrl := UserController{
				DB: mockInterface,
			}
			ctrl.Register(c)

			t.Logf("Response Body: %s", w.Body.String())
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
