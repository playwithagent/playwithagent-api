package httpserver

import (
	"fmt"
	"playwithagent-xo/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config config.Config
	
}

func NewServer(config config.Config) *Server {	
	return &Server{
		config: config,
	}
}

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// gameGroup := e.Group("/game")

	e.GET("/health", s.healthCheck)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
	// gameGroup.POST("/start", s.startGame)
}