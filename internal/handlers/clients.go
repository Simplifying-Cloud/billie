package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/internal/services"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) ClientsList(c echo.Context) error {
	ctx := c.Request().Context()

	clients, err := h.Services.Clients.List(ctx)
	if err != nil {
		return err
	}

	return render.Render(c, http.StatusOK, pages.ClientsList(ConvertClientsToMock(clients)))
}

func (h *Handlers) ClientsSearch(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")

	clients, err := h.Services.Clients.Search(ctx, query)
	if err != nil {
		return err
	}

	return render.Render(c, http.StatusOK, pages.ClientsTableBody(ConvertClientsToMock(clients)))
}

func (h *Handlers) ClientDetail(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	client, err := h.Services.Clients.Get(ctx, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Client not found")
	}

	// Get client's invoices
	invoices, err := h.Services.Invoices.ListByClient(ctx, id)
	if err != nil {
		return err
	}

	mockClient := ConvertClientToMock(*client)
	mockInvoices := ConvertInvoicesToMock(invoices)

	return render.Render(c, http.StatusOK, pages.ClientDetail(mockClient, mockInvoices))
}

// CreateClient creates a new client (POST /clients)
func (h *Handlers) CreateClient(c echo.Context) error {
	ctx := c.Request().Context()

	params := services.CreateClientParams{
		Name:           c.FormValue("name"),
		Email:          c.FormValue("email"),
		Phone:          c.FormValue("phone"),
		Company:        c.FormValue("company"),
		AddressStreet:  c.FormValue("address_street"),
		AddressCity:    c.FormValue("address_city"),
		AddressState:   c.FormValue("address_state"),
		AddressZip:     c.FormValue("address_zip"),
		AddressCountry: c.FormValue("address_country"),
	}

	client, err := h.Services.Clients.Create(ctx, params)
	if err != nil {
		// Return error message for HTMX
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to create client: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// For HTMX: redirect to clients list or return updated table
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Trigger", "client-created")
		c.Response().Header().Set("HX-Redirect", "/clients/"+client.ID)
		return c.NoContent(http.StatusOK)
	}

	// For JSON API
	return c.JSON(http.StatusCreated, ConvertClientToMock(*client))
}

// UpdateClient updates an existing client (PUT /clients/:id)
func (h *Handlers) UpdateClient(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	params := services.UpdateClientParams{
		Name:           c.FormValue("name"),
		Email:          c.FormValue("email"),
		Phone:          c.FormValue("phone"),
		Company:        c.FormValue("company"),
		AddressStreet:  c.FormValue("address_street"),
		AddressCity:    c.FormValue("address_city"),
		AddressState:   c.FormValue("address_state"),
		AddressZip:     c.FormValue("address_zip"),
		AddressCountry: c.FormValue("address_country"),
	}

	client, err := h.Services.Clients.Update(ctx, id, params)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to update client: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Trigger", "client-updated")
		// Return the updated client detail page
		invoices, _ := h.Services.Invoices.ListByClient(ctx, id)
		return render.Render(c, http.StatusOK, pages.ClientDetail(ConvertClientToMock(*client), ConvertInvoicesToMock(invoices)))
	}

	return c.JSON(http.StatusOK, ConvertClientToMock(*client))
}

// DeleteClient deletes a client (DELETE /clients/:id)
func (h *Handlers) DeleteClient(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Clients.Delete(ctx, id); err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to delete client: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/clients")
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusNoContent)
}
