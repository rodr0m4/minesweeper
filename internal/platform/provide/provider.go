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

func GameHolder() game.Holder {
	return game.NewInMemoryHolder()
}

func BoardDrawer(revealEverything bool) operation.DefaultBoardDrawer {
	return operation.DefaultBoardDrawer{RevealEverything: revealEverything}
}

func CreateGame(holder game.Holder) operation.CreateGame {
	return operation.CreateGame{
		Holder: holder,
		Rand:   random.Real{},
	}
}

func CreateGameHandler(holder game.Holder, boardDrawer operation.BoardDrawer) *rest.CreateGameHandler {
	return &rest.CreateGameHandler{
		GameCreator: operation.CreateGame{
			Holder: holder,
			Rand:   random.Real{},
		},
		BoardDrawer: boardDrawer,
	}
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

func ModifyTileHandler(game game.Game, boardDrawer operation.BoardDrawer) *rest.ModifyTileHandler {
	return &rest.ModifyTileHandler{
		Game:        game,
		BoardDrawer: boardDrawer,
		Tapper: operation.Tap{
			GameFinisher: operation.FinishGame{},
			TileRevealer: operation.RevealAdjacent{},
		},
		Marker: nil,
	}
}
