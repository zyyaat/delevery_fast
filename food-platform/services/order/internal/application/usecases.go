// Package application contains use cases for the Order Service.
package application

import (
	"context"
	"fmt"

	"github.com/food-platform/services/order/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	FindByOrderNumber(ctx context.Context, number string) (*domain.Order, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error)
	FindActiveByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error)
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID, limit, offset int) ([]*domain.Order, error)
	FindActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
}

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error
	PublishOrderConfirmed(ctx context.Context, event OrderConfirmedEvent) error
	PublishOrderCancelled(ctx context.Context, event OrderCancelledEvent) error
	PublishOrderStatusChanged(ctx context.Context, event OrderStatusChangedEvent) error
}

// ============ DTOs ============

type OrderItemDTO struct {
	ID          uuid.UUID `json:"id"`
	MenuItemID  uuid.UUID `json:"menu_item_id"`
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	LineTotal   float64   `json:"line_total"`
	Notes       string    `json:"notes,omitempty"`
}

type OrderDTO struct {
	ID              uuid.UUID       `json:"id"`
	OrderNumber     string          `json:"order_number"`
	CustomerID      uuid.UUID       `json:"customer_id"`
	RestaurantID    uuid.UUID       `json:"restaurant_id"`
	DriverID        *uuid.UUID      `json:"driver_id,omitempty"`
	Status          string          `json:"status"`
	Items           []OrderItemDTO  `json:"items"`
	Subtotal        float64         `json:"subtotal"`
	DeliveryFee     float64         `json:"delivery_fee"`
	ServiceFee      float64         `json:"service_fee"`
	VAT             float64         `json:"vat"`
	Discount        float64         `json:"discount"`
	Total           float64         `json:"total"`
	PaymentMethod   string          `json:"payment_method"`
	PaymentStatus   string          `json:"payment_status"`
	DeliveryAddress string          `json:"delivery_address"`
	ETAMinutes      int             `json:"eta_minutes"`
	CancelReason    string          `json:"cancel_reason,omitempty"`
	Notes           string          `json:"notes,omitempty"`
	CashbackEarned  float64         `json:"cashback_earned,omitempty"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

// ============ Events ============

type OrderCreatedEvent struct {
	EventID      uuid.UUID
	OrderID      uuid.UUID
	CustomerID   uuid.UUID
	RestaurantID uuid.UUID
	TotalAmount  float64
	PaymentMethod string
	ItemsCount   int
}

type OrderConfirmedEvent struct {
	EventID   uuid.UUID
	OrderID   uuid.UUID
	DriverID  *uuid.UUID
	ETAMinutes int
}

type OrderCancelledEvent struct {
	EventID    uuid.UUID
	OrderID    uuid.UUID
	Reason     string
	CancelledBy string
}

type OrderStatusChangedEvent struct {
	EventID       uuid.UUID
	OrderID       uuid.UUID
	PreviousStatus string
	NewStatus      string
}

// ============ Commands ============

type CreateOrderItemCommand struct {
	MenuItemID uuid.UUID
	Name       string
	Quantity   int
	UnitPrice  float64
	Modifiers  string
	Notes      string
}

type CreateOrderCommand struct {
	CustomerID      uuid.UUID
	RestaurantID    uuid.UUID
	Items           []CreateOrderItemCommand
	DeliveryAddress string
	Latitude        float64
	Longitude       float64
	PaymentMethod   domain.PaymentMethod
	DeliveryFee     float64
	ServiceFeeRate  float64
	VATRate         float64
	Discount        float64
	Notes           string
}

type CancelOrderCommand struct {
	OrderID uuid.UUID
	Reason  string
}

type UpdateStatusCommand struct {
	OrderID uuid.UUID
	Status  domain.OrderStatus
}

// ============ Use Cases ============

type CreateOrderUseCase struct {
	repo     OrderRepository
	publisher EventPublisher
}

func NewCreateOrderUseCase(repo OrderRepository, publisher EventPublisher) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: repo, publisher: publisher}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*OrderDTO, error) {
	// Convert items
	items := make([]*domain.OrderItem, len(cmd.Items))
	for i, itemCmd := range cmd.Items {
		item, err := domain.NewOrderItem(
			itemCmd.MenuItemID, itemCmd.Name, itemCmd.Quantity,
			itemCmd.UnitPrice, itemCmd.Modifiers, itemCmd.Notes,
		)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	// Create order
	order, err := domain.NewOrder(
		cmd.CustomerID, cmd.RestaurantID, items,
		cmd.DeliveryAddress, cmd.Latitude, cmd.Longitude,
		cmd.PaymentMethod, cmd.DeliveryFee, cmd.ServiceFeeRate,
		cmd.VATRate, cmd.Discount,
	)
	if err != nil {
		return nil, err
	}

	if cmd.Notes != "" {
		order.SetNotes(cmd.Notes)
	}

	// Save
	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// Publish event
	_ = uc.publisher.PublishOrderCreated(ctx, OrderCreatedEvent{
		EventID:      uuid.New(),
		OrderID:      order.ID(),
		CustomerID:   order.CustomerID(),
		RestaurantID: order.RestaurantID(),
		TotalAmount:  order.Total(),
		PaymentMethod: string(order.PaymentMethod()),
		ItemsCount:   len(order.Items()),
	})

	return toOrderDTO(order), nil
}

type GetOrderUseCase struct {
	repo OrderRepository
}

func NewGetOrderUseCase(repo OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{repo: repo}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, id uuid.UUID) (*OrderDTO, error) {
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toOrderDTO(order), nil
}

type GetActiveOrdersUseCase struct {
	repo OrderRepository
}

func NewGetActiveOrdersUseCase(repo OrderRepository) *GetActiveOrdersUseCase {
	return &GetActiveOrdersUseCase{repo: repo}
}

func (uc *GetActiveOrdersUseCase) Execute(ctx context.Context, customerID uuid.UUID) ([]*OrderDTO, error) {
	orders, err := uc.repo.FindActiveByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*OrderDTO, 0, len(orders))
	for _, o := range orders {
		dtos = append(dtos, toOrderDTO(o))
	}
	return dtos, nil
}

type GetOrderHistoryUseCase struct {
	repo OrderRepository
}

func NewGetOrderHistoryUseCase(repo OrderRepository) *GetOrderHistoryUseCase {
	return &GetOrderHistoryUseCase{repo: repo}
}

func (uc *GetOrderHistoryUseCase) Execute(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*OrderDTO, error) {
	if limit == 0 {
		limit = 20
	}
	orders, err := uc.repo.FindByCustomer(ctx, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	dtos := make([]*OrderDTO, 0, len(orders))
	for _, o := range orders {
		dtos = append(dtos, toOrderDTO(o))
	}
	return dtos, nil
}

type CancelOrderUseCase struct {
	repo     OrderRepository
	publisher EventPublisher
}

func NewCancelOrderUseCase(repo OrderRepository, publisher EventPublisher) *CancelOrderUseCase {
	return &CancelOrderUseCase{repo: repo, publisher: publisher}
}

func (uc *CancelOrderUseCase) Execute(ctx context.Context, cmd CancelOrderCommand) error {
	order, err := uc.repo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(cmd.Reason); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	_ = uc.publisher.PublishOrderCancelled(ctx, OrderCancelledEvent{
		EventID:    uuid.New(),
		OrderID:    order.ID(),
		Reason:     cmd.Reason,
		CancelledBy: "customer",
	})

	return nil
}

type UpdateOrderStatusUseCase struct {
	repo     OrderRepository
	publisher EventPublisher
}

func NewUpdateOrderStatusUseCase(repo OrderRepository, publisher EventPublisher) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{repo: repo, publisher: publisher}
}

func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, cmd UpdateStatusCommand) error {
	order, err := uc.repo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	prevStatus := order.Status()
	if err := order.TransitionTo(cmd.Status); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	_ = uc.publisher.PublishOrderStatusChanged(ctx, OrderStatusChangedEvent{
		EventID:        uuid.New(),
		OrderID:        order.ID(),
		PreviousStatus: string(prevStatus),
		NewStatus:      string(cmd.Status),
	})

	return nil
}

// ============ Helpers ============

func toOrderDTO(order *domain.Order) *OrderDTO {
	items := make([]OrderItemDTO, len(order.Items()))
	for i, item := range order.Items() {
		items[i] = OrderItemDTO{
			ID:         item.ID(),
			MenuItemID: item.MenuItemID(),
			Name:       item.Name(),
			Quantity:   item.Quantity(),
			UnitPrice:  item.UnitPrice(),
			LineTotal:  item.LineTotal(),
			Notes:      item.Notes(),
		}
	}

	dto := &OrderDTO{
		ID:              order.ID(),
		OrderNumber:     order.OrderNumber(),
		CustomerID:      order.CustomerID(),
		RestaurantID:    order.RestaurantID(),
		DriverID:        order.DriverID(),
		Status:          string(order.Status()),
		Items:           items,
		Subtotal:        order.Subtotal(),
		DeliveryFee:     order.DeliveryFee(),
		ServiceFee:      order.ServiceFee(),
		VAT:             order.VAT(),
		Discount:        order.Discount(),
		Total:           order.Total(),
		PaymentMethod:   string(order.PaymentMethod()),
		PaymentStatus:   string(order.PaymentStatus()),
		DeliveryAddress: order.DeliveryAddress(),
		ETAMinutes:      order.ETAMinutes(),
		CancelReason:    order.CancelReason(),
		Notes:           order.Notes(),
		CashbackEarned:  order.CashbackEarned(),
		CreatedAt:       order.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       order.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	return dto
}
