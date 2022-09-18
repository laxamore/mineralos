package rigs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type dbMock struct {
	mock.Mock
}

func (a dbMock) Save(value interface{}) (tx *gorm.DB) {
	return &gorm.DB{Error: nil}
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

			db := dbMock{}
			cntrl := NewRigController{}

			cntrl.TryNewRig(c, db)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}