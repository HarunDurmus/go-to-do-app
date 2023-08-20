package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/harundurmus/go-to-do-app/internal/config"
	errs "github.com/harundurmus/go-to-do-app/pkg/errors"
	log "github.com/harundurmus/go-to-do-app/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	app    *fiber.App
	config *config.Server
	logger *zap.Logger
}

func New(serverConfig *config.Server, handlers []Handler, logger *zap.Logger) Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: errs.Handler(zap.NewNop()),
	})

	server := Server{app: app, config: serverConfig, logger: logger}
	server.app.Use(cors.New())
	server.app.Use(log.Middleware(logger))
	server.addRoutes()

	for _, handler := range handlers {
		handler.RegisterRoutes(server.app)
	}

	return server
}

func (s *Server) Run() error {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		err := s.app.Shutdown()
		if err != nil {
			s.logger.Error("Graceful shutdown failed")
		}
	}()

	return s.app.Listen(s.config.Port)
}

func healthCheck(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return nil
}
func (s *Server) addRoutes() {
	s.app.Get("/health", healthCheck)
}
