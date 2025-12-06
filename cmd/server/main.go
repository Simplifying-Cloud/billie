package main

import (
	"log"

	"github.com/dscott/invoicey/internal/database"
	"github.com/dscott/invoicey/internal/database/db"
	"github.com/dscott/invoicey/internal/handlers"
	"github.com/dscott/invoicey/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Open database
	sqlDB, err := database.Open("./billie.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := database.RunMigrations(sqlDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize SQLC queries
	queries := db.New(sqlDB)

	// Initialize services
	svc := services.New(queries, "./uploads")

	// Initialize handlers
	h := handlers.New(svc)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Static files
	e.Static("/static", "static")
	e.Static("/uploads", "uploads")

	// HTML Routes (HTMX)
	e.GET("/", h.Dashboard)
	e.GET("/clients", h.ClientsList)
	e.GET("/clients/search", h.ClientsSearch)
	e.GET("/clients/:id", h.ClientDetail)
	e.POST("/clients", h.CreateClient)
	e.PUT("/clients/:id", h.UpdateClient)
	e.DELETE("/clients/:id", h.DeleteClient)

	e.GET("/invoices", h.InvoicesList)
	e.GET("/invoices/new", h.InvoiceEditor)
	e.GET("/invoices/:id", h.InvoiceDetail)
	e.POST("/invoices", h.CreateInvoice)
	e.PUT("/invoices/:id", h.UpdateInvoice)
	e.DELETE("/invoices/:id", h.DeleteInvoice)
	e.POST("/invoices/:id/send", h.SendInvoice)
	e.POST("/invoices/:id/mark-paid", h.MarkInvoicePaid)
	e.POST("/invoices/:id/line-items", h.AddLineItem)
	e.DELETE("/invoices/:id/line-items/:itemId", h.RemoveLineItem)

	e.GET("/expenses", h.ExpensesList)
	e.POST("/expenses", h.CreateExpense)
	e.PUT("/expenses/:id", h.UpdateExpense)
	e.DELETE("/expenses/:id", h.DeleteExpense)
	e.POST("/expenses/upload-receipt", h.UploadReceipt)

	e.GET("/p/:token", h.PublicInvoice)

	// API Routes (JSON)
	api := e.Group("/api/v1")

	// Clients API
	api.GET("/clients", h.APIListClients)
	api.GET("/clients/:id", h.APIGetClient)
	api.POST("/clients", h.APICreateClient)
	api.PUT("/clients/:id", h.APIUpdateClient)
	api.DELETE("/clients/:id", h.APIDeleteClient)

	// Invoices API
	api.GET("/invoices", h.APIListInvoices)
	api.GET("/invoices/:id", h.APIGetInvoice)
	api.POST("/invoices", h.APICreateInvoice)
	api.PUT("/invoices/:id", h.APIUpdateInvoice)
	api.DELETE("/invoices/:id", h.APIDeleteInvoice)
	api.POST("/invoices/:id/send", h.APISendInvoice)
	api.POST("/invoices/:id/mark-paid", h.APIMarkPaid)

	// Expenses API
	api.GET("/expenses", h.APIListExpenses)
	api.GET("/expenses/:id", h.APIGetExpense)
	api.POST("/expenses", h.APICreateExpense)
	api.PUT("/expenses/:id", h.APIUpdateExpense)
	api.DELETE("/expenses/:id", h.APIDeleteExpense)

	// Lookups API
	api.GET("/categories", h.ListCategories)
	api.GET("/merchants", h.ListMerchants)
	api.GET("/merchants/search", h.SearchMerchants)
	api.POST("/merchants", h.CreateMerchant)
	api.GET("/clients-select", h.GetClientsForSelect)

	// Metrics API
	api.GET("/metrics/dashboard", h.APIDashboardMetrics)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
