// Package domain contains the core business logic of the Menu Service.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMenuItemNotFound = errors.New("menu item not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrItemUnavailable  = errors.New("menu item is not available")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrInvalidName      = errors.New("invalid name")
)

// MenuItem represents a food item on a restaurant's menu.
type MenuItem struct {
	id             uuid.UUID
	restaurantID   uuid.UUID
	categoryID     uuid.UUID
	name           string
	description    string
	price          float64
	imageURL       string
	isAvailable    bool
	prepTimeMin    int
	rating         float64
	ratingCount    int
	isMostOrdered  bool
	displayOrder   int
	modifiers      []Modifier
	createdAt      time.Time
	updatedAt      time.Time
}

// Modifier represents a customization option (e.g., "Size", "Extra Cheese").
type Modifier struct {
	id             uuid.UUID
	name           string
	required       bool
	multipleChoice bool
	options        []ModifierOption
}

// ModifierOption represents a choice within a modifier.
type ModifierOption struct {
	id          uuid.UUID
	name        string
	priceDelta  float64
}

// Category represents a menu section (e.g., "Pizzas", "Drinks").
type Category struct {
	id           uuid.UUID
	restaurantID uuid.UUID
	name         string
	displayOrder int
}

// NewMenuItem creates a new menu item with validation.
func NewMenuItem(restaurantID, categoryID uuid.UUID, name string, price float64) (*MenuItem, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if price < 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now().UTC()
	return &MenuItem{
		id:           uuid.New(),
		restaurantID: restaurantID,
		categoryID:   categoryID,
		name:         name,
		price:        price,
		isAvailable:  true,
		prepTimeMin:  10,
		displayOrder: 0,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// NewCategory creates a new menu category.
func NewCategory(restaurantID uuid.UUID, name string, displayOrder int) *Category {
	return &Category{
		id:           uuid.New(),
		restaurantID: restaurantID,
		name:         name,
		displayOrder: displayOrder,
	}
}

// NewModifier creates a new modifier.
func NewModifier(name string, required, multipleChoice bool, options []ModifierOption) Modifier {
	return Modifier{
		id:             uuid.New(),
		name:           name,
		required:       required,
		multipleChoice: multipleChoice,
		options:        options,
	}
}

// NewModifierOption creates a new modifier option.
func NewModifierOption(name string, priceDelta float64) ModifierOption {
	return ModifierOption{
		id:         uuid.New(),
		name:       name,
		priceDelta: priceDelta,
	}
}

// ============ Getters ============

func (m *MenuItem) ID() uuid.UUID         { return m.id }
func (m *MenuItem) RestaurantID() uuid.UUID { return m.restaurantID }
func (m *MenuItem) CategoryID() uuid.UUID { return m.categoryID }
func (m *MenuItem) Name() string          { return m.name }
func (m *MenuItem) Description() string   { return m.description }
func (m *MenuItem) Price() float64        { return m.price }
func (m *MenuItem) ImageURL() string      { return m.imageURL }
func (m *MenuItem) IsAvailable() bool     { return m.isAvailable }
func (m *MenuItem) PrepTimeMin() int      { return m.prepTimeMin }
func (m *MenuItem) Rating() float64       { return m.rating }
func (m *MenuItem) RatingCount() int      { return m.ratingCount }
func (m *MenuItem) IsMostOrdered() bool   { return m.isMostOrdered }
func (m *MenuItem) DisplayOrder() int     { return m.displayOrder }
func (m *MenuItem) Modifiers() []Modifier { return m.modifiers }
func (m *MenuItem) CreatedAt() time.Time  { return m.createdAt }
func (m *MenuItem) UpdatedAt() time.Time  { return m.updatedAt }

func (c *Category) ID() uuid.UUID         { return c.id }
func (c *Category) RestaurantID() uuid.UUID { return c.restaurantID }
func (c *Category) Name() string          { return c.name }
func (c *Category) DisplayOrder() int     { return c.displayOrder }

// ============ Setters ============

func (m *MenuItem) SetName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	m.name = name
	m.updatedAt = time.Now().UTC()
	return nil
}

func (m *MenuItem) SetDescription(desc string) {
	m.description = desc
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetPrice(price float64) error {
	if price < 0 {
		return ErrInvalidPrice
	}
	m.price = price
	m.updatedAt = time.Now().UTC()
	return nil
}

func (m *MenuItem) SetImageURL(url string) {
	m.imageURL = url
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetAvailable(available bool) {
	m.isAvailable = available
	m.updatedAt = time.Now().UTC()
}

// SetUnavailable marks the item as 86'd (sold out).
func (m *MenuItem) SetUnavailable() {
	m.isAvailable = false
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetPrepTime(minutes int) {
	m.prepTimeMin = minutes
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetMostOrdered(val bool) {
	m.isMostOrdered = val
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetDisplayOrder(order int) {
	m.displayOrder = order
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetModifiers(mods []Modifier) {
	m.modifiers = mods
	m.updatedAt = time.Now().UTC()
}

func (m *MenuItem) SetCategoryID(catID uuid.UUID) {
	m.categoryID = catID
	m.updatedAt = time.Now().UTC()
}

// CalculatePrice calculates the total price for a given quantity and selected modifiers.
func (m *MenuItem) CalculatePrice(quantity int, selectedOptions [][]uuid.UUID) float64 {
	if quantity < 1 {
		quantity = 1
	}

	total := m.price * float64(quantity)

	for _, optionIDs := range selectedOptions {
		for _, optID := range optionIDs {
			for _, mod := range m.modifiers {
				for _, opt := range mod.options {
					if opt.id == optID {
						total += opt.priceDelta * float64(quantity)
					}
				}
			}
		}
	}

	return total
}

// Getters for Modifier (needed by application layer)
func (m Modifier) GetID() uuid.UUID { return m.id }
func (m Modifier) GetName() string { return m.name }
func (m Modifier) GetRequired() bool { return m.required }
func (m Modifier) GetMultipleChoice() bool { return m.multipleChoice }
func (m Modifier) GetOptions() []ModifierOption { return m.options }

// Getters for ModifierOption
func (o ModifierOption) GetID() uuid.UUID { return o.id }
func (o ModifierOption) GetName() string { return o.name }
func (o ModifierOption) GetPriceDelta() float64 { return o.priceDelta }

// ReconstructMenuItem creates a MenuItem from persisted data.
func ReconstructMenuItem(
	id, restaurantID, categoryID uuid.UUID, name, description string,
	price float64, imageURL string, isAvailable bool, prepTimeMin int,
	rating float64, ratingCount int, isMostOrdered bool, displayOrder int,
	modifiers []Modifier, createdAt, updatedAt time.Time,
) *MenuItem {
	return &MenuItem{
		id: id, restaurantID: restaurantID, categoryID: categoryID,
		name: name, description: description, price: price,
		imageURL: imageURL, isAvailable: isAvailable, prepTimeMin: prepTimeMin,
		rating: rating, ratingCount: ratingCount, isMostOrdered: isMostOrdered,
		displayOrder: displayOrder, modifiers: modifiers,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

// ReconstructCategory creates a Category from persisted data.
func ReconstructCategory(id, restaurantID uuid.UUID, name string, displayOrder int) *Category {
	return &Category{id: id, restaurantID: restaurantID, name: name, displayOrder: displayOrder}
}
