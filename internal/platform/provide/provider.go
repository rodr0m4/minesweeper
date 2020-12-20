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

func BoardDrawer() operation.DefaultBoardDrawer {
	return operation.DefaultBoardDrawer{}
}

func StartGameHandler(game game.Game, boardDrawer operation.BoardDrawer) *rest.StartGameHandler {
	return &rest.StartGameHandler{
		Game:        game,
		GameStarter: operation.StartGame{Rand: random.Real{}},
		BoardDrawer: boardDrawer,
	}
}

func ShowGameHandler(game game.Game, boardDrawer operation.BoardDrawer) *rest.ShowGameHandler {
	return &rest.ShowGameHandler{
		Game:        game,
		BoardDrawer: boardDrawer,
	}
}

func TapHandler(game game.Game, boardDrawer operation.BoardDrawer) *rest.TapHandler {
	return &rest.TapHandler{
		Game: game,
		Tapper: operation.Tap{
			GameFinisher: operation.FinishGame{},
			TileRevealer: operation.RevealAdjacent{},
		},
		BoardDrawer: boardDrawer,
	}
}
