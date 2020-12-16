package provide

import (
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/rest"

	"github.com/gin-gonic/gin"
)

func GinEngine() *gin.Engine {
	return gin.Default()
}

func Game() game.Game {
	return game.NewInMemory(nil)
}

func StartGameHandler(game game.Game) *rest.StartGameHandler {
	return &rest.StartGameHandler{
		Game:        game,
		GameStarter: operation.StartGame{},
	}
}
