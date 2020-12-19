package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	startGameHandler := s.StartGameHandler
	showGameHandler := s.ShowGameHandler
	tapHandler := s.TapHandler

	r.POST("/games", startGameHandler.StartGame)
	r.GET("/games", showGameHandler.ShowGame)
	r.PATCH("/games", tapHandler.Tap)
}
