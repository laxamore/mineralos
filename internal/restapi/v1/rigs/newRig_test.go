package rigs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/laxamore/mineralos/internal/db"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type newRigMock struct {
	mock.Mock
	db.IDB
}

func (m newRigMock) Create(value interface{}) (tx *gorm.DB) {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestNewRig(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		privilege    string
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "OKNewRIG",
			expectedCode: http.StatusOK,
			bodyData: gin.H{
				"rig_name": "test",
			},
		},
		{
			testName:     "BadRequestNewRig",
			expectedCode: http.StatusBadRequest,
			bodyData: gin.H{
				"name": "test",
			},
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Mocking
			mockInterface := &newRigMock{}
			mockInterface.On("Create", mock.Anything).Return(&gorm.DB{})
			// End of mocking

			// Setup
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
			// End of setup

			// Run Test
			ctrl := RigController{
				DB: mockInterface,
			}
			ctrl.NewRig(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
