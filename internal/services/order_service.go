package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"lgtm-lab/internal/models"
	"time"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) error
}

type orderService struct {
	logger zerolog.Logger
}

func NewOrderService(logger zerolog.Logger) OrderService {
	return &orderService{
		logger: logger,
	}
}

func (s *orderService) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	traceID := ctx.Value("trace_id").(string)

	s.logger.Info().
		Str("trace_id", traceID).
		Str("order_id", orderID).
		Str("operation", "get_order").
		Str("service", "order_service").
		Msg("Fetching order")

	// Simulate database fetch
	if orderID == "123" {
		order := &models.Order{
			ID:         orderID,
			CustomerID: "cust_456",
			Amount:     99.99,
			Status:     "completed",
			CreatedAt:  time.Now(),
		}

		s.logger.Info().
			Str("trace_id", traceID).
			Str("order_id", orderID).
			Str("operation", "get_order").
			Str("service", "order_service").
			Str("status", "success").
			Msg("Order found successfully")

		return order, nil
	}

	s.logger.Warn().
		Str("trace_id", traceID).
		Str("order_id", orderID).
		Str("operation", "get_order").
		Str("service", "order_service").
		Str("status", "not_found").
		Msg("Order not found")

	return nil, fmt.Errorf("order not found")
}

func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	traceID := ctx.Value("trace_id").(string)

	s.logger.Info().
		Str("trace_id", traceID).
		Str("operation", "create_order").
		Str("service", "order_service").
		Str("customer_id", order.CustomerID).
		Float64("amount", order.Amount).
		Msg("Creating new order")

	// Simulate order creation
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.Status = "pending"

	s.logger.Info().
		Str("trace_id", traceID).
		Str("order_id", order.ID).
		Str("operation", "create_order").
		Str("service", "order_service").
		Str("status", "success").
		Msg("Order created successfully")

	return nil
}
