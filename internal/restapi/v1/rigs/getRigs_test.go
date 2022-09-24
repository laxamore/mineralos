package rigs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type getRigsMock struct {
	mock.Mock
	db.IDB
}

func (m getRigsMock) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func TestGetRigs(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
	}

	testData := []TestData{
		{
			testName:     "SuccessGetRigs",
			expectedCode: http.StatusOK,
		},
		{
			testName:     "NoRecordGetRigs",
			expectedCode: http.StatusNoContent,
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
			// End of Setup

			// Mocking
			mockInterface := &getRigsMock{}
			switch td.testName {
			case "SuccessGetRigs":
				mockInterface.On("Find", mock.Anything, mock.Anything).Return(&gorm.DB{}).Run(func(args mock.Arguments) {
					*args.Get(0).(*[]models.Rig) = []models.Rig{
						{
							ID: 1,
						},
					}
				})
			case "NoRecordGetRigs":
				mockInterface.On("Find", mock.Anything, mock.Anything).Return(&gorm.DB{})
			}
			// End of Mocking

			// Run Test
			ctrl := RigController{
				DB: mockInterface,
			}
			ctrl.GetRigs(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
