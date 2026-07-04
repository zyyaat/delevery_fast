// Package application contains use cases for the Menu Service.
package application

import (
	"context"
	"fmt"

	"github.com/food-platform/services/menu/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type MenuItemRepository interface {
	Create(ctx context.Context, item *domain.MenuItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.MenuItem, error)
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.MenuItem, error)
	FindByCategory(ctx context.Context, restaurantID, categoryID uuid.UUID) ([]*domain.MenuItem, error)
	Update(ctx context.Context, item *domain.MenuItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetAvailability(ctx context.Context, id uuid.UUID, available bool) error
}

type CategoryRepository interface {
	Create(ctx context.Context, cat *domain.Category) error
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Category, error)
	Update(ctx context.Context, cat *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ============ DTOs ============

type ModifierOptionDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceDelta float64   `json:"price_delta"`
}

type ModifierDTO struct {
	ID             uuid.UUID           `json:"id"`
	Name           string              `json:"name"`
	Required       bool                `json:"required"`
	MultipleChoice bool                `json:"multiple_choice"`
	Options        []ModifierOptionDTO `json:"options"`
}

type MenuItemDTO struct {
	ID            uuid.UUID     `json:"id"`
	RestaurantID  uuid.UUID     `json:"restaurant_id"`
	CategoryID    uuid.UUID     `json:"category_id"`
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	Price         float64       `json:"price"`
	ImageURL      string        `json:"image_url,omitempty"`
	IsAvailable   bool          `json:"is_available"`
	PrepTimeMin   int           `json:"prep_time_minutes"`
	Rating        float64       `json:"rating,omitempty"`
	RatingCount   int           `json:"rating_count,omitempty"`
	IsMostOrdered bool          `json:"is_most_ordered,omitempty"`
	Modifiers     []ModifierDTO `json:"modifiers,omitempty"`
}

type CategoryDTO struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
}

type CategoryWithItemsDTO struct {
	Category CategoryDTO     `json:"category"`
	Items    []MenuItemDTO   `json:"items"`
}

// ============ Commands ============

type CreateMenuItemCommand struct {
	RestaurantID uuid.UUID
	CategoryID   uuid.UUID
	Name         string
	Description  string
	Price        float64
	ImageURL     string
	PrepTimeMin  int
	Modifiers    []CreateModifierCommand
}

type CreateModifierCommand struct {
	Name           string
	Required       bool
	MultipleChoice bool
	Options        []CreateModifierOptionCommand
}

type CreateModifierOptionCommand struct {
	Name       string
	PriceDelta float64
}

type UpdateMenuItemCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Price       *float64
	ImageURL    *string
	PrepTimeMin *int
	IsAvailable *bool
}

type CreateCategoryCommand struct {
	RestaurantID uuid.UUID
	Name         string
	DisplayOrder int
}

// ============ Use Cases ============

type CreateMenuItemUseCase struct {
	itemRepo MenuItemRepository
}

func NewCreateMenuItemUseCase(repo MenuItemRepository) *CreateMenuItemUseCase {
	return &CreateMenuItemUseCase{itemRepo: repo}
}

func (uc *CreateMenuItemUseCase) Execute(ctx context.Context, cmd CreateMenuItemCommand) (*MenuItemDTO, error) {
	item, err := domain.NewMenuItem(cmd.RestaurantID, cmd.CategoryID, cmd.Name, cmd.Price)
	if err != nil {
		return nil, err
	}

	item.SetDescription(cmd.Description)
	if cmd.ImageURL != "" {
		item.SetImageURL(cmd.ImageURL)
	}
	if cmd.PrepTimeMin > 0 {
		item.SetPrepTime(cmd.PrepTimeMin)
	}

	// Convert modifier commands to domain modifiers
	if len(cmd.Modifiers) > 0 {
		mods := make([]domain.Modifier, len(cmd.Modifiers))
		for i, mCmd := range cmd.Modifiers {
			options := make([]domain.ModifierOption, len(mCmd.Options))
			for j, oCmd := range mCmd.Options {
				options[j] = domain.NewModifierOption(oCmd.Name, oCmd.PriceDelta)
			}
			mods[i] = domain.NewModifier(mCmd.Name, mCmd.Required, mCmd.MultipleChoice, options)
		}
		item.SetModifiers(mods)
	}

	if err := uc.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create menu item: %w", err)
	}

	return toItemDTO(item), nil
}

type GetMenuItemUseCase struct {
	itemRepo MenuItemRepository
}

func NewGetMenuItemUseCase(repo MenuItemRepository) *GetMenuItemUseCase {
	return &GetMenuItemUseCase{itemRepo: repo}
}

func (uc *GetMenuItemUseCase) Execute(ctx context.Context, id uuid.UUID) (*MenuItemDTO, error) {
	item, err := uc.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toItemDTO(item), nil
}

type GetMenuUseCase struct {
	itemRepo     MenuItemRepository
	categoryRepo CategoryRepository
}

