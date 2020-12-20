package rest

import (
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateGameHandler struct {
	GameCreator GameCreator
	BoardDrawer operation.BoardDrawer
}

type GameCreator interface {
	CreateGame(rows, columns, bombs int) (game.ID, internal.Board, error)
}

type createGameRequest struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`
	Bombs   int `json:"bombs"`
}

type createGameResponse struct {
	ID   int                  `json:"id"`
	Game operation.ShowedGame `json:"game"`
}

func (h CreateGameHandler) CreateGame(ctx *gin.Context) {
	var request createGameRequest
	if err := ctx.BindJSON(&request); err != nil {
		return
	}

	id, board, err := h.GameCreator.CreateGame(request.Rows, request.Columns, request.Bombs)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := createGameResponse{
		ID:   int(id),
		Game: h.BoardDrawer.DrawBoard(board),
	}

	ctx.JSON(http.StatusCreated, response)
}
