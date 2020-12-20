package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_CreateGame_Should_Fail_When_Given_Invalid_JSON(t *testing.T) {
	handler := CreateGameHandler{}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerCreateGame(r, handler)

	r.ServeHTTP(rr, newCreateGameRequestFromBytes([]byte("<title />")))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_CreateGame_Should_Fail_If_GameCreator_Fails(t *testing.T) {
	rows := 2
	columns := 2
	bombs := 1

	handler := CreateGameHandler{
		GameCreator: gameCreatorThatExpectsAndReturns(t, rows, columns, bombs, func() (int, internal.Board, error) {
			return 0, internal.Board{}, errors.New("oh no")
		}),
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerCreateGame(r, handler)

	r.ServeHTTP(rr, newCreateGameRequest(rows, columns, bombs))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_CreateGame_Should_Call_DrawBoard_If_Can_CreateGame(t *testing.T) {
	rows := 2
	columns := 2
	bombs := 1

	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{})
	lines := []string{"hello", "world"}

	handler := CreateGameHandler{
		GameCreator: gameCreatorThatExpectsAndReturns(t, rows, columns, bombs, func() (int, internal.Board, error) {
			return 0, board, nil
		}),
		BoardDrawer: boardDrawerMock(func(actual internal.Board) operation.ShowedGame {
			assert.Equal(t, board, actual)
			return operation.ShowedGame{Lines: lines}
		}),
	}

	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	registerCreateGame(r, handler)

	r.ServeHTTP(rr, newCreateGameRequest(rows, columns, bombs))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

// Helpers

func registerCreateGame(r *gin.Engine, handler CreateGameHandler) {
	r.Use(middleware.ErrorLogger())
	r.POST("/games", handler.CreateGame)
}

func newCreateGameRequest(rows, columns, bombs int) *http.Request {
	body := gin.H{
		"rows":    rows,
		"columns": columns,
		"bombs":   bombs,
	}

	buf, _ := json.Marshal(body)
	return newCreateGameRequestFromBytes(buf)
}

func newCreateGameRequestFromBytes(buf []byte) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(buf))
}

// Mocks

type gameCreatorMock func(rows, columns, bombs int) (int, internal.Board, error)

func (g gameCreatorMock) CreateGame(rows, columns, bombs int) (int, internal.Board, error) {
	return g(rows, columns, bombs)
}

func gameCreatorThatExpectsAndReturns(t *testing.T, rows, columns, bombs int, f func() (int, internal.Board, error)) gameCreatorMock {
	return func(actualRows, actualColumns, actualBombs int) (int, internal.Board, error) {
		assert.Equal(t, rows, actualRows)
		assert.Equal(t, columns, actualColumns)
		assert.Equal(t, bombs, actualBombs)

		return f()
	}
}
