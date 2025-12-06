package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) PublicInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.Param("token")

	invoice, err := h.Services.Invoices.GetByToken(ctx, token)
	if err != nil {
		return c.String(http.StatusNotFound, "Invoice not found")
	}

	mockInvoice := ConvertInvoiceToMock(*invoice)
	return render.Render(c, http.StatusOK, pages.PublicInvoice(mockInvoice))
}

// Legacy function for backwards compatibility
func PublicInvoice(c echo.Context) error {
	token := c.Param("token")
	invoice := mock.GetInvoiceByToken(token)
	if invoice == nil {
		return c.String(http.StatusNotFound, "Invoice not found")
	}
	return render.Render(c, http.StatusOK, pages.PublicInvoice(*invoice))
}
