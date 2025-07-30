package main

import (
	"context"
	"lgtm-lab/internal/controllers"
	mdlwr "lgtm-lab/internal/middleware"
	"lgtm-lab/internal/models"
	"lgtm-lab/internal/services"
	"lgtm-lab/internal/utils"
	mtrx "lgtm-lab/metrics"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthHandler(c echo.Context) error {
	traceID := c.Get("trace_id").(string)
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service is healthy",
		TraceID: traceID,
	})
}

func helpHandler(c echo.Context) error {
	traceID := c.Get("trace_id").(string)
	help := map[string]string{
		"GET /help":              "Show this help message",
		"GET /health":            "Health check endpoint",
		"GET /metrics":           "Prometheus metrics endpoint (for monitoring)",
		"GET /api/v1/orders/:id": "Get order by ID (try ID: 123)",
		"POST /api/v1/orders":    "Create a new order",
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    help,
		TraceID: traceID,
	})
}

func main() {
	// Setup Zerolog
	logger := utils.SetupZerologLogger()

	logger.Info().
		Str("component", "application").
		Str("version", "1.0.0").
		Str("stage", "startup").
		Str("logger", "zerolog").
		Msg("Starting Order Management API with Zerolog")

	// Initialize services with dependency injection
	orderService := services.NewOrderService(logger)
	orderController := controllers.NewOrderController(orderService, logger)
	// Setup Echo framework
	e := echo.New()
	e.HideBanner = true

	// Apply middleware
	e.Use(mdlwr.TraceIDMiddleware())       // Generate trace IDs for request tracking
	e.Use(mdlwr.LoggingMiddleware(logger)) // Structured HTTP access logs with Zerolog
	e.Use(mtrx.MetricsMiddleware())        // Prometheus metrics collection
	e.Use(middleware.Recover())            // Panic recovery
	e.Use(middleware.CORS())               // CORS headers

	// Routes
	e.GET("/help", helpHandler)
	e.GET("/health", healthHandler)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler())) // Prometheus metrics endpoint

	// API routes
	api := e.Group("/api/v1")
	api.GET("/orders/:id", orderController.GetOrder)
	api.POST("/orders", orderController.CreateOrder)

	// Start server
	logger.Info().
		Str("component", "application").
		Int("port", 8080).
		Str("stage", "startup").
		Msg("Server starting on port 8080")

	// Graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			logger.Fatal().
				Str("component", "application").
				Err(err).
				Str("stage", "startup").
				Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().
		Str("component", "application").
		Str("stage", "shutdown").
		Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().
			Str("component", "application").
			Err(err).
			Str("stage", "shutdown").
			Msg("Server forced to shutdown")
	}

	logger.Info().
		Str("component", "application").
		Str("stage", "shutdown").
		Msg("Server exited gracefully")
}
