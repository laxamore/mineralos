package rigs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type newWalletMocking struct {
	mock.Mock
	db.IDB
}

func (m newWalletMocking) Create(value interface{}) (tx *gorm.DB) {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestNewWallet(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		bodyData     gin.H
	}

	testData := []TestData{
		{
			testName:     "SuccessNewRIG",
			expectedCode: http.StatusOK,
			bodyData: gin.H{
				"wallet_coin":    "eth",
				"wallet_name":    "test",
				"wallet_address": "0x9062bd26e6a84086634bb4322cd526f368e4e402",
			},
		},
		{
			testName:     "FailedRequest",
			expectedCode: http.StatusBadRequest,
			bodyData: gin.H{
				"wallet_coin": "eth",
				"wallet_name": "test",
			},
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

			jsonbytes, err := json.Marshal(td.bodyData)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			// End of Setup

			// Mocking
			mockingInterface := &newWalletMocking{}
			mockingInterface.On("Create", mock.Anything).Return(&gorm.DB{})
			// End of Mocking

			// Run Test
			ctrl := RigController{
				DB: mockingInterface,
			}
			ctrl.NewWallet(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
