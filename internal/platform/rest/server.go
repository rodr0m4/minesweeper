package rest

import (
	"minesweeper/internal/platform/game"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine           *gin.Engine
	Game             game.Game
	StartGameHandler *StartGameHandler
	ShowGameHandler  *ShowGameHandler
	TapHandler       *TapHandler
}
