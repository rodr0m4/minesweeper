package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	startGameHandler := s.StartGameHandler
	showGameHandler := s.ShowGameHandler
	tapHandler := s.TapHandler

	r.POST("/game", startGameHandler.StartGame)
	r.GET("/game", showGameHandler.ShowGame)
	r.POST("/game/tap", tapHandler.Tap)
}
