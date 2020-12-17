package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	startGameHandler := s.StartGameHandler
	showGameHandler := s.ShowGameHandler

	r.POST("/games", startGameHandler.StartGame)
	r.GET("/games", showGameHandler.ShowGame)
}
