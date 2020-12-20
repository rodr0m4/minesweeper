package rest

import (
	"minesweeper/internal/platform/game"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine            *gin.Engine
	GameHolder        game.Holder
	CreateGameHandler *CreateGameHandler
	ShowGameHandler   *ShowGameHandler
	ModifyTileHandler *ModifyTileHandler
}
