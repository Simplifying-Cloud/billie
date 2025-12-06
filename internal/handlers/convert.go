package handlers

import (
	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/services"
)

// Convert service types to mock types for template compatibility
// This allows us to keep the templates unchanged while using the database

// ConvertClientToMock converts a service ClientWithStats to mock.Client
func ConvertClientToMock(c services.ClientWithStats) mock.Client {
	return mock.Client{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
		Phone: c.Phone,
		Company: c.Company,
		Address: mock.Address{
			Street:  c.AddressStreet,
			City:    c.AddressCity,
			State:   c.AddressState,
			Zip:     c.AddressZip,
			Country: c.AddressCountry,
		},
		TotalBilled:  c.TotalBilled,
		InvoiceCount: int(c.InvoiceCount),
		CreatedAt:    c.CreatedAt,
	}
}

// ConvertClientsToMock converts a slice of service clients to mock clients
func ConvertClientsToMock(clients []services.ClientWithStats) []mock.Client {
	result := make([]mock.Client, len(clients))
	for i, c := range clients {
		result[i] = ConvertClientToMock(c)
	}
	return result
}

// ConvertInvoiceToMock converts a service InvoiceWithDetails to mock.Invoice
func ConvertInvoiceToMock(inv services.InvoiceWithDetails) mock.Invoice {
	var client *mock.Client
	if inv.Client != nil {
		c := ConvertClientToMock(*inv.Client)
		client = &c
	}

	lineItems := make([]mock.LineItem, len(inv.LineItems))
	for i, li := range inv.LineItems {
		lineItems[i] = mock.LineItem{
			Description: li.Description,
			Quantity:    li.Quantity,
			Rate:        li.Rate,
			Amount:      li.Amount,
		}
	}

	return mock.Invoice{
		ID:          inv.ID,
		Number:      inv.Number,
		ClientID:    inv.ClientID,
		Client:      client,
		Status:      mock.InvoiceStatus(inv.Status),
		IssueDate:   inv.IssueDate,
		DueDate:     inv.DueDate,
		LineItems:   lineItems,
		Subtotal:    inv.Subtotal,
		TaxRate:     inv.TaxRate,
		TaxAmount:   inv.TaxAmount,
		Total:       inv.Total,
		Notes:       inv.Notes,
		PublicToken: inv.PublicToken,
	}
}

// ConvertInvoicesToMock converts a slice of service invoices to mock invoices
func ConvertInvoicesToMock(invoices []services.InvoiceWithDetails) []mock.Invoice {
	result := make([]mock.Invoice, len(invoices))
	for i, inv := range invoices {
		result[i] = ConvertInvoiceToMock(inv)
	}
	return result
}

// ConvertExpenseToMock converts a service ExpenseWithDetails to mock.Expense
func ConvertExpenseToMock(e services.ExpenseWithDetails) mock.Expense {
	return mock.Expense{
		ID:          e.ID,
		Description: e.Description,
		Category:    mock.ExpenseCategory(e.CategorySlug),
		Merchant:    e.MerchantName,
		Amount:      e.Amount,
		Date:        e.Date,
		ReceiptURL:  e.ReceiptPath,
		Notes:       e.Notes,
	}
}

// ConvertExpensesToMock converts a slice of service expenses to mock expenses
func ConvertExpensesToMock(expenses []services.ExpenseWithDetails) []mock.Expense {
	result := make([]mock.Expense, len(expenses))
	for i, e := range expenses {
		result[i] = ConvertExpenseToMock(e)
	}
	return result
}

// ConvertMetricsToMock converts service DashboardMetrics to mock.DashboardMetrics
func ConvertMetricsToMock(m *services.DashboardMetrics) mock.DashboardMetrics {
	return mock.DashboardMetrics{
		TotalRevenue:      m.TotalRevenue,
		RevenueChange:     m.RevenueChange,
		OutstandingAmount: m.OutstandingAmount,
		OutstandingChange: m.OutstandingChange,
		TotalExpenses:     m.TotalExpenses,
		ExpensesChange:    m.ExpensesChange,
		UpcomingExpenses:  m.UpcomingExpenses,
		UpcomingChange:    m.UpcomingChange,
	}
}

// ConvertChartDataToMock converts service ChartData to mock.ChartData
func ConvertChartDataToMock(c *services.ChartData) mock.ChartData {
	return mock.ChartData{
		Labels:   c.Labels,
		Revenue:  c.Revenue,
		Expenses: c.Expenses,
	}
}

// ConvertCategoriesToMock converts service categories to the mock category labels map
func ConvertCategoriesToMock(categories []services.CategoryData) map[mock.ExpenseCategory]string {
	result := make(map[mock.ExpenseCategory]string)
	for _, c := range categories {
		result[mock.ExpenseCategory(c.Slug)] = c.Name
	}
	return result
}
