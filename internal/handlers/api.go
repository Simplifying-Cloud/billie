package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dscott/invoicey/internal/services"
	"github.com/labstack/echo/v4"
)

// APIResponse is the standard JSON API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// === Clients API ===

func (h *Handlers) APIListClients(c echo.Context) error {
	ctx := c.Request().Context()

	clients, err := h.Services.Clients.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertClientsToMock(clients)})
}

func (h *Handlers) APIGetClient(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	client, err := h.Services.Clients.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "Client not found"})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertClientToMock(*client)})
}

func (h *Handlers) APICreateClient(c echo.Context) error {
	ctx := c.Request().Context()

	var params services.CreateClientParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	client, err := h.Services.Clients.Create(ctx, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, APIResponse{Success: true, Data: ConvertClientToMock(*client)})
}

func (h *Handlers) APIUpdateClient(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var params services.UpdateClientParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	client, err := h.Services.Clients.Update(ctx, id, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertClientToMock(*client)})
}

func (h *Handlers) APIDeleteClient(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Clients.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true})
}

// === Invoices API ===

func (h *Handlers) APIListInvoices(c echo.Context) error {
	ctx := c.Request().Context()
	status := c.QueryParam("status")

	var invoices []services.InvoiceWithDetails
	var err error

	if status == "" || status == "all" {
		invoices, err = h.Services.Invoices.List(ctx)
	} else {
		invoices, err = h.Services.Invoices.ListByStatus(ctx, status)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertInvoicesToMock(invoices)})
}

func (h *Handlers) APIGetInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "Invoice not found"})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertInvoiceToMock(*invoice)})
}

// CreateInvoiceRequest represents the JSON request for creating an invoice
type CreateInvoiceRequest struct {
	ClientID  string                    `json:"client_id"`
	IssueDate string                    `json:"issue_date"`
	DueDate   string                    `json:"due_date"`
	TaxRate   float64                   `json:"tax_rate"`
	Notes     string                    `json:"notes"`
	LineItems []services.LineItemInput  `json:"line_items"`
}

func (h *Handlers) APICreateInvoice(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	issueDate, _ := time.Parse("2006-01-02", req.IssueDate)
	dueDate, _ := time.Parse("2006-01-02", req.DueDate)

	params := services.CreateInvoiceParams{
		ClientID:  req.ClientID,
		IssueDate: issueDate,
		DueDate:   dueDate,
		TaxRate:   req.TaxRate,
		Notes:     req.Notes,
		LineItems: req.LineItems,
	}

	invoice, err := h.Services.Invoices.Create(ctx, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, APIResponse{Success: true, Data: ConvertInvoiceToMock(*invoice)})
}

// UpdateInvoiceRequest represents the JSON request for updating an invoice
type UpdateInvoiceRequest struct {
	ClientID  string                   `json:"client_id"`
	IssueDate string                   `json:"issue_date"`
	DueDate   string                   `json:"due_date"`
	TaxRate   float64                  `json:"tax_rate"`
	Notes     string                   `json:"notes"`
	LineItems []services.LineItemInput `json:"line_items"`
}

func (h *Handlers) APIUpdateInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var req UpdateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	issueDate, _ := time.Parse("2006-01-02", req.IssueDate)
	dueDate, _ := time.Parse("2006-01-02", req.DueDate)

	params := services.UpdateInvoiceParams{
		ClientID:  req.ClientID,
		IssueDate: issueDate,
		DueDate:   dueDate,
		TaxRate:   req.TaxRate,
		Notes:     req.Notes,
		LineItems: req.LineItems,
	}

	invoice, err := h.Services.Invoices.Update(ctx, id, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertInvoiceToMock(*invoice)})
}

func (h *Handlers) APIDeleteInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Invoices.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true})
}

func (h *Handlers) APISendInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.Send(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertInvoiceToMock(*invoice)})
}

func (h *Handlers) APIMarkPaid(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.MarkAsPaid(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertInvoiceToMock(*invoice)})
}

// === Expenses API ===

func (h *Handlers) APIListExpenses(c echo.Context) error {
	ctx := c.Request().Context()
	category := c.QueryParam("category")

	var expenses []services.ExpenseWithDetails
	var err error

	if category == "" || category == "all" {
		expenses, err = h.Services.Expenses.List(ctx)
	} else {
		expenses, err = h.Services.Expenses.ListByCategory(ctx, category)
	}

	if err != nil {
		expenses = []services.ExpenseWithDetails{}
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertExpensesToMock(expenses)})
}

func (h *Handlers) APIGetExpense(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	expense, err := h.Services.Expenses.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "Expense not found"})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertExpenseToMock(*expense)})
}

// CreateExpenseRequest represents the JSON request for creating an expense
type CreateExpenseRequest struct {
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Merchant    string  `json:"merchant"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	ReceiptPath string  `json:"receipt_path"`
	Notes       string  `json:"notes"`
}

func (h *Handlers) APICreateExpense(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	// Handle merchant - get or create
	var merchantID string
	if req.Merchant != "" {
		merchant, err := h.Services.Expenses.GetOrCreateMerchant(ctx, req.Merchant)
		if err == nil {
			merchantID = merchant.ID
		}
	}

	params := services.CreateExpenseParams{
		Description: req.Description,
		CategoryID:  req.CategoryID,
		MerchantID:  merchantID,
		Amount:      req.Amount,
		Date:        date,
		ReceiptPath: req.ReceiptPath,
		Notes:       req.Notes,
	}

	expense, err := h.Services.Expenses.Create(ctx, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, APIResponse{Success: true, Data: ConvertExpenseToMock(*expense)})
}

func (h *Handlers) APIUpdateExpense(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var req CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	// Handle merchant - get or create
	var merchantID string
	if req.Merchant != "" {
		merchant, err := h.Services.Expenses.GetOrCreateMerchant(ctx, req.Merchant)
		if err == nil {
			merchantID = merchant.ID
		}
	}

	params := services.UpdateExpenseParams{
		Description: req.Description,
		CategoryID:  req.CategoryID,
		MerchantID:  merchantID,
		Amount:      req.Amount,
		Date:        date,
		ReceiptPath: req.ReceiptPath,
		Notes:       req.Notes,
	}

	expense, err := h.Services.Expenses.Update(ctx, id, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: ConvertExpenseToMock(*expense)})
}

func (h *Handlers) APIDeleteExpense(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Expenses.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true})
}

// === Metrics API ===

func (h *Handlers) APIDashboardMetrics(c echo.Context) error {
	ctx := c.Request().Context()

	metrics, err := h.Services.Metrics.GetDashboardMetrics(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
	}

	chartData, err := h.Services.Metrics.GetChartData(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
	}

	data := map[string]interface{}{
		"metrics":   ConvertMetricsToMock(metrics),
		"chartData": ConvertChartDataToMock(chartData),
	}

	return c.JSON(http.StatusOK, APIResponse{Success: true, Data: data})
}

// Ensure json import is used
var _ = json.Unmarshal
