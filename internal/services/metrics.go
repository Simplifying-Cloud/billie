package services

import (
	"context"
	"time"

	"github.com/dscott/invoicey/internal/database/db"
)

// MetricsService handles dashboard metrics and analytics
type MetricsService struct {
	queries *db.Queries
}

// NewMetricsService creates a new metrics service
func NewMetricsService(queries *db.Queries) *MetricsService {
	return &MetricsService{queries: queries}
}

// DashboardMetrics holds all dashboard metric data
type DashboardMetrics struct {
	TotalRevenue      float64
	RevenueChange     float64
	OutstandingAmount float64
	OutstandingChange float64
	TotalExpenses     float64
	ExpensesChange    float64
	UpcomingExpenses  float64
	UpcomingChange    float64
}

// ChartDataPoint represents a single data point for charts
type ChartDataPoint struct {
	Label    string
	Revenue  float64
	Expenses float64
}

// ChartData holds chart data for visualization
type ChartData struct {
	Labels   []string
	Revenue  []float64
	Expenses []float64
}

// StatusCount holds invoice count by status
type StatusCount struct {
	Status string
	Count  int64
}

// CategoryExpense holds expense totals by category
type CategoryExpense struct {
	CategoryName string
	CategorySlug string
	Total        float64
}

// GetDashboardMetrics returns all dashboard metrics
func (s *MetricsService) GetDashboardMetrics(ctx context.Context) (*DashboardMetrics, error) {
	now := time.Now()
	currentYearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	lastYearStart := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)

	// Get total revenue (all time for paid invoices)
	totalRevenueVal, err := s.queries.GetTotalRevenue(ctx)
	totalRevenue := interfaceToFloat64(totalRevenueVal, err)

	// Get current and previous year revenue for change calculation
	currentYearRevenueVal, _ := s.queries.GetTotalRevenueByYear(ctx, currentYearStart)
	currentYearRevenue := interfaceToFloat64(currentYearRevenueVal, nil)

	lastYearRevenueVal, _ := s.queries.GetTotalRevenueByYear(ctx, lastYearStart)
	lastYearRevenue := interfaceToFloat64(lastYearRevenueVal, nil)

	var revenueChange float64
	if lastYearRevenue > 0 {
		revenueChange = ((currentYearRevenue - lastYearRevenue) / lastYearRevenue) * 100
	}

	// Get outstanding amount
	outstandingVal, _ := s.queries.GetOutstandingAmount(ctx)
	outstandingAmount := interfaceToFloat64(outstandingVal, nil)

	// Get total expenses
	currentYearExpensesVal, _ := s.queries.GetTotalExpensesByYear(ctx, currentYearStart)
	currentYearExpenses := interfaceToFloat64(currentYearExpensesVal, nil)

	lastYearExpensesVal, _ := s.queries.GetTotalExpensesByYear(ctx, lastYearStart)
	lastYearExpenses := interfaceToFloat64(lastYearExpensesVal, nil)

	var expensesChange float64
	if lastYearExpenses > 0 {
		expensesChange = ((currentYearExpenses - lastYearExpenses) / lastYearExpenses) * 100
	}

	// For upcoming expenses, we'll use the current month's expenses as a proxy
	upcomingExpenses := currentYearExpenses / 12 // Monthly average

	return &DashboardMetrics{
		TotalRevenue:      totalRevenue,
		RevenueChange:     revenueChange,
		OutstandingAmount: outstandingAmount,
		OutstandingChange: 0,
		TotalExpenses:     currentYearExpenses,
		ExpensesChange:    expensesChange,
		UpcomingExpenses:  upcomingExpenses,
		UpcomingChange:    0,
	}, nil
}

// GetChartData returns revenue and expense data for charts
func (s *MetricsService) GetChartData(ctx context.Context) (*ChartData, error) {
	now := time.Now()
	currentYearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	// Get monthly revenue
	monthlyRevenue, err := s.queries.GetMonthlyRevenue(ctx, currentYearStart)
	if err != nil {
		return nil, err
	}

	// Get monthly expenses
	monthlyExpenses, err := s.queries.GetMonthlyExpenses(ctx, currentYearStart)
	if err != nil {
		return nil, err
	}

	// Build chart data with all 12 months
	labels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	revenue := make([]float64, 12)
	expenses := make([]float64, 12)

	// Map revenue data to months
	for _, r := range monthlyRevenue {
		monthStr := interfaceToString(r.Month)
		if month := parseMonth(monthStr); month >= 1 && month <= 12 {
			revenue[month-1] = interfaceToFloat64(r.Revenue, nil)
		}
	}

	// Map expense data to months
	for _, e := range monthlyExpenses {
		monthStr := interfaceToString(e.Month)
		if month := parseMonth(monthStr); month >= 1 && month <= 12 {
			expenses[month-1] = interfaceToFloat64(e.Expenses, nil)
		}
	}

	return &ChartData{
		Labels:   labels,
		Revenue:  revenue,
		Expenses: expenses,
	}, nil
}

// parseMonth parses a month string (01-12) to an integer
func parseMonth(s string) int {
	if len(s) == 0 {
		return 0
	}
	var month int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			month = month*10 + int(c-'0')
		}
	}
	return month
}

// GetInvoiceStatusCounts returns invoice counts by status
func (s *MetricsService) GetInvoiceStatusCounts(ctx context.Context) ([]StatusCount, error) {
	counts, err := s.queries.GetInvoiceStatusCounts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]StatusCount, len(counts))
	for i, c := range counts {
		result[i] = StatusCount{
			Status: c.Status,
			Count:  c.Count,
		}
	}
	return result, nil
}

// GetExpensesByCategory returns expense totals grouped by category
func (s *MetricsService) GetExpensesByCategory(ctx context.Context) ([]CategoryExpense, error) {
	expenses, err := s.queries.GetExpensesByCategory(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryExpense, len(expenses))
	for i, e := range expenses {
		result[i] = CategoryExpense{
			CategoryName: e.CategoryName,
			CategorySlug: e.CategorySlug,
			Total:        interfaceToFloat64(e.Total, nil),
		}
	}
	return result, nil
}

// GetUpcomingDueInvoices returns invoices due within 7 days
func (s *MetricsService) GetUpcomingDueInvoices(ctx context.Context) ([]db.Invoice, error) {
	return s.queries.GetUpcomingDueInvoices(ctx)
}
