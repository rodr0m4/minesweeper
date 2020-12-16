package rest

import (
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type startGameRequest struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`
	Bombs   int `json:"bombs"`
}

type startGameResponse struct {
	Lines []string `json:"lines"`
}

type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GameStarter interface {
	StartGame(game game.Game, rows, columns, bombs int) error
}

type StartGameHandler struct {
	Game        game.Game
	GameStarter GameStarter
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

	var response startGameResponse

	board, err := s.Game.Board()

	if err != nil {
		abortWithError(ctx, err)
		return
	}

	response.Lines = internal.DrawBoardIntoStringArray(board, true)

	ctx.JSON(http.StatusCreated, response)
}

func errorToStatusCode(err error) int {
	if errors.Is(err, internal.ErrInvalidOperation) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func abortWithError(ctx *gin.Context, err error) {
	abortWithErrorAndStatus(ctx, errorToStatusCode(err), err)
}

func abortWithErrorAndStatus(ctx *gin.Context, code int, err error) {
	var json JSONError

	json.Code = code
	json.Message = err.Error()

	ctx.AbortWithStatusJSON(code, json)
}
