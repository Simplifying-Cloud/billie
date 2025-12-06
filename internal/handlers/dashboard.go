package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) Dashboard(c echo.Context) error {
	ctx := c.Request().Context()

	// Get metrics
	metrics, err := h.Services.Metrics.GetDashboardMetrics(ctx)
	if err != nil {
		return err
	}

	// Get chart data
	chartData, err := h.Services.Metrics.GetChartData(ctx)
	if err != nil {
		return err
	}

	// Get recent invoices
	recentInvoices, err := h.Services.Invoices.GetRecent(ctx, 3)
	if err != nil {
		return err
	}

	// Get recent clients (all clients, limited to 3)
	clients, err := h.Services.Clients.List(ctx)
	if err != nil {
		return err
	}
	recentClients := clients
	if len(recentClients) > 3 {
		recentClients = recentClients[:3]
	}

	// Get recent expenses
	recentExpenses, err := h.Services.Expenses.GetRecent(ctx, 3)
	if err != nil {
		return err
	}

	data := pages.DashboardData{
		Metrics:        ConvertMetricsToMock(metrics),
		RecentInvoices: ConvertInvoicesToMock(recentInvoices),
		RecentClients:  ConvertClientsToMock(recentClients),
		RecentExpenses: ConvertExpensesToMock(recentExpenses),
		ChartData:      ConvertChartDataToMock(chartData),
	}

	return render.Render(c, http.StatusOK, pages.Dashboard(data))
}
