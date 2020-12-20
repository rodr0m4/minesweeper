package main

import (
	"fmt"
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

	gameHolder := provide.GameHolder()
	boardDrawer := provide.BoardDrawer(shouldRevealEverything())

	server := &rest.Server{
		Engine:            provide.GinEngine(),
		GameHolder:        gameHolder,
		CreateGameHandler: provide.CreateGameHandler(gameHolder, boardDrawer),
		DeleteGameHandler: provide.DeleteGameHandler(gameHolder, boardDrawer),
		ShowGameHandler:   provide.ShowGameHandler(gameHolder, boardDrawer),
		ModifyTileHandler: provide.ModifyTileHandler(gameHolder, boardDrawer),
	}

	server.RegisterRoutes()

	return server.Engine.Run(":" + port)
}

func shouldRevealEverything() bool {
	env := strings.ToLower(os.Getenv("REVEAL_EVERYTHING"))
	fmt.Printf("env is %s", env)

	switch env {
	case "", "false":
		return false
	default:
		return true
	}
}
