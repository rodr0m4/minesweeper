package rest

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TapHandler struct {
	Game        game.Game
	Tapper      Tapper
	BoardDrawer operation.BoardDrawer
}

type tapHandlerRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type tapHandlerResponse struct {
	Result string               `json:"result"`
	Game   operation.ShowedGame `json:"game"`
}

type Tapper interface {
	Tap(game game.Game, row, column int) (internal.TapResult, error)
}

func (h TapHandler) Tap(ctx *gin.Context) {
	var request tapHandlerRequest

	if err := ctx.BindJSON(&request); err != nil {
		return
	}

	result, err := h.Tapper.Tap(h.Game, request.Row, request.Column)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	showed, err := operation.DrawGame(h.Game, h.BoardDrawer)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var response tapHandlerResponse

	response.Result = tapResultToString(result)
	response.Game = showed

	ctx.JSON(http.StatusOK, response)
}

func tapResultToString(result internal.TapResult) string {
	switch result {
	case internal.NothingResult:
		return "Nothing"
	case internal.ExplosionResult:
		return "BOOM!"
	}

	panic(fmt.Sprintf("unrecheable code, invalid result %d", result))
}
