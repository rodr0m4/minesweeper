package rest

import (
	"errors"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteGame_Should_Fail_With_BadRequest_If_Not_Valid_ID(t *testing.T) {
	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	var handler DeleteGameHandler

	registerDeleteGame(r, handler)

	r.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/games/asd", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func Test_DeleteGame_Should_Call_Deleter_And_Fail_If_Deleter_Fails(t *testing.T) {
	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	handler := DeleteGameHandler{
		GameDeleter: gameDeleterMock(func(id game.ID) (operation.ShowedGame, error) {
			return operation.ShowedGame{}, errors.New("oh no")
		}),
	}

	registerDeleteGame(r, handler)

	r.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/games/1", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_DeleteGame_Should_Call_Deleter(t *testing.T) {
	rr := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rr)

	handler := DeleteGameHandler{
		GameDeleter: gameDeleterMock(func(id game.ID) (operation.ShowedGame, error) {
			return operation.ShowedGame{}, nil
		}),
	}

	registerDeleteGame(r, handler)

	r.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/games/1", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

// Helpers

func registerDeleteGame(r *gin.Engine, handler DeleteGameHandler) {
	r.Use(middleware.ErrorLogger())
	r.DELETE("/games/:id", handler.DeleteGame)
}

// Mocks

type gameDeleterMock func(game.ID) (operation.ShowedGame, error)

func (g gameDeleterMock) DeleteGame(id game.ID) (operation.ShowedGame, error) {
	return g(id)
}
