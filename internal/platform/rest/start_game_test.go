package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_StartGame_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := StartGameHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerStartGame(r, handler)

	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString("<title/>"))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_Should_Fail_When_Game_Cant_Start(t *testing.T) {
	type Case struct {
		starter      GameStarter
		expectedCode int
	}

	cases := []Case{
		{
			starter:      gameStarterThatFailsWith(internal.NewInvalidOperation("oh no")),
			expectedCode: http.StatusBadRequest,
		},
		{
			starter:      gameStarterThatFailsWith(errors.New("oh no")),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		rr := httptest.NewRecorder()
		_, r := gin.CreateTestContext(rr)

		handler := StartGameHandler{
			GameStarter: tt.starter,
		}

		registerStartGame(r, handler)

		req := newStartGameRequest()

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func Test_Should_Fail_When_ShowGame_Fails(t *testing.T) {
	g := game.Fake{}
	err := errors.New("oh no")

	starter := gameStarterMock(func(_ game.Game, _, _, _ int) error {
		return nil
	})
	shower := gameShowerMock(func(actual game.Game) (operation.ShowedGame, error) {
		assert.Equal(t, g, actual)

		return operation.ShowedGame{}, err
	})

	handler := StartGameHandler{
		Game:        g,
		GameStarter: starter,
		GameShower:  shower,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerStartGame(r, handler)

	r.ServeHTTP(rr, newStartGameRequest())

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_Should_Return_Showed_Game_And_Created_When_Passes(t *testing.T) {
	g := game.Fake{}
	lines := []string{"Hello", "World"}

	starter := gameStarterMock(func(_ game.Game, _, _, _ int) error {
		return nil
	})
	shower := gameShowerMock(func(actual game.Game) (operation.ShowedGame, error) {
		assert.Equal(t, g, actual)

		return operation.ShowedGame{
			Lines: lines,
		}, nil
	})

	handler := StartGameHandler{
		Game:        g,
		GameStarter: starter,
		GameShower:  shower,
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerStartGame(r, handler)

	r.ServeHTTP(rr, newStartGameRequest())

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response operation.ShowedGame
	_ = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, lines, response.Lines)
}

// Helpers

func newStartGameRequest() *http.Request {
	body, _ := json.Marshal(gin.H{
		"rows":    2,
		"columns": 2,
		"bombs":   1,
	})

	buf := bytes.NewBuffer(body)
	req := httptest.NewRequest(http.MethodPost, "/games", buf)

	return req
}

func registerStartGame(r *gin.Engine, handler StartGameHandler) {
	r.Use(middleware.ErrorLogger())
	r.POST("/games", handler.StartGame)
}

// Mocks

type gameShowerMock func(game.Game) (operation.ShowedGame, error)

func (g gameShowerMock) ShowGame(game game.Game) (operation.ShowedGame, error) {
	return g(game)
}

type gameStarterMock func(game game.Game, rows, columns, bombs int) error

func (g gameStarterMock) StartGame(game game.Game, rows, columns, bombs int) error {
	return g(game, rows, columns, bombs)
}

func gameStarterThatFailsWith(err error) gameStarterMock {
	return func(game game.Game, rows, columns, bombs int) error {
		return err
	}
}
