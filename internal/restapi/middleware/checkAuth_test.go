package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/databases/models/user"
	JWT "github.com/laxamore/mineralos/internal/jwt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckAuth(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		roleToCheck  *user.Role
		actualRole   *user.Role
	}

	testData := []TestData{
		{
			testName:     "OKCheckAuthRole",
			expectedCode: http.StatusOK,
			roleToCheck:  &user.RoleAdmin,
			actualRole:   &user.RoleAdmin,
		},
		{
			testName:     "UnauthorizedCheckAuthRole",
			expectedCode: http.StatusUnauthorized,
			roleToCheck:  &user.RoleAdmin,
			actualRole:   &user.RoleUser,
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

			config.Config = &config.ConfigStruct{
				JWT_SECRET:           "123",
				ACCESS_TOKEN_EXPIRED: 3600,
			}

			token, err := JWT.SignJWT(JWT.JWTCustomClaims{
				Username: "test",
				Email:    "test@test.com",
				Role:     td.actualRole,
			}, time.Now().Unix()+config.Config.ACCESS_TOKEN_EXPIRED)
			require.NoError(t, err)

			// Set Authorization header
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			ctrl := MiddlewareController{}
			ctrl.CheckAuthRole(c, td.roleToCheck)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
