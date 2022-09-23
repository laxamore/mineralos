package rigs

import (
	"fmt"
	"github.com/laxamore/mineralos/internal/db"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type getRigMock struct {
	mock.Mock
	db.IDB
}

func (m *getRigMock) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetRig(t *testing.T) {
	type TestData struct {
		testName     string
		rig_id       string
		expectedCode int
	}

	testData := []TestData{
		{
			testName:     "SuccessGetRig",
			expectedCode: http.StatusOK,
		},
		{
			testName:     "FailGetRig",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Mocking
			mockInterface := &getRigMock{}

			switch td.testName {
			case "SuccessGetRig":
				mockInterface.On("First", mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{Error: nil})
			case "FailGetRig":
				mockInterface.On("First", mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			}
			// End of mocking

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = []gin.Param{
				{
					Key:   "rig_id",
					Value: td.rig_id,
				},
			}
			// End of setup

			// Run Test
			cntrl := RigController{
				DB: mockInterface,
			}
			cntrl.GetRig(c)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
