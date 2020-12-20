package rest

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModifyTileHandler struct {
	Game        game.Game
	BoardDrawer operation.BoardDrawer
	Tapper      Tapper
	Marker      Marker
}

type Tapper interface {
	Tap(game game.Game, row, column int) (internal.TapResult, error)
}

type Marker interface {
	Mark(game game.Game, row, column int) error
}

type modifyTileRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type modifyTileResponse struct {
	Result string   `json:"result"`
	Lines  []string `json:"lines"`
}

func (h ModifyTileHandler) Mark(ctx *gin.Context) {
	request, err := h.bindRequest(ctx)

	if err != nil {
		return
	}

	if err := h.Marker.Mark(h.Game, request.Row, request.Column); err != nil {
		_ = ctx.Error(err)
		return
	}

	h.renderBoard(ctx, "")
}

func (h ModifyTileHandler) Tap(ctx *gin.Context) {
	request, err := h.bindRequest(ctx)

	if err != nil {
		return
	}

	result, err := h.Tapper.Tap(h.Game, request.Row, request.Column)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	h.renderBoard(ctx, tapResultToString(result))
}

func (h ModifyTileHandler) renderBoard(ctx *gin.Context, result string) {
	showed, err := operation.DrawGame(h.Game, h.BoardDrawer)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var response modifyTileResponse

	response.Result = result
	response.Lines = showed.Lines

	ctx.JSON(http.StatusOK, response)
}

func (h ModifyTileHandler) bindRequest(ctx *gin.Context) (modifyTileRequest, error) {
	var request modifyTileRequest
	err := ctx.BindJSON(&request)

	return request, err
}

func tapResultToString(result internal.TapResult) string {
	switch result {
	case internal.NothingResult:
		return "Nothing"
	case internal.ExplosionResult:
		return "BOOM!"
	}

	panic(fmt.Sprintf("unrecheable code, invalid tap result %d", result))
}
