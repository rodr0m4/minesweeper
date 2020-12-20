package main

import (
	"log"
	"minesweeper/internal/platform/provide"
	"minesweeper/internal/platform/rest"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	port := "8080"

	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	game := provide.Game()
	boardDrawer := provide.BoardDrawer()

	server := &rest.Server{
		Engine:            provide.GinEngine(),
		Game:              game,
		StartGameHandler:  provide.StartGameHandler(game, boardDrawer),
		ShowGameHandler:   provide.ShowGameHandler(game, boardDrawer),
		ModifyTileHandler: provide.ModifyTileHandler(game, boardDrawer),
	}

	server.RegisterRoutes()

	return server.Engine.Run(":" + port)
}
