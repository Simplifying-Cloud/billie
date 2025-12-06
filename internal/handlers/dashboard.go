package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func Dashboard(c echo.Context) error {
	data := pages.DashboardData{
		Metrics:        mock.Metrics,
		RecentInvoices: mock.GetRecentInvoices(3),
		RecentClients:  mock.Clients[:3],
		RecentExpenses: mock.GetRecentExpenses(3),
		ChartData:      mock.RevenueChartData,
	}
	return render.Render(c, http.StatusOK, pages.Dashboard(data))
}
