package rest

func (s *Server) RegisterRoutes() {
	r := s.Engine
	startGameHandler := s.StartGameHandler

	r.POST("/games", startGameHandler.StartGame)
}
