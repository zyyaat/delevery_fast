package domain

import (
	"testing"
)

func TestNewRestaurant_Valid(t *testing.T) {
	r, err := NewRestaurant(
		"Pizza Hut Maadi",
		30.0444,
		31.2357,
		"15 Road 9, Maadi",
		"Cairo",
		[]CuisineType{CuisineItalian, CuisineFastFood},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if r.Name() != "Pizza Hut Maadi" {
		t.Errorf("expected name, got %s", r.Name())
	}
	if r.Slug() != "pizza-hut-maadi" {
		t.Errorf("expected slug 'pizza-hut-maadi', got %s", r.Slug())
	}
	if len(r.CuisineTypes()) != 2 {
		t.Errorf("expected 2 cuisines, got %d", len(r.CuisineTypes()))
	}
	if r.Status() != RestaurantStatusPendingVerification {
		t.Errorf("expected pending_verification, got %s", r.Status())
	}
	if r.IsOpen() != true {
		t.Error("expected default open")
	}
}

func TestNewRestaurant_EmptyName(t *testing.T) {
	_, err := NewRestaurant("", 30, 31, "addr", "Cairo", nil)
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

func TestNewRestaurant_InvalidCoordinates(t *testing.T) {
	tests := []struct {
		lat, lng float64
	}{
		{91, 31},   // lat too high
		{-91, 31},  // lat too low
		{30, 181},  // lng too high
		{30, -181}, // lng too low
	}

	for _, tt := range tests {
		_, err := NewRestaurant("Test", tt.lat, tt.lng, "addr", "Cairo", nil)
		if err != ErrInvalidCoordinates {
			t.Errorf("expected ErrInvalidCoordinates for (%f, %f), got %v", tt.lat, tt.lng, err)
		}
	}
}

func TestRestaurant_DistanceTo(t *testing.T) {
	r, _ := NewRestaurant("Test", 30.0444, 31.2357, "addr", "Cairo", nil)

	// Same location → 0 distance
	dist := r.DistanceTo(30.0444, 31.2357)
	if dist > 0.01 {
		t.Errorf("expected ~0 distance, got %f", dist)
	}

	// ~1km away (Cairo center to nearby)
	dist = r.DistanceTo(30.0500, 31.2400)
	if dist < 0.5 || dist > 2 {
		t.Errorf("expected ~1km, got %f", dist)
	}
}

func TestRestaurant_IsAcceptingOrders(t *testing.T) {
	r, _ := NewRestaurant("Test", 30, 31, "addr", "Cairo", nil)

	// Pending verification → not accepting
	if r.IsAcceptingOrders() {
		t.Error("expected not accepting when pending")
	}

	// Active + open → accepting
	r.SetStatus(RestaurantStatusActive)
	if !r.IsAcceptingOrders() {
		t.Error("expected accepting when active + open")
	}

	// Active + closed → not accepting
	r.SetOpen(false)
	if r.IsAcceptingOrders() {
		t.Error("expected not accepting when closed")
	}

	// Suspended → not accepting
	r.SetOpen(true)
	r.SetStatus(RestaurantStatusSuspended)
	if r.IsAcceptingOrders() {
		t.Error("expected not accepting when suspended")
	}
}

func TestRestaurant_UpdateRating(t *testing.T) {
	r, _ := NewRestaurant("Test", 30, 31, "addr", "Cairo", nil)

	r.UpdateRating(4.0)
	if r.Rating() != 4.0 {
		t.Errorf("expected 4.0, got %f", r.Rating())
	}
	if r.RatingCount() != 1 {
		t.Errorf("expected 1 count, got %d", r.RatingCount())
	}

	r.UpdateRating(5.0)
	// (4.0 + 5.0) / 2 = 4.5
	if r.Rating() != 4.5 {
		t.Errorf("expected 4.5, got %f", r.Rating())
	}
	if r.RatingCount() != 2 {
		t.Errorf("expected 2 count, got %d", r.RatingCount())
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Pizza Hut", "pizza-hut"},
		{"McDonald's", "mcdonald-s"},
		{"KFC - Egyptian", "kfc-egyptian"},
		{"  Multiple   Spaces  ", "multiple-spaces"},
		{"مطعم المصري", ""}, // Arabic → falls back to UUID
	}

	for _, tt := range tests {
		result := generateSlug(tt.input)
		if tt.expected == "" {
			if result == "" {
				t.Errorf("expected non-empty slug for %q", tt.input)
			}
			continue
		}
		if result != tt.expected {
			t.Errorf("generateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
