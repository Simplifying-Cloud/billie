package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func ClientsList(c echo.Context) error {
	clients := mock.Clients
	return render.Render(c, http.StatusOK, pages.ClientsList(clients))
}

func ClientsSearch(c echo.Context) error {
	query := c.QueryParam("q")
	clients := mock.SearchClients(query)
	return render.Render(c, http.StatusOK, pages.ClientsTableBody(clients))
}

func ClientDetail(c echo.Context) error {
	id := c.Param("id")
	client := mock.GetClient(id)
	if client == nil {
		return c.String(http.StatusNotFound, "Client not found")
	}

	// Get client's invoices
	var clientInvoices []mock.Invoice
	for _, inv := range mock.Invoices {
		if inv.ClientID == id {
			clientInvoices = append(clientInvoices, inv)
		}
	}

	return render.Render(c, http.StatusOK, pages.ClientDetail(*client, clientInvoices))
}
