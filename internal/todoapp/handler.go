package todoapp

import (
	"context"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	LocationService interface {
		InitializeData(ctx context.Context) error
	}

	Handler struct {
		logger  *zap.Logger
		service LocationService
	}
)

func NewHandler(logger *zap.Logger, driverlocationService LocationService) *Handler {
	handler := &Handler{
		logger:  logger,
		service: driverlocationService,
	}
	hystrix.ConfigureCommand("go-to-do-ap", hystrix.CommandConfig{
		Timeout:               5000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	return handler
}

// Import data from csv file godoc
// @Summary  Import data from csv file.
// @Description Import data from uploaded csv file.
// @Accept */*
// @Produce json
// @Success 200
// @Failure 400 {object} Response
// @Router /init [post]
func (h *Handler) InitializeDatabase(ctx echo.Context) error {
	output := make(chan bool, 1)
	errs := hystrix.Go("go-to-do-app", func() error {

		h.logger.Info("Initializing data ...")
		err := h.service.InitializeData(ctx.Request().Context())
		if err != nil {
			h.logger.Error("Data could not be imported", zap.Error(err))
			return err
		}
		h.logger.Info("Data successfully imported")
		output <- true
		return err
	}, nil)
	select {
	case _ = <-output:
		return ctx.JSON(http.StatusOK, http.NoBody)

	case err := <-errs:
		h.logger.Error("Failed to initialize data", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, Response{Error: err.Error()})

	}
}