func NewGetMenuUseCase(itemRepo MenuItemRepository, categoryRepo CategoryRepository) *GetMenuUseCase {
	return &GetMenuUseCase{itemRepo: itemRepo, categoryRepo: categoryRepo}
}

func (uc *GetMenuUseCase) Execute(ctx context.Context, restaurantID uuid.UUID) ([]CategoryWithItemsDTO, error) {
	categories, err := uc.categoryRepo.FindByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	items, err := uc.itemRepo.FindByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	// Group items by category
	itemsByCategory := make(map[uuid.UUID][]*domain.MenuItem)
	for _, item := range items {
		itemsByCategory[item.CategoryID()] = append(itemsByCategory[item.CategoryID()], item)
	}

	result := make([]CategoryWithItemsDTO, 0, len(categories))
	for _, cat := range categories {
		catItems := itemsByCategory[cat.ID()]
		itemDTOs := make([]MenuItemDTO, 0, len(catItems))
		for _, item := range catItems {
			itemDTOs = append(itemDTOs, *toItemDTO(item))
		}
		result = append(result, CategoryWithItemsDTO{
			Category: CategoryDTO{
				ID:           cat.ID(),
				RestaurantID: cat.RestaurantID(),
				Name:         cat.Name(),
				DisplayOrder: cat.DisplayOrder(),
			},
			Items: itemDTOs,
		})
	}

	return result, nil
}

type UpdateMenuItemUseCase struct {
	itemRepo MenuItemRepository
}

func NewUpdateMenuItemUseCase(repo MenuItemRepository) *UpdateMenuItemUseCase {
	return &UpdateMenuItemUseCase{itemRepo: repo}
}

func (uc *UpdateMenuItemUseCase) Execute(ctx context.Context, cmd UpdateMenuItemCommand) (*MenuItemDTO, error) {
	item, err := uc.itemRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if cmd.Name != nil {
		if err := item.SetName(*cmd.Name); err != nil {
			return nil, err
		}
	}
	if cmd.Description != nil {
		item.SetDescription(*cmd.Description)
	}
	if cmd.Price != nil {
		if err := item.SetPrice(*cmd.Price); err != nil {
			return nil, err
		}
	}
	if cmd.ImageURL != nil {
		item.SetImageURL(*cmd.ImageURL)
	}
	if cmd.PrepTimeMin != nil {
		item.SetPrepTime(*cmd.PrepTimeMin)
	}
	if cmd.IsAvailable != nil {
		item.SetAvailable(*cmd.IsAvailable)
	}

	if err := uc.itemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("update menu item: %w", err)
	}

	return toItemDTO(item), nil
}

type ToggleAvailabilityUseCase struct {
	itemRepo MenuItemRepository
}

func NewToggleAvailabilityUseCase(repo MenuItemRepository) *ToggleAvailabilityUseCase {
	return &ToggleAvailabilityUseCase{itemRepo: repo}
}

func (uc *ToggleAvailabilityUseCase) Execute(ctx context.Context, id uuid.UUID, available bool) error {
	return uc.itemRepo.SetAvailability(ctx, id, available)
}

type CreateCategoryUseCase struct {
	repo CategoryRepository
}

func NewCreateCategoryUseCase(repo CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{repo: repo}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, cmd CreateCategoryCommand) (*CategoryDTO, error) {
	cat := domain.NewCategory(cmd.RestaurantID, cmd.Name, cmd.DisplayOrder)
	if err := uc.repo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return &CategoryDTO{
		ID:           cat.ID(),
		RestaurantID: cat.RestaurantID(),
		Name:         cat.Name(),
		DisplayOrder: cat.DisplayOrder(),
	}, nil
}

// ============ Helpers ============

func toItemDTO(item *domain.MenuItem) *MenuItemDTO {
	dto := &MenuItemDTO{
		ID:            item.ID(),
		RestaurantID:  item.RestaurantID(),
		CategoryID:    item.CategoryID(),
		Name:          item.Name(),
		Description:   item.Description(),
		Price:         item.Price(),
		ImageURL:      item.ImageURL(),
		IsAvailable:   item.IsAvailable(),
		PrepTimeMin:   item.PrepTimeMin(),
		Rating:        item.Rating(),
		RatingCount:   item.RatingCount(),
		IsMostOrdered: item.IsMostOrdered(),
	}

	if len(item.Modifiers()) > 0 {
		dto.Modifiers = make([]ModifierDTO, len(item.Modifiers()))
		for i, m := range item.Modifiers() {
			opts := make([]ModifierOptionDTO, len(m.Options()))
			for j, o := range m.Options() {
				opts[j] = ModifierOptionDTO{
					ID:         o.ID,
					Name:       o.Name,
					PriceDelta: o.PriceDelta,
				}
			}
			dto.Modifiers[i] = ModifierDTO{
				ID:             m.ID,
				Name:           m.Name,
				Required:       m.Required,
				MultipleChoice: m.MultipleChoice,
				Options:        opts,
			}
		}
	}

	return dto
}
