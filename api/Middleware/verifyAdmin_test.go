package Middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/utils/JWT"
	"github.com/stretchr/testify/require"
)

func TestVerifyAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type TestData struct {
		testName      string
		userPrivilege string
		expectedCode  int
	}

	testData := []TestData{
		{
			testName:      "SuccessVerifyAdmin",
			userPrivilege: "admin",
			expectedCode:  http.StatusOK,
		},
		{
			testName:      "UnprivilegeUser",
			userPrivilege: "readOnly",
			expectedCode:  http.StatusUnauthorized,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}

			// Create test jwt
			claims := jwt.MapClaims{
				"username":  "test",
				"email":     "test@test.com",
				"privilege": td.userPrivilege,
			}
			token, err := JWT.SignJWT(claims, time.Now().Unix()+5)
			if err != nil {
				t.Fatalf("Login Token Sign Failed:\n%v", err)
			}
			//

			c.Request.Header.Set("token", token)

			cntrl := VerifyAdminController{}
			cntrl.TryVerifyAdmin(c)
			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
