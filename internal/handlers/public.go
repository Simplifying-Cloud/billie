package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func PublicInvoice(c echo.Context) error {
	token := c.Param("token")
	invoice := mock.GetInvoiceByToken(token)
	if invoice == nil {
		return c.String(http.StatusNotFound, "Invoice not found")
	}
	return render.Render(c, http.StatusOK, pages.PublicInvoice(*invoice))
}
