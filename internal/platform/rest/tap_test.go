package rest

import (
	"bytes"
	"encoding/json"
	"minesweeper/internal/platform/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_Tap_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := TapHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerTap(r, handler)

	req := newTapRequestFromBytes([]byte("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Helpers

func newTapRequestFromBytes(buf []byte) *http.Request {
	return httptest.NewRequest(http.MethodPatch, "/games", bytes.NewBuffer(buf))
}

func newTapRequest(body gin.H) *http.Request {
	buf, _ := json.Marshal(body)
	return newTapRequestFromBytes(buf)
}

func registerTap(r *gin.Engine, handler TapHandler) {
	r.Use(middleware.ErrorLogger())
	r.PATCH("/games", handler.Tap)
}
