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

func ShowGame() operation.ShowGame {
	return operation.ShowGame{
		BoardDrawer: operation.BoardDrawer{},
	}
}

func StartGameHandler(game game.Game, showGame operation.ShowGame) *rest.StartGameHandler {
	return &rest.StartGameHandler{
		Game:        game,
		GameStarter: operation.StartGame{},
		GameShower:  showGame,
	}
}

func ShowGameHandler(game game.Game, showGame operation.ShowGame) *rest.ShowGameHandler {
	return &rest.ShowGameHandler{
		Game:       game,
		GameShower: showGame,
	}
}
