package rest

import (
	"minesweeper/internal/platform/game"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine            *gin.Engine
	Game              game.Game
	GameHolder        game.Holder
	CreateGameHandler *CreateGameHandler
	StartGameHandler  *StartGameHandler
	ShowGameHandler   *ShowGameHandler
	ModifyTileHandler *ModifyTileHandler
}
