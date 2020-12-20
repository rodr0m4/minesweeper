package rest

import (
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteGameHandler struct {
	GameDeleter GameDeleter
}

type GameDeleter interface {
	DeleteGame(id game.ID) (operation.ShowedGame, error)
}

func (h DeleteGameHandler) DeleteGame(ctx *gin.Context) {
	if err := h.deleteGame(ctx); err != nil {
		_ = ctx.Error(err)
	}
}

func (h DeleteGameHandler) deleteGame(ctx *gin.Context) error {
	id, err := ExtractIDFromPath(ctx)

	if err != nil {
		return err
	}

	sg, err := h.GameDeleter.DeleteGame(id)

	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, sg)
	return nil
}
