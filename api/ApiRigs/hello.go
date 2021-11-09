package ApiRigs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/api"
)

func Hello(c *gin.Context) {
	response := api.Result{
		Code:     http.StatusBadRequest,
		Response: "Bad Request",
	}

	response.Code = http.StatusOK
	response.Response = gin.H{
		// "Auth":   strings.Split(c.GetHeader("Authorization"), " ")[1],
		"Auth": c.GetHeader("Authorization"),
		"msg":  "Hello World!",
	}

	c.JSON(response.Code, response.Response)
}
