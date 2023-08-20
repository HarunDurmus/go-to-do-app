package log

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var loggerKey = "logger"

// Middleware adds logger to fiber context
// after using this middleware, FromContext can be used to get logger
// logger := log.FromContext(c.Context())
// using this logger, all further logs will contain requestID
func Middleware(logger *zap.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("x-request-id", uuid.New().String())
		c.Locals(loggerKey, logger.With(zap.String("requestID", requestID)))
		return c.Next()
	}
}

// Inject adds logger to given context
func Inject(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger
func FromContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		if len(fields) > 0 {
			return l.With(fields...)
		}
		return l
	}

	return zap.New(zap.L().Core())
}

func AddFields(c *fiber.Ctx, fields ...zap.Field) {
	if l, ok := c.Context().Value(loggerKey).(*zap.Logger); ok {
		newLogger := l.With(fields...)
		c.Locals(loggerKey, newLogger)
	}
}
