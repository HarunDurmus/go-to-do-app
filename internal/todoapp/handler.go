package todoapp

import (
	"context"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type (
	TaskService interface {
		InitializeData(ctx context.Context) error
		CreateTaskData(ctx context.Context, model ToDoList) (*ToDoList, error)
	}

	Handler struct {
		logger  *zap.Logger
		service TaskService
	}
)

func NewHandler(logger *zap.Logger, taskService TaskService) *Handler {
	handler := &Handler{
		logger:  logger,
		service: taskService,
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

// -------------------
// @Summary  -------------------------
// @Description -----------------
// @Accept */*
// @Produce json
// @Success 200
// @Failure 400 {object} Response
// @Router /create-task [post]
// @param id body ToDoList true "User"
func (h *Handler) CreateTaskData(ctx echo.Context) error {
	output := make(chan bool, 1)
	resp := make(chan *ToDoList)
	errs := hystrix.Go("go-to-do-app", func() error {
		var too ToDoList
		err := ctx.Bind(&too)
		if err != nil {
			h.logger.Error("Given coordinate could not be mapped to Coordinate struct")
			return err
		}
		result, err := h.service.CreateTaskData(ctx.Request().Context(), too)
		if err != nil {
			h.logger.Error("the nearest point could not be calculate", zap.Error(err))
			return err
		}
		h.logger.Info("FindNearestPoint handler executed successfully")
		output <- true
		resp <- result
		return err
	}, nil)
	select {
	case _ = <-output:
		return ctx.JSON(http.StatusOK, Response{Data: <-resp})
	case err := <-errs:
		h.logger.Error("Failed to search nearest point", zap.Error(err))
		return ctx.JSON(http.StatusNotFound, Response{Error: err.Error()})
	}
}

//// TaskComplete update task route
//func (h *Handler) TaskComplete(w http.ResponseWriter, r *http.Request) {
//
//}
//
//// UndoTask undo the complete task route
//func (h *Handler) UndoTask(w http.ResponseWriter, r *http.Request) {
//
//}
//
//// DeleteTask delete one task route
//func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
//
//}
//
//// DeleteAllTask delete all tasks route
//func (h *Handler) DeleteAllTask(w http.ResponseWriter, r *http.Request) {
//
//}
//
//// get all task from the DB and return it
//func (h *Handler) getAllTask() []primitive.M {
//
//}
//
//// Insert one task in the DB
//func (h *Handler) insertOneTask() {
//
//}
//
//// task complete method, update task's status to true
//func (h *Handler) taskComplete(task string) {
//
//}
//
//// task undo method, update task's status to false
//func (h *Handler) undoTask(task string) {
//
//}
//
//// delete one task from the DB, delete by ID
//func (h *Handler) deleteOneTask(task string) {
//
//}
//
//// delete all the tasks from the DB
//func (h *Handler) deleteAllTask() int64 {
//
//}
