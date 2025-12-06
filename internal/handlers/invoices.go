package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

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
