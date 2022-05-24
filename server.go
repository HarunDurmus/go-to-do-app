package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/harundurmus/go-to-do-app/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
)

type Server struct {
	e      *echo.Echo
	config *config.Config
}

var timeout = 10 * time.Second

func NewServer(c *config.Config) *Server {
	server := &Server{}
	e := echo.New()
	e.Use(middleware.Recover())
	server.e = e
	server.config = c

	return server
}

func (s *Server) Start() error {
	s.e.Server.Addr = fmt.Sprintf(":%d", s.config.Server.Port)

	s.e.GET("/health", s.healthCheck)

	return graceful.ListenAndServe(s.e.Server, timeout)
}

func (s *Server) healthCheck(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
