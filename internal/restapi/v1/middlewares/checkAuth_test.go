package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type checkAuthMock struct {
	mock.Mock
	jwt.Cmd
}

func (m checkAuthMock) VerifyJWT(token string, claims jwt.Claims) (*jwt.Token, error) {
	args := m.Called(token, claims)
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func TestCheckAuth(t *testing.T) {
	type TestData struct {
		testName     string
		expectedCode int
		roleToCheck  *models.Role
		actualRole   *models.Role
	}

	testData := []TestData{
		{
			testName:     "OKCheckAuthRole",
			expectedCode: http.StatusOK,
			roleToCheck:  &models.RoleAdmin,
			actualRole:   &models.RoleAdmin,
		},
		{
			testName:     "UnauthorizedCheckAuthRole",
			expectedCode: http.StatusUnauthorized,
			roleToCheck:  &models.RoleAdmin,
			actualRole:   &models.RoleUser,
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			// Mocking
			mockInterface := &checkAuthMock{}
			switch td.testName {
			case "OKCheckAuthRole":
				returnToken := jwt.Token{}
				returnToken.Valid = true
				mockInterface.On("VerifyJWT", mock.Anything, mock.Anything).Return(&returnToken, nil).Run(func(args mock.Arguments) {
					claims := args.Get(1).(*jwt.LoginClaims)
					claims.Role = td.actualRole
				})
			case "UnauthorizedCheckAuthRole":
				returnToken := jwt.Token{}
				returnToken.Valid = false
				mockInterface.On("VerifyJWT", mock.Anything, mock.Anything).Return(&returnToken, nil)
			}
			// End of mocking

			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Header: make(http.Header),
			}
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "tokentest"))
			// End of setup

			// Run Test
			ctrl := MiddlewareController{
				JWTService: mockInterface,
			}
			ctrl.CheckAuthRole(c, td.roleToCheck)

			require.EqualValues(t, fmt.Sprintf("HTTP Status Code: %d", td.expectedCode), fmt.Sprintf("HTTP Status Code: %d", w.Code))
		})
	}
}
