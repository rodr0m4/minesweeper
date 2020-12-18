package rest

import (
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
		_ = ctx.Error(err)
		return
	}

	sg, err := s.GameShower.ShowGame(s.Game)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, sg)
}
