package provide

import (
	"minesweeper/internal/operation"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/random"
	"minesweeper/internal/platform/rest"
	"minesweeper/internal/platform/rest/middleware"

	"github.com/gin-gonic/gin"
)

func GinEngine() *gin.Engine {
	r := gin.Default()

	// r.Use(gin.ErrorLogger())
	r.Use(middleware.ErrorLogger())

	return r
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
		GameStarter: operation.StartGame{Rand: random.Real{}},
		GameShower:  showGame,
	}
}

func ShowGameHandler(game game.Game, showGame operation.ShowGame) *rest.ShowGameHandler {
	return &rest.ShowGameHandler{
		Game:       game,
		GameShower: showGame,
	}
}

func TapHandler(game game.Game, showGame operation.ShowGame) *rest.TapHandler {
	return &rest.TapHandler{
		Game: game,
		Tapper: operation.Tap{
			GameFinisher: operation.FinishGame{},
		},
		GameShower: showGame,
	}
}
