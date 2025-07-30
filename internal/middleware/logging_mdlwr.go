package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func TraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceID := uuid.New().String()
			c.Set("trace_id", traceID)
			c.Response().Header().Set("X-Trace-ID", traceID)
			return next(c)
		}
	}
}

func LoggingMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			traceID := c.Get("trace_id")
			if traceID == nil {
				traceID = "unknown"
			}

			logger.Info().
				Str("trace_id", traceID.(string)).
				Str("method", values.Method).
				Str("uri", values.URI).
				Int("status", values.Status).
				Int64("latency_ms", values.Latency.Milliseconds()).
				Str("remote_ip", values.RemoteIP).
				Str("component", "http_access_log").
				Msg("HTTP request processed")

			return nil
		},
	})
}
