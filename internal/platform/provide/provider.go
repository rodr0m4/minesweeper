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

func GameHolder() game.Holder {
	return game.NewInMemoryHolder()
}

func BoardDrawer(revealEverything bool) operation.DefaultBoardDrawer {
	return operation.DefaultBoardDrawer{
		RevealEverything: revealEverything,
		ShowTiles:        true,
		ShowLines:        true,
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

func DeleteGameHandler(holder game.Holder, boardDrawer operation.BoardDrawer) *rest.DeleteGameHandler {
	return &rest.DeleteGameHandler{
		GameDeleter: operation.DeleteGame{
			Holder:      holder,
			BoardDrawer: boardDrawer,
		},
	}
}

func ShowGameHandler(holder game.Holder, boardDrawer operation.BoardDrawer) *rest.ShowGameHandler {
	return &rest.ShowGameHandler{
		GameHolder:  holder,
		BoardDrawer: boardDrawer,
	}
}

func ModifyTileHandler(holder game.Holder, boardDrawer operation.BoardDrawer) *rest.ModifyTileHandler {
	return &rest.ModifyTileHandler{
		GameHolder:  holder,
		BoardDrawer: boardDrawer,
		Tapper: operation.Tap{
			GameFinisher: operation.FinishGame{},
			TileRevealer: operation.RevealAdjacent{},
		},
		Marker: operation.Mark{},
	}
}
