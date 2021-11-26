package ApiUsers

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api/api"
)

func Logout(c *gin.Context) {
	response := api.Result{
		Code:     http.StatusOK,
		Response: "logout success",
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "rtoken",
		Value:    url.QueryEscape("rtoken"),
		MaxAge:   -1,
		Path:     "/",
		Domain:   os.Getenv("domain"),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(response.Code, response.Response)
}
