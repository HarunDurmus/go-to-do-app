package todo

import (
	"bytes"
	"encoding/json"
	"github.com/harundurmus/go-to-do-app/pkg/errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TodoService interface {
	InsertOrUpdateTodo(todo Todo) error
}

type Handler struct {
	todoService TodoService
	logger      *zap.Logger
}

func NewHandler(todoService TodoService, logger *zap.Logger) *Handler {
	return &Handler{
		todoService: todoService,
		logger:      logger,
	}
}

func (h *Handler) CreateOrUpdate(c *fiber.Ctx) error {
	todo := Todo{}
	requestBody := c.Body()
	h.logger.Debug("create todo request arrived", zap.ByteString("requestBody", requestBody))
	if err := json.NewDecoder(bytes.NewBuffer(requestBody)).Decode(&todo); err != nil {
		h.logger.Error("request json decoding error", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(errors.BadRequest(err.Error()))
	}

	if err := h.todoService.InsertOrUpdateTodo(todo); err != nil {
		h.logger.Error("todo creation error", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(errors.InternalServerError(err.Error()))
	}
	h.logger.Debug("create todo handler executed successfully")

	return c.SendStatus(http.StatusCreated)
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	router := app.Group("api/v1")
	router.Post("/todo", h.CreateOrUpdate)
}
