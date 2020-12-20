package main

import (
	"log"
	"minesweeper/internal/platform/provide"
	"minesweeper/internal/platform/rest"
	"os"
	"strings"
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
	gameHolder := provide.GameHolder()
	boardDrawer := provide.BoardDrawer(shouldRevealEverything())

	server := &rest.Server{
		Engine:            provide.GinEngine(),
		Game:              game,
		GameHolder:        gameHolder,
		CreateGameHandler: provide.CreateGameHandler(gameHolder, boardDrawer),
		StartGameHandler:  provide.StartGameHandler(game, boardDrawer),
		ShowGameHandler:   provide.ShowGameHandler(game, boardDrawer),
		ModifyTileHandler: provide.ModifyTileHandler(game, boardDrawer),
	}

	server.RegisterRoutes()

	return server.Engine.Run(":" + port)
}

func shouldRevealEverything() bool {
	env := strings.ToLower(os.Getenv("REVEAL_EVERYTHING"))

	switch env {
	case "", "false":
		return false
	default:
		return true
	}
}
