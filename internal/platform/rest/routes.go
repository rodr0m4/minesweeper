package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	createGameHandler := s.CreateGameHandler
	showGameHandler := s.ShowGameHandler
	modifyTileHandler := s.ModifyTileHandler

	r.POST("/games", createGameHandler.CreateGame)
	r.GET("/games/:id", showGameHandler.ShowGame)
	r.POST("/games/:id/tap", modifyTileHandler.Tap)
	r.POST("/games/:id/mark", modifyTileHandler.Mark)
}
