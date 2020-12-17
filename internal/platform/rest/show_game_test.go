package rest

import (
	"encoding/json"
	"errors"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_ShowGame_Should_Fail_If_Shower_Fails(t *testing.T) {
	g := game.Fake{}
	err := errors.New("oh no")

	shower := gameShowerMock(func(actual game.Game) (operation.ShowedGame, error) {
		assert.Equal(t, g, actual)
		return operation.ShowedGame{}, err
	})

	handler := ShowGameHandler{
		Game:       g,
		GameShower: shower,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	r.GET("/games", handler.ShowGame)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_ShowGame_Should_Return_Showed_Game_When_Does_Not_Fail(t *testing.T) {
	g := game.Fake{}
	lines := []string{"Hello", "World"}

	shower := gameShowerMock(func(actual game.Game) (operation.ShowedGame, error) {
		assert.Equal(t, g, actual)
		return operation.ShowedGame{
			Lines: lines,
		}, nil
	})

	handler := ShowGameHandler{
		Game:       g,
		GameShower: shower,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	r.GET("/games", handler.ShowGame)

	req := httptest.NewRequest(http.MethodGet, "/games", nil)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var sg operation.ShowedGame
	_ = json.Unmarshal(rr.Body.Bytes(), &sg)

	assert.Equal(t, lines, sg.Lines)
}
