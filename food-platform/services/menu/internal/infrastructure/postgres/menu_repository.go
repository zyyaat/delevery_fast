// Package postgres implements the Menu Service repositories.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/menu/internal/domain"
	"github.com/google/uuid"
)

// MenuItemRepository implements application.MenuItemRepository.
type MenuItemRepository struct {
	db *sql.DB
}

func NewMenuItemRepository(db *sql.DB) *MenuItemRepository {
	return &MenuItemRepository{db: db}
}

func (r *MenuItemRepository) Create(ctx context.Context, item *domain.MenuItem) error {
	query := `
		INSERT INTO menu_items (
			id, restaurant_id, category_id, name, description, price,
			image_url, is_available, prep_time_minutes, rating, rating_count,
			is_most_ordered, display_order, modifiers, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID(), item.RestaurantID(), item.CategoryID(),
		item.Name(), item.Description(), item.Price(),
		nullableString(item.ImageURL()), item.IsAvailable(), item.PrepTimeMin(),
		item.Rating(), item.RatingCount(), item.IsMostOrdered(),
		item.DisplayOrder(), modifiersToJSON(item.Modifiers()),
		item.CreatedAt(), item.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("menu_item_repo.Create: %w", err)
	}
	return nil
}

func (r *MenuItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, category_id, name, COALESCE(description, ''), price,
		       COALESCE(image_url, ''), is_available, prep_time_minutes,
		       rating, rating_count, is_most_ordered, display_order,
		       COALESCE(modifiers, '[]'::jsonb), created_at, updated_at
		FROM menu_items
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanMenuItem(row)
}

func (r *MenuItemRepository) FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, category_id, name, COALESCE(description, ''), price,
		       COALESCE(image_url, ''), is_available, prep_time_minutes,
		       rating, rating_count, is_most_ordered, display_order,
		       COALESCE(modifiers, '[]'::jsonb), created_at, updated_at
		FROM menu_items
		WHERE restaurant_id = $1
		ORDER BY category_id, display_order
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("menu_item_repo.FindByRestaurant: %w", err)
	}
	defer rows.Close()
	return scanMenuItems(rows)
}

func (r *MenuItemRepository) FindByCategory(ctx context.Context, restaurantID, categoryID uuid.UUID) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, category_id, name, COALESCE(description, ''), price,
		       COALESCE(image_url, ''), is_available, prep_time_minutes,
		       rating, rating_count, is_most_ordered, display_order,
		       COALESCE(modifiers, '[]'::jsonb), created_at, updated_at
		FROM menu_items
		WHERE restaurant_id = $1 AND category_id = $2
		ORDER BY display_order
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("menu_item_repo.FindByCategory: %w", err)
	}
	defer rows.Close()
	return scanMenuItems(rows)
}

func (r *MenuItemRepository) Update(ctx context.Context, item *domain.MenuItem) error {
	query := `
		UPDATE menu_items SET
			name = $2, description = $3, price = $4,
			image_url = $5, is_available = $6, prep_time_minutes = $7,
			is_most_ordered = $8, display_order = $9, modifiers = $10,
			updated_at = $11
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID(), item.Name(), item.Description(), item.Price(),
		nullableString(item.ImageURL()), item.IsAvailable(), item.PrepTimeMin(),
		item.IsMostOrdered(), item.DisplayOrder(), modifiersToJSON(item.Modifiers()),
		time.Now().UTC(),
	)
	return err
}

func (r *MenuItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu_items WHERE id = $1`, id)
	return err
}

func (r *MenuItemRepository) SetAvailability(ctx context.Context, id uuid.UUID, available bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu_items SET is_available = $2, updated_at = $3 WHERE id = $1`,
		id, available, time.Now().UTC())
	return err
}

// ============ Category Repository ============

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, cat *domain.Category) error {
	query := `
		INSERT INTO menu_categories (id, restaurant_id, name, display_order)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, cat.ID(), cat.RestaurantID(), cat.Name(), cat.DisplayOrder())
	return err
}

func (r *CategoryRepository) FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Category, error) {
	query := `
		SELECT id, restaurant_id, name, display_order
		FROM menu_categories
		WHERE restaurant_id = $1
		ORDER BY display_order
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		cats = append(cats, &c)
	}
	return cats, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, cat *domain.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE menu_categories SET name = $2, display_order = $3 WHERE id = $1`,
		cat.ID(), cat.Name(), cat.DisplayOrder())
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu_categories WHERE id = $1`, id)
	return err
}

// ============ Helpers ============

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func modifiersToJSON(mods []domain.Modifier) interface{} {
	if len(mods) == 0 {
		return nil
	}
	// In production, marshal to JSON properly
	// For now, return nil and handle in application layer
	return nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanMenuItem(s scanner) (*domain.MenuItem, error) {
	var (
		id            uuid.UUID
		restaurantID  uuid.UUID
		categoryID    uuid.UUID
		name          string
		description   string
		price         float64
		imageURL      string
		isAvailable   bool
		prepTimeMin   int
		rating        float64
		ratingCount   int
		isMostOrdered bool
		displayOrder  int
		modifiersJSON []byte
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := s.Scan(
		&id, &restaurantID, &categoryID, &name, &description, &price,
		&imageURL, &isAvailable, &prepTimeMin,
		&rating, &ratingCount, &isMostOrdered, &displayOrder,
		&modifiersJSON, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrMenuItemNotFound
		}
		return nil, fmt.Errorf("scanMenuItem: %w", err)
	}

	return &domain.MenuItem{
		id:            id,
		restaurantID:  restaurantID,
		categoryID:    categoryID,
		name:          name,
		description:   description,
		price:         price,
		imageURL:      imageURL,
		isAvailable:   isAvailable,
		prepTimeMin:   prepTimeMin,
		rating:        rating,
		ratingCount:   ratingCount,
		isMostOrdered: isMostOrdered,
		displayOrder:  displayOrder,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}, nil
}

func scanMenuItems(rows *sql.Rows) ([]*domain.MenuItem, error) {
	var items []*domain.MenuItem
	for rows.Next() {
		item, err := scanMenuItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
