package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
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

	r.PATCH("/games", handler.Tap)

	req := httptest.NewRequest(http.MethodPatch, "/games", bytes.NewBufferString("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_Tap_Fails_If_Game_Board_Fails(t *testing.T) {
	g := game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return internal.Board{}, errors.New("oh no")
		},
	}

	handler := TapHandler{
		Game: g,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	r.PATCH("/games", handler.Tap)

	body := gin.H{
		"row":    5,
		"column": 5,
	}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/games", bytes.NewBuffer(buf))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Tap_Fails_If_Invalid_Position(t *testing.T) {
	board := internal.NewBoard(2, 2, 1)
	g := game.NewInMemory(&board)

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	handler := TapHandler{Game: g}
	r.PATCH("/games", handler.Tap)

	body := gin.H{
		"row":    5,
		"column": 5,
	}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, "/games", bytes.NewBuffer(buf))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
