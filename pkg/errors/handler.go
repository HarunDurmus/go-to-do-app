package errors

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/harundurmus/go-to-do-app/pkg/logger"
	"go.uber.org/zap"
)

func Handler(zapLogger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err != nil {
			if _, ok := err.(*fiber.Error); ok {
				return fiber.DefaultErrorHandler(c, err)
			}
			l := log.FromContext(c.Context())
			res := buildErrorResponse(err)
			l.Error("encountered an error: "+err.Error(), zap.Error(getError(err, res)), zap.Int("status code", res.StatusCode()))

			if err = c.Status(res.StatusCode()).JSON(res); err != nil {
				l.Error("failed writing error response: "+err.Error(), zap.Error(err))
				return err
			}
		}
		return nil
	}
}
