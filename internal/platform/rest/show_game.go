package rest

import (
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShowGameHandler struct {
	GameHolder  game.Holder
	BoardDrawer operation.BoardDrawer
}

func (h ShowGameHandler) ShowGame(ctx *gin.Context) {
	id, err := ExtractIDFromPath(ctx)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	g, err := h.GameHolder.Get(id)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	sg, err := operation.DrawGame(g, h.BoardDrawer)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, sg)
}
