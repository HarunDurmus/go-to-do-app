package main

import (
	"github.com/harundurmus/go-to-do-app/internal/config"
	"github.com/harundurmus/go-to-do-app/internal/todo"
	"github.com/harundurmus/go-to-do-app/pkg/couchbase"
	log "github.com/harundurmus/go-to-do-app/pkg/logger"
	"github.com/harundurmus/go-to-do-app/pkg/server"
	"os"

	"github.com/yudai/pp"
	"go.uber.org/zap"
)

const (
	todoAppBucketName = "todoapp"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		err := os.Setenv("APP_ENV", "dev")
		if err != nil {
			panic(err)
		}
	}
	conf, err := config.New(".config", os.Getenv("APP_ENV"))
	if err != nil {

		panic(err)
	}
	_, _ = pp.Println(conf)
	logger := log.NewWith(log.Config{
		Level: conf.LogLevel,
	})
	cb, err := couchbase.New(conf.Couchbase)
	if err != nil {
		_, _ = pp.Println(err)
		panic(err)
	}
	goTodoApp, err := cb.Bucket(todoAppBucketName)
	if err != nil {
		panic(err)
	}
	todoRepository := todo.NewRepository(cb.Cluster(), goTodoApp)
	todoService := todo.NewService(todoRepository, logger)
	todoHandler := todo.NewHandler(&todoService, logger)

	runServer(&conf.Server, logger, todoHandler)
}

func runServer(conf *config.Server, logger *zap.Logger, handlers ...server.Handler) {
	s := server.New(conf, handlers, logger)
	if err := s.Run(); err != nil {
		panic(err)
	}
}
