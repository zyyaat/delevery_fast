package application

import (
	"context"
	"math/rand"
	"time"

	"github.com/food-platform/analytics/internal/domain"
)

type GetDashboardStatsUseCase struct{}

func NewGetDashboardStatsUseCase() *GetDashboardStatsUseCase {
	return &GetDashboardStatsUseCase{}
}

func (uc *GetDashboardStatsUseCase) Execute(ctx context.Context) *domain.DashboardStats {
	// In production, query ClickHouse for real metrics
	// For MVP, return mock data
	return &domain.DashboardStats{
		OrdersActive:      423,
		DriversOnline:     187,
		RestaurantsActive: 142,
		AlertsActive:      3,
		GMVToday:          84250,
		OrdersPerMin:      7.2,
		AvgDeliveryTime:   32,
		DriverUtilization: 0.78,
		OrderCompletion:   0.94,
		CancellationRate:  0.04,
		PaymentSuccess:    0.985,
		AvgDriverRating:   4.7,
		P95Latency:        380,
	}
}

type GetZoneMetricsUseCase struct{}

func NewGetZoneMetricsUseCase() *GetZoneMetricsUseCase {
	return &GetZoneMetricsUseCase{}
}

func (uc *GetZoneMetricsUseCase) Execute(ctx context.Context) []domain.ZoneMetric {
	return []domain.ZoneMetric{
		{ZoneID: "maadi", ZoneName: "معادي", Orders: 72, Drivers: 31, Gap: -8, AvgTime: 28, DemandLevel: "high"},
		{ZoneID: "zamalek", ZoneName: "الزمالك", Orders: 58, Drivers: 25, Gap: -5, AvgTime: 30, DemandLevel: "high"},
		{ZoneID: "nasr", ZoneName: "مدينة نصر", Orders: 45, Drivers: 38, Gap: 3, AvgTime: 35, DemandLevel: "medium"},
		{ZoneID: "downtown", ZoneName: "وسط البلد", Orders: 32, Drivers: 28, Gap: 4, AvgTime: 32, DemandLevel: "medium"},
		{ZoneID: "heliopolis", ZoneName: "مصر الجديدة", Orders: 28, Drivers: 30, Gap: 2, AvgTime: 38, DemandLevel: "low"},
		{ZoneID: "tagamoa", ZoneName: "التجمع", Orders: 22, Drivers: 18, Gap: -4, AvgTime: 42, DemandLevel: "medium"},
	}
}

type GetIncidentsUseCase struct{}

func NewGetIncidentsUseCase() *GetIncidentsUseCase {
	return &GetIncidentsUseCase{}
}

func (uc *GetIncidentsUseCase) Execute(ctx context.Context) []domain.Incident {
	return []domain.Incident{
		{ID: "INC-001", Severity: "P0", Title: "Vodafone Cash API Down", Description: "Payment success rate dropped to 65%", Status: "mitigating", StartedAt: time.Now().Add(-15 * time.Minute), Duration: 15 * time.Minute, Owner: "Ahmed K."},
		{ID: "INC-002", Severity: "P1", Title: "معادي - Driver Shortage", Description: "Gap: 18 orders without drivers", Status: "investigating", StartedAt: time.Now().Add(-25 * time.Minute), Duration: 25 * time.Minute, Owner: "Omar T."},
		{ID: "INC-003", Severity: "P2", Title: "Pizza Hut Maadi - 5 delayed orders", Description: "Orders delayed >15 min", Status: "acknowledged", StartedAt: time.Now().Add(-45 * time.Minute), Duration: 45 * time.Minute, Owner: "Sarah M."},
	}
}

type GetForecastUseCase struct{}

func NewGetForecastUseCase() *GetForecastUseCase {
	return &GetForecastUseCase{}
}

func (uc *GetForecastUseCase) Execute(ctx context.Context) []domain.ForecastPoint {
	var points []domain.ForecastPoint
	hours := []string{"6am", "7am", "8am", "9am", "10am", "11am", "12pm", "1pm", "2pm", "3pm", "4pm", "5pm", "6pm", "7pm", "8pm", "9pm", "10pm", "11pm"}
	forecast := []int{80, 120, 160, 220, 280, 340, 480, 440, 320, 280, 250, 300, 400, 450, 420, 350, 280, 200}

	for i, h := range hours {
		points = append(points, domain.ForecastPoint{
			Hour:       i + 6,
			Day:        h,
			Forecast:   forecast[i],
			Confidence: 0.85 + rand.Float64()*0.1,
		})
	}
	return points
}
