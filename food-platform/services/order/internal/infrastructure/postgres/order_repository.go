// Package postgres implements the Order repository using PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/services/order/internal/domain"
	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (
			id, order_number, customer_id, restaurant_id, driver_id, status,
			subtotal, delivery_fee, service_fee, vat, discount, total,
			payment_method, payment_status, delivery_address, latitude, longitude,
			eta_minutes, prep_started_at, picked_up_at, delivered_at,
			cancel_reason, notes, cashback_earned, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	`
	var driverID interface{}
	if order.DriverID() != nil {
		driverID = *order.DriverID()
	}

	_, err = tx.ExecContext(ctx, orderQuery,
		order.ID(), order.OrderNumber(), order.CustomerID(), order.RestaurantID(),
		driverID, string(order.Status()),
		order.Subtotal(), order.DeliveryFee(), order.ServiceFee(), order.VAT(),
		order.Discount(), order.Total(),
		string(order.PaymentMethod()), string(order.PaymentStatus()),
		order.DeliveryAddress(), order.Latitude(), order.Longitude(),
		order.ETAMinutes(), order.PrepStartedAt(), order.PickedUpAt(), order.DeliveredAt(),
		order.CancelReason(), order.Notes(), order.CashbackEarned(),
		order.CreatedAt(), order.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	// Insert items
	itemQuery := `
		INSERT INTO order_items (id, order_id, menu_item_id, name, quantity, unit_price, modifiers, notes, line_total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	for _, item := range order.Items() {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID(), item.OrderID(), item.MenuItemID(),
			item.Name(), item.Quantity(), item.UnitPrice(),
			item.Modifiers(), item.Notes(), item.LineTotal(),
		)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	order, err := scanOrder(row)
	if err != nil {
		return nil, err
	}

	// Load items
	items, err := r.findItemsByOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	// Attach items to order
	for _, item := range items {
		// We need to set items via reflection or a setter — for now, reconstruct
		_ = item
	}

	return order, nil
}

func (r *OrderRepository) FindByOrderNumber(ctx context.Context, number string) (*domain.Order, error) {
	query := `SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders WHERE order_number = $1`
	row := r.db.QueryRowContext(ctx, query, number)
	return scanOrder(row)
}

func (r *OrderRepository) FindByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
	query := `SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders WHERE customer_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *OrderRepository) FindActiveByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	query := `SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders WHERE customer_id = $1 AND status IN ('pending', 'confirmed', 'preparing', 'ready', 'picked_up')
		ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *OrderRepository) FindByRestaurant(ctx context.Context, restaurantID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
	query := `SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders WHERE restaurant_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, restaurantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *OrderRepository) FindActiveByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Order, error) {
	query := `SELECT id, order_number, customer_id, restaurant_id, driver_id, status,
		       subtotal, delivery_fee, service_fee, vat, discount, total,
		       payment_method, payment_status, delivery_address, latitude, longitude,
		       eta_minutes, prep_started_at, picked_up_at, delivered_at,
		       cancel_reason, notes, cashback_earned, created_at, updated_at
		FROM orders WHERE restaurant_id = $1 AND status IN ('pending', 'confirmed', 'preparing', 'ready')
		ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanOrders(rows)
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `UPDATE orders SET
		driver_id = $2, status = $3,
		prep_started_at = $4, picked_up_at = $5, delivered_at = $6,
		cancel_reason = $7, notes = $8, cashback_earned = $9,
		payment_status = $10, eta_minutes = $11, updated_at = $12
		WHERE id = $1`

	var driverID interface{}
	if order.DriverID() != nil {
		driverID = *order.DriverID()
	}

	_, err := r.db.ExecContext(ctx, query,
		order.ID(), driverID, string(order.Status()),
		order.PrepStartedAt(), order.PickedUpAt(), order.DeliveredAt(),
		order.CancelReason(), order.Notes(), order.CashbackEarned(),
		string(order.PaymentStatus()), order.ETAMinutes(), time.Now().UTC(),
	)
	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $2, updated_at = $3 WHERE id = $1`,
		id, string(status), time.Now().UTC())
	return err
}

func (r *OrderRepository) findItemsByOrder(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	query := `SELECT id, order_id, menu_item_id, name, quantity, unit_price, modifiers, notes, line_total
		FROM order_items WHERE order_id = $1 ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var orderID, menuItemID uuid.UUID
		var modifiers string
		err := rows.Scan(
			&item, &orderID, &menuItemID,
			&item, &item, &item, &modifiers, &item, &item,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// ============ Helpers ============

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanOrder(s scanner) (*domain.Order, error) {
	var (
		id              uuid.UUID
		orderNumber     string
		customerID      uuid.UUID
		restaurantID    uuid.UUID
		driverID        *uuid.UUID
		status          string
		subtotal        float64
		deliveryFee     float64
		serviceFee      float64
		vat             float64
		discount        float64
		total           float64
		paymentMethod   string
		paymentStatus   string
		deliveryAddress string
		latitude        float64
		longitude       float64
		etaMinutes      int
		prepStartedAt   *time.Time
		pickedUpAt      *time.Time
		deliveredAt     *time.Time
		cancelReason    string
		notes           string
		cashbackEarned  float64
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := s.Scan(
		&id, &orderNumber, &customerID, &restaurantID, &driverID, &status,
		&subtotal, &deliveryFee, &serviceFee, &vat, &discount, &total,
		&paymentMethod, &paymentStatus, &deliveryAddress, &latitude, &longitude,
		&etaMinutes, &prepStartedAt, &pickedUpAt, &deliveredAt,
		&cancelReason, &notes, &cashbackEarned, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("scanOrder: %w", err)
	}

	return &domain.Order{
		id:              id,
		orderNumber:     orderNumber,
		customerID:      customerID,
		restaurantID:    restaurantID,
		driverID:        driverID,
		status:          domain.OrderStatus(status),
		subtotal:        subtotal,
		deliveryFee:     deliveryFee,
		serviceFee:      serviceFee,
		vat:             vat,
		discount:        discount,
		total:           total,
		paymentMethod:   domain.PaymentMethod(paymentMethod),
		paymentStatus:   domain.PaymentStatus(paymentStatus),
		deliveryAddress: deliveryAddress,
		latitude:        latitude,
		longitude:       longitude,
		etaMinutes:      etaMinutes,
		prepStartedAt:   prepStartedAt,
		pickedUpAt:      pickedUpAt,
		deliveredAt:     deliveredAt,
		cancelReason:    cancelReason,
		notes:           notes,
		cashbackEarned:  cashbackEarned,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}, nil
}

func scanOrders(rows *sql.Rows) ([]*domain.Order, error) {
	var orders []*domain.Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
