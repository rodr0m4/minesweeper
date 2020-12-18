package rest

import (
	"minesweeper/internal/platform/game"

	"github.com/gin-gonic/gin"
)

type TapHandler struct {
	Game game.Game
}

type tapHandlerRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

func (h TapHandler) Tap(ctx *gin.Context) {
	var request tapHandlerRequest

	if err := ctx.BindJSON(&request); err != nil {
		return
	}
}
