package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/internal/services"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) InvoicesList(c echo.Context) error {
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
		return err
	}

	mockInvoices := ConvertInvoicesToMock(invoices)

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		return render.Render(c, http.StatusOK, pages.InvoicesTableBody(mockInvoices))
	}

	return render.Render(c, http.StatusOK, pages.InvoicesList(mockInvoices, status))
}

func (h *Handlers) InvoiceDetail(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.Get(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Invoice not found")
	}

	mockInvoice := ConvertInvoiceToMock(*invoice)
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
}

func (h *Handlers) InvoiceEditor(c echo.Context) error {
	// For new invoice, pass nil
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(nil))
}

// CreateInvoice creates a new invoice (POST /invoices)
func (h *Handlers) CreateInvoice(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse dates
	issueDate, _ := time.Parse("2006-01-02", c.FormValue("issue_date"))
	dueDate, _ := time.Parse("2006-01-02", c.FormValue("due_date"))

	// Parse tax rate
	var taxRate float64
	if tr := c.FormValue("tax_rate"); tr != "" {
		json.Unmarshal([]byte(tr), &taxRate)
	}

	// Parse line items from JSON
	var lineItems []services.LineItemInput
	lineItemsJSON := c.FormValue("line_items")
	if lineItemsJSON != "" {
		json.Unmarshal([]byte(lineItemsJSON), &lineItems)
	}

	params := services.CreateInvoiceParams{
		ClientID:  c.FormValue("client_id"),
		IssueDate: issueDate,
		DueDate:   dueDate,
		TaxRate:   taxRate,
		Notes:     c.FormValue("notes"),
		LineItems: lineItems,
	}

	invoice, err := h.Services.Invoices.Create(ctx, params)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to create invoice: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/invoices/"+invoice.ID)
		return c.NoContent(http.StatusOK)
	}

	return c.JSON(http.StatusCreated, ConvertInvoiceToMock(*invoice))
}

// UpdateInvoice updates an existing invoice (PUT /invoices/:id)
func (h *Handlers) UpdateInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	// Parse dates
	issueDate, _ := time.Parse("2006-01-02", c.FormValue("issue_date"))
	dueDate, _ := time.Parse("2006-01-02", c.FormValue("due_date"))

	// Parse tax rate
	var taxRate float64
	if tr := c.FormValue("tax_rate"); tr != "" {
		json.Unmarshal([]byte(tr), &taxRate)
	}

	// Parse line items from JSON
	var lineItems []services.LineItemInput
	lineItemsJSON := c.FormValue("line_items")
	if lineItemsJSON != "" {
		json.Unmarshal([]byte(lineItemsJSON), &lineItems)
	}

	params := services.UpdateInvoiceParams{
		ClientID:  c.FormValue("client_id"),
		IssueDate: issueDate,
		DueDate:   dueDate,
		TaxRate:   taxRate,
		Notes:     c.FormValue("notes"),
		LineItems: lineItems,
	}

	invoice, err := h.Services.Invoices.Update(ctx, id, params)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to update invoice: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		mockInvoice := ConvertInvoiceToMock(*invoice)
		return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
	}

	return c.JSON(http.StatusOK, ConvertInvoiceToMock(*invoice))
}

// DeleteInvoice deletes an invoice (DELETE /invoices/:id)
func (h *Handlers) DeleteInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Invoices.Delete(ctx, id); err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to delete invoice: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/invoices")
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusNoContent)
}

// SendInvoice marks an invoice as sent (POST /invoices/:id/send)
func (h *Handlers) SendInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.Send(ctx, id)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to send invoice: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		mockInvoice := ConvertInvoiceToMock(*invoice)
		return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
	}

	return c.JSON(http.StatusOK, ConvertInvoiceToMock(*invoice))
}

// MarkInvoicePaid marks an invoice as paid (POST /invoices/:id/mark-paid)
func (h *Handlers) MarkInvoicePaid(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	invoice, err := h.Services.Invoices.MarkAsPaid(ctx, id)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to mark invoice as paid: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		mockInvoice := ConvertInvoiceToMock(*invoice)
		return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
	}

	return c.JSON(http.StatusOK, ConvertInvoiceToMock(*invoice))
}

// AddLineItem adds a line item to an invoice (POST /invoices/:id/line-items)
func (h *Handlers) AddLineItem(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var quantity, rate float64
	json.Unmarshal([]byte(c.FormValue("quantity")), &quantity)
	json.Unmarshal([]byte(c.FormValue("rate")), &rate)

	item := services.LineItemInput{
		Description: c.FormValue("description"),
		Quantity:    quantity,
		Rate:        rate,
	}

	invoice, err := h.Services.Invoices.AddLineItem(ctx, id, item)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to add line item: "+err.Error())
	}

	mockInvoice := ConvertInvoiceToMock(*invoice)
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
}

// RemoveLineItem removes a line item from an invoice (DELETE /invoices/:id/line-items/:itemId)
func (h *Handlers) RemoveLineItem(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	itemID := c.Param("itemId")

	invoice, err := h.Services.Invoices.RemoveLineItem(ctx, id, itemID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to remove line item: "+err.Error())
	}

	mockInvoice := ConvertInvoiceToMock(*invoice)
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(&mockInvoice))
}

// GetClientsForSelect returns a list of clients for the invoice form select
func (h *Handlers) GetClientsForSelect(c echo.Context) error {
	ctx := c.Request().Context()

	clients, err := h.Services.Clients.List(ctx)
	if err != nil {
		return err
	}

	// Return as JSON for the select dropdown
	result := make([]map[string]string, len(clients))
	for i, client := range clients {
		result[i] = map[string]string{
			"id":      client.ID,
			"name":    client.Name,
			"company": client.Company,
		}
	}

	return c.JSON(http.StatusOK, result)
}

// Legacy functions for backwards compatibility (used by mock-based routes)
func InvoicesList(c echo.Context) error {
	status := c.QueryParam("status")
	var invoices []mock.Invoice

	if status == "" || status == "all" {
		invoices = mock.Invoices
	} else {
		invoices = mock.GetInvoicesByStatus(mock.InvoiceStatus(status))
	}

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		return render.Render(c, http.StatusOK, pages.InvoicesTableBody(invoices))
	}

	return render.Render(c, http.StatusOK, pages.InvoicesList(invoices, status))
}

func InvoiceDetail(c echo.Context) error {
	id := c.Param("id")
	invoice := mock.GetInvoice(id)
	if invoice == nil {
		return c.String(http.StatusNotFound, "Invoice not found")
	}
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(invoice))
}

func InvoiceEditor(c echo.Context) error {
	// For new invoice, pass nil
	return render.Render(c, http.StatusOK, pages.InvoiceEditor(nil))
}
