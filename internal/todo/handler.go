package todo

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/harundurmus/go-to-do-app/pkg/errors"
	"go.uber.org/zap"
)

// TodoService defines the interface for operations related to TODOs.
type TodoService interface {
	InsertOrUpdateTodo(todo Todo) error
	DeleteTodo(ID string) error
}

// Handler handles HTTP requests related to TODO operations.
type Handler struct {
	todoService TodoService
	logger      *zap.Logger
}

// NewHandler creates a new instance of Handler with the given TodoService and logger.
func NewHandler(todoService TodoService, logger *zap.Logger) *Handler {
	return &Handler{
		todoService: todoService,
		logger:      logger,
	}
}

// CreateOrUpdate handles the creation or update of a TODO item.
// @Summary  Create Or Update todo data
// @Description Create Or Update todo data
// @Accept */*
// @Produce json
// @Param body body Todo true "Todo"
// @Router /api/v1/todo [post]
func (h *Handler) CreateOrUpdate(c *fiber.Ctx) error {
	todo := Todo{}
	if err := json.Unmarshal(c.Body(), &todo); err != nil {
		h.logger.Error("request JSON decoding error", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(errors.BadRequest(err.Error()))
	}

	if err := h.todoService.InsertOrUpdateTodo(todo); err != nil {
		h.logger.Error("todo creation error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.InternalServerError(err.Error()))
	}
	h.logger.Debug("create todo handler executed successfully")

	return c.SendStatus(fiber.StatusCreated)
}

// Delete handles the deletion of a TODO item by its ID.
// @Summary Delete todo data
// @Description Delete todo data by ID
// @Accept */*
// @Param id path string true "TODO ID"
// @Router /api/v1/todo/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	// Get TODO ID from path parameters
	ID := c.Params("id")

	// Call the service layer to delete the TODO item
	if err := h.todoService.DeleteTodo(ID); err != nil {
		h.logger.Error("error deleting todo", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(errors.InternalServerError(err.Error()))
	}

	// Log success message
	h.logger.Debug("delete todo handler executed successfully")

	// Return success response
	return c.SendStatus(fiber.StatusOK)
}

// RegisterRoutes registers the TODO-related routes to the given Fiber app.
func (h *Handler) RegisterRoutes(app *fiber.App) {
	router := app.Group("/api/v1")
	router.Post("/todo", h.CreateOrUpdate)
}
