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

type deleteWalletMocking struct {
	mock.Mock
	db.IDB
}

func (m deleteWalletMocking) Delete(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func TestDeleteWallet(t *testing.T) {
	type TestData struct {
		testName            string
		expectedCode        int
		deleteWalletRequest DeleteWalletRequest
	}

	testData := []TestData{
		{
			testName:     "SuccessDeleteWallet",
			expectedCode: http.StatusOK,
			deleteWalletRequest: DeleteWalletRequest{
				WalletId: 1234,
			},
		},
		{
			testName:     "FailedDeleteRequest",
			expectedCode: http.StatusBadRequest,
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

			jsonbytes, err := json.Marshal(td.deleteWalletRequest)
			if err != nil {
				panic(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
			// End of Setup

			// Mocking
			mockInterface := &deleteWalletMocking{}
			mockInterface.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{})
			// End of Mocking

			// Run Test
			ctrl := RigController{
				DB: mockInterface,
			}
			ctrl.DeleteWallet(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
