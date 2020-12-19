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
	showGame := provide.ShowGame()

	server := &rest.Server{
		Engine:           provide.GinEngine(),
		Game:             game,
		StartGameHandler: provide.StartGameHandler(game, showGame),
		ShowGameHandler:  provide.ShowGameHandler(game, showGame),
		TapHandler:       provide.TapHandler(game, showGame),
	}

	server.RegisterRoutes()

	return server.Engine.Run(":" + port)
}
