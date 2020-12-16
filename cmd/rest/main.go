package main

import (
	"log"
	"minesweeper/internal/platform/provide"
	"minesweeper/internal/platform/rest"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	game := provide.Game()

	server := &rest.Server{
		Engine:           provide.GinEngine(),
		Game:             game,
		StartGameHandler: provide.StartGameHandler(game),
	}

	server.RegisterRoutes()

	return server.Engine.Run()
}
