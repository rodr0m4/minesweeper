package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := StartGameHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	r.POST("/games", handler.StartGame)

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
