package services

import (
	"github.com/dscott/invoicey/internal/database/db"
)

// Services holds all business logic services
type Services struct {
	Clients  *ClientService
	Invoices *InvoiceService
	Expenses *ExpenseService
	Metrics  *MetricsService
	Uploads  *UploadService
}

// New creates a new Services instance with all dependencies
func New(queries *db.Queries, uploadDir string) *Services {
	return &Services{
		Clients:  NewClientService(queries),
		Invoices: NewInvoiceService(queries),
		Expenses: NewExpenseService(queries),
		Metrics:  NewMetricsService(queries),
		Uploads:  NewUploadService(uploadDir),
	}
}
