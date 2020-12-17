package rest

import (
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShowGameHandler struct {
	Game       game.Game
	GameShower GameShower
}

func (h ShowGameHandler) ShowGame(ctx *gin.Context) {
	sg, err := h.GameShower.ShowGame(h.Game)

	if err != nil {
		abortWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, sg)
}
