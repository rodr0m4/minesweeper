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

	gameHolder := provide.GameHolder()
	boardDrawer := provide.BoardDrawer(
		parseBooleanEnvVar("REVEAL_EVERYTHING", false),
		parseBooleanEnvVar("SHOW_TILES", false),
		parseBooleanEnvVar("SHOW_LINES", true),
	)

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

func parseBooleanEnvVar(name string, defaultValue bool) bool {
	env := strings.ToLower(os.Getenv(name))

	switch env {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}
