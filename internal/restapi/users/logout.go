package users

import (
	"github.com/laxamore/mineralos/internal/restapi"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	response := restapi.Result{
		Code:     http.StatusOK,
		Response: "logout success",
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "rtoken",
		Value:    url.QueryEscape("rtoken"),
		MaxAge:   -1,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(response.Code, response.Response)
}
