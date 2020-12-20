package rest

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ExtractIDFromPath(ctx *gin.Context) (game.ID, error) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		return 0, internal.NewInvalidOperation(err.Error())
	}

	return game.ID(id), nil
}
