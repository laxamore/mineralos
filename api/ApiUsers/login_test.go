package ApiUsers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginRepositoryMock struct {
	mock.Mock
}

func (r LoginRepositoryMock) FindOne(db_name string, collection_name string, filter interface{}) map[string]interface{} {
	users := map[string]interface{}{
		"username":  "laxamore",
		"email":     "laxamore@gmail.com",
		"password":  "eb676aabfb58c2fc243b570214b658f51353002cd647190eb700f877bdd3ed0b",
		"privilege": "admin",
	}

	var input map[string]interface{}

	filterBytes, _ := bson.Marshal(filter)
	bson.Unmarshal(filterBytes, &input)

	if (input["username"] == users["username"] && input["password"] == users["password"]) ||
		(input["email"] == users["email"] && input["password"] == users["password"]) {
		return users
	}

	return nil
}

func TestLoginUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	jsonbytes, err := json.Marshal(gin.H{
		"username": "laxamore",
		"password": "39393939",
	})
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

	repo := LoginRepositoryMock{}
	ctrl := LoginController{}

	ctrl.TryLogin(c, repo)
	assert.EqualValues(t, http.StatusOK, w.Code)
}
