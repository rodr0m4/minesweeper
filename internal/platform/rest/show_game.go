package rest

import (
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShowGameHandler struct {
	Game        game.Game
	BoardDrawer operation.BoardDrawer
}

func (h ShowGameHandler) ShowGame(ctx *gin.Context) {
	sg, err := operation.DrawGame(h.Game, h.BoardDrawer)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, sg)
}
