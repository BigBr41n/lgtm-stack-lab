package controllers

import (
	"context"
	"lgtm-lab/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"lgtm-lab/internal/models"
)

type OrderController struct {
	orderService services.OrderService
	logger       zerolog.Logger
}

func NewOrderController(orderService services.OrderService, logger zerolog.Logger) *OrderController {
	return &OrderController{
		orderService: orderService,
		logger:       logger,
	}
}

func (oc *OrderController) GetOrder(c echo.Context) error {
	traceID := c.Get("trace_id").(string)
	orderID := c.Param("id")

	oc.logger.Info().
		Str("trace_id", traceID).
		Str("order_id", orderID).
		Str("operation", "get_order_endpoint").
		Str("controller", "order_controller").
		Str("method", "GET").
		Str("path", c.Path()).
		Msg("Received get order request")

	ctx := context.WithValue(c.Request().Context(), "trace_id", traceID)
	order, err := oc.orderService.GetOrder(ctx, orderID)

	if err != nil {
		oc.logger.Error().
			Str("trace_id", traceID).
			Str("order_id", orderID).
			Str("operation", "get_order_endpoint").
			Str("controller", "order_controller").
			Err(err).
			Str("status", "error").
			Msg("Failed to get order")

		return c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Order not found",
			TraceID: traceID,
		})
	}

	oc.logger.Info().
		Str("trace_id", traceID).
		Str("order_id", orderID).
		Str("operation", "get_order_endpoint").
		Str("controller", "order_controller").
		Str("status", "success").
		Msg("Order retrieved successfully")

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    order,
		TraceID: traceID,
	})
}

func (oc *OrderController) CreateOrder(c echo.Context) error {
	traceID := c.Get("trace_id").(string)

	var order models.Order
	if err := c.Bind(&order); err != nil {
		oc.logger.Error().
			Str("trace_id", traceID).
			Str("operation", "create_order_endpoint").
			Str("controller", "order_controller").
			Err(err).
			Str("status", "bind_error").
			Msg("Failed to bind request")

		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			TraceID: traceID,
		})
	}

	ctx := context.WithValue(c.Request().Context(), "trace_id", traceID)
	if err := oc.orderService.CreateOrder(ctx, &order); err != nil {
		oc.logger.Error().
			Str("trace_id", traceID).
			Str("operation", "create_order_endpoint").
			Str("controller", "order_controller").
			Err(err).
			Str("status", "error").
			Msg("Failed to create order")

		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create order",
			TraceID: traceID,
		})
	}

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    order,
		TraceID: traceID,
	})
}
