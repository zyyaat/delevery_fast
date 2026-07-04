package domain

import "time"

// Metric represents a single data point.
type Metric struct {
	Name      string
	Value     float64
	Timestamp time.Time
	Tags      map[string]string
}

// MetricSeries represents a time series of metrics.
type MetricSeries struct {
	Name   string
	Points []MetricPoint
}

type MetricPoint struct {
	Timestamp time.Time
	Value     float64
}

// DashboardStats represents real-time platform metrics.
type DashboardStats struct {
	OrdersActive     int
	DriversOnline    int
	RestaurantsActive int
	AlertsActive     int
	GMVToday         float64
	OrdersPerMin     float64
	AvgDeliveryTime  float64
	DriverUtilization float64
	OrderCompletion  float64
	CancellationRate float64
	PaymentSuccess   float64
	AvgDriverRating  float64
	P95Latency       float64
}

// ZoneMetric represents metrics for a specific zone.
type ZoneMetric struct {
	ZoneID     string
	ZoneName   string
	Orders     int
	Drivers    int
	Gap        int // orders - drivers (negative = need more drivers)
	AvgTime    float64
	DemandLevel string // low, medium, high
}

// Incident represents a platform incident.
type Incident struct {
	ID          string
	Severity    string // P0, P1, P2, P3
	Title       string
	Description string
	Status      string // detected, acknowledged, investigating, resolved
	StartedAt   time.Time
	Duration    time.Duration
	Owner       string
}

// ForecastPoint represents a predicted demand point.
type ForecastPoint struct {
	Hour     int
	Day      string
	Forecast int
	Confidence float64
}
