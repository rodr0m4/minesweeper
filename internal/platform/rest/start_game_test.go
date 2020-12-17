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

func Test_Should_Fail_When_Passed_Invalid_JSON(t *testing.T) {
	handler := StartGameHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	r.POST("/games", handler.StartGame)

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
		r.POST("/games", handler.StartGame)

		body, _ := json.Marshal(gin.H{
			"rows":    2,
			"columns": 2,
			"bombs":   1,
		})
		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func Test_Should_Fail_When_Game_Board_Fails(t *testing.T) {
	type Case struct {
		game         game.Game
		expectedCode int
	}

	cases := []Case{
		{
			game:         gameWhoseBoardFailsWith(internal.NewInvalidOperation("oh no")),
			expectedCode: http.StatusBadRequest,
		},
		{
			game:         gameWhoseBoardFailsWith(errors.New("oh no")),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		rr := httptest.NewRecorder()
		_, r := gin.CreateTestContext(rr)

		handler := StartGameHandler{
			Game:        tt.game,
			GameStarter: gameStarterThatFailsWith(nil),
		}
		r.POST("/games", handler.StartGame)

		body, _ := json.Marshal(gin.H{
			"rows":    2,
			"columns": 2,
			"bombs":   1,
		})
		req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func Test_Should_Draw_Board_And_Return_Created(t *testing.T) {
	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	board := internal.NewBoard(2, 2, 1)
	g := game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return board, nil
		},
	}

	handler := StartGameHandler{
		Game:        g,
		GameStarter: gameStarterThatFailsWith(nil),
		BoardDrawer: boardDrawerMock(func(actual internal.Board, revealEverything bool) []string {
			assert.Equal(t, board, actual)
			return []string{"Hello", "World"}
		}),
	}

	r.POST("/games", handler.StartGame)

	body, _ := json.Marshal(gin.H{
		"rows":    2,
		"columns": 2,
		"bombs":   1,
	})
	req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response startGameResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, []string{"Hello", "World"}, response.Lines)
}

// Mocks

type boardDrawerMock func(board internal.Board, revealEverything bool) []string

func (b boardDrawerMock) DrawBoardIntoStringArray(board internal.Board, revealEverything bool) []string {
	return b(board, revealEverything)
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

func gameWhoseBoardFailsWith(err error) game.Game {
	return game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return internal.Board{}, err
		},
	}
}
