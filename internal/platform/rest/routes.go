package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	startGameHandler := s.StartGameHandler
	showGameHandler := s.ShowGameHandler
	modifyTileHandler := s.ModifyTileHandler

	r.POST("/game", startGameHandler.StartGame)
	r.GET("/game", showGameHandler.ShowGame)
	r.POST("/game/tap", modifyTileHandler.Tap)
	r.POST("/game/mark", modifyTileHandler.Mark)
}
