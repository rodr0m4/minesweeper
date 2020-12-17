package rest

import (
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type startGameRequest struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`
	Bombs   int `json:"bombs"`
}

type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GameStarter interface {
	StartGame(game game.Game, rows, columns, bombs int) error
}

type GameShower interface {
	ShowGame(game game.Game) (operation.ShowedGame, error)
}

type StartGameHandler struct {
	Game        game.Game
	GameStarter GameStarter
	GameShower  GameShower
}

func (s StartGameHandler) StartGame(ctx *gin.Context) {
	var request startGameRequest

	if err := ctx.BindJSON(&request); err != nil {
		return
	}

	if err := s.GameStarter.StartGame(s.Game, request.Rows, request.Columns, request.Bombs); err != nil {
		abortWithError(ctx, err)
		return
	}

	sg, err := s.GameShower.ShowGame(s.Game)

	if err != nil {
		abortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, sg)
}

func errorToStatusCode(err error) int {
	if errors.Is(err, internal.ErrInvalidOperation) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

// TODO: Refactor these into a middleware
func abortWithError(ctx *gin.Context, err error) {
	abortWithErrorAndStatus(ctx, errorToStatusCode(err), err)
}

func abortWithErrorAndStatus(ctx *gin.Context, code int, err error) {
	var json JSONError

	json.Code = code
	json.Message = err.Error()

	ctx.AbortWithStatusJSON(code, json)
}
