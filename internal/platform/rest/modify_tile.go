package rest

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ModifyTileHandler struct {
	GameHolder  game.Holder
	BoardDrawer operation.BoardDrawer
	Tapper      Tapper
	Marker      Marker
}

type Tapper interface {
	Tap(game game.Game, row, column int) (internal.TapResult, error)
}

type Marker interface {
	Mark(game game.Game, row, column int, mark internal.TileMark) error
}

type modifyTileRequest struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Mark   string `json:"mark"`
}

type modifyTileResponse struct {
	Result string               `json:"result"`
	Game   operation.ShowedGame `json:"game"`
}

func (h ModifyTileHandler) Mark(ctx *gin.Context) {
	request, err := h.bindRequest(ctx)

	if err != nil {
		return
	}

	g, err := h.getGame(ctx)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	mark, err := stringToTileMark(request.Mark)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := h.Marker.Mark(g, request.Row, request.Column, mark); err != nil {
		_ = ctx.Error(err)
		return
	}

	h.renderBoard(ctx, g, "")
}

func (h ModifyTileHandler) Tap(ctx *gin.Context) {
	request, err := h.bindRequest(ctx)

	if err != nil {
		return
	}

	g, err := h.getGame(ctx)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	result, err := h.Tapper.Tap(g, request.Row, request.Column)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	h.renderBoard(ctx, g, tapResultToString(result))
}

func (h ModifyTileHandler) renderBoard(ctx *gin.Context, game game.Game, result string) {
	showed, err := operation.DrawGame(game, h.BoardDrawer)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var response modifyTileResponse

	response.Result = result
	response.Game = showed

	ctx.JSON(http.StatusOK, response)
}

func (h ModifyTileHandler) getGame(ctx *gin.Context) (game.Game, error) {
	id, err := ExtractIDFromPath(ctx)

	if err != nil {
		return nil, err
	}

	return h.GameHolder.Get(id)
}

func (h ModifyTileHandler) bindRequest(ctx *gin.Context) (modifyTileRequest, error) {
	var request modifyTileRequest
	err := ctx.BindJSON(&request)

	return request, err
}

func tapResultToString(result internal.TapResult) string {
	switch result {
	case internal.NothingResult:
		return "Nothing Happened"
	case internal.ExplosionResult:
		return "BOOM!"
	}

	panic(fmt.Sprintf("unrecheable code, invalid tap result %d", result))
}

func stringToTileMark(s string) (internal.TileMark, error) {
	s = strings.ToLower(s)

	switch s {
	case "flag":
		return internal.FlagMark, nil
	case "question":
		return internal.QuestionMark, nil
	}

	return 0, internal.NewInvalidOperation(fmt.Sprintf("invalid tile mark '%s'", s))
}
