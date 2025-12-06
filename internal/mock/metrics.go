package mock

import "time"

type DashboardMetrics struct {
	TotalRevenue       float64 // Calendar year total
	RevenueChange      float64
	OutstandingAmount  float64 // Unpaid invoices
	OutstandingChange  float64
	TotalExpenses      float64 // Calendar year total
	ExpensesChange     float64
	UpcomingExpenses   float64 // Next 30 days
	UpcomingChange     float64
}

type ChartDataPoint struct {
	Label string
	Value float64
}

type ChartData struct {
	Labels   []string
	Revenue  []float64
	Expenses []float64
}

var Metrics = DashboardMetrics{
	TotalRevenue:      91308.38,
	RevenueChange:     12.5,
	OutstandingAmount: 35181.25,
	OutstandingChange: -5.2,
	TotalExpenses:     1449.80,
	ExpensesChange:    8.3,
	UpcomingExpenses:  850.00,
	UpcomingChange:    15.2,
}

var RevenueChartData = ChartData{
	Labels:   []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	Revenue:  []float64{8500, 12200, 9800, 15100, 11500, 14200, 12500, 18200, 15800, 22100, 19500, 23200},
	Expenses: []float64{1200, 980, 1450, 1100, 1350, 1600, 2100, 1800, 2400, 2200, 2800, 1450},
}

var StatusDistribution = []ChartDataPoint{
	{Label: "Paid", Value: 2},
	{Label: "Sent", Value: 2},
	{Label: "Draft", Value: 2},
	{Label: "Overdue", Value: 1},
}

// GetUpcomingExpenses returns expenses due in the next 30 days
func GetUpcomingExpenses() float64 {
	now := time.Now()
	thirtyDaysFromNow := now.AddDate(0, 0, 30)
	var total float64
	for _, exp := range Expenses {
		if exp.Date.After(now) && exp.Date.Before(thirtyDaysFromNow) {
			total += exp.Amount
		}
	}
	// Return mock value if no upcoming expenses found
	if total == 0 {
		return 850.00
	}
	return total
}
