package main

import (
	"net/http"

	"github.com/dscott/invoicey/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", handlers.Dashboard)
	e.GET("/clients", handlers.ClientsList)
	e.GET("/clients/search", handlers.ClientsSearch)
	e.GET("/clients/:id", handlers.ClientDetail)
	e.GET("/invoices", handlers.InvoicesList)
	e.GET("/invoices/new", handlers.InvoiceEditor)
	e.GET("/invoices/:id", handlers.InvoiceDetail)
	e.GET("/expenses", handlers.ExpensesList)
	e.GET("/p/:token", handlers.PublicInvoice)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

// Custom error handler
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	c.String(code, "Something went wrong")
}
