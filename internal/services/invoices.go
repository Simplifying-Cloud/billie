package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dscott/invoicey/internal/database/db"
	"github.com/google/uuid"
)

// InvoiceService handles invoice business logic
type InvoiceService struct {
	queries *db.Queries
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(queries *db.Queries) *InvoiceService {
	return &InvoiceService{queries: queries}
}

// InvoiceStatus represents invoice status
type InvoiceStatus string

const (
	StatusDraft   InvoiceStatus = "draft"
	StatusSent    InvoiceStatus = "sent"
	StatusPaid    InvoiceStatus = "paid"
	StatusOverdue InvoiceStatus = "overdue"
)

// LineItemData represents a line item for display
type LineItemData struct {
	ID          string
	Description string
	Quantity    float64
	Rate        float64
	Amount      float64
	SortOrder   int64
}

// InvoiceWithDetails combines invoice with client and line items
type InvoiceWithDetails struct {
	ID          string
	Number      string
	ClientID    string
	Client      *ClientWithStats
	Status      InvoiceStatus
	IssueDate   time.Time
	DueDate     time.Time
	LineItems   []LineItemData
	Subtotal    float64
	TaxRate     float64
	TaxAmount   float64
	Total       float64
	Notes       string
	PublicToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToInvoiceWithDetails converts db types to InvoiceWithDetails
func (s *InvoiceService) ToInvoiceWithDetails(ctx context.Context, inv db.Invoice) (*InvoiceWithDetails, error) {
	// Get client
	client, err := s.queries.GetClient(ctx, inv.ClientID)
	if err != nil {
		return nil, err
	}

	totalBilled := interfaceToFloat64(s.queries.GetClientTotalBilled(ctx, client.ID))
	invoiceCount, _ := s.queries.GetClientInvoiceCount(ctx, client.ID)
	clientData := ToClientWithStats(client, totalBilled, invoiceCount)

	// Get line items
	lineItems, err := s.queries.ListLineItemsByInvoice(ctx, inv.ID)
	if err != nil {
		return nil, err
	}

	items := make([]LineItemData, len(lineItems))
	for i, li := range lineItems {
		items[i] = LineItemData{
			ID:          li.ID,
			Description: li.Description,
			Quantity:    li.Quantity,
			Rate:        li.Rate,
			Amount:      li.Amount,
			SortOrder:   li.SortOrder,
		}
	}

	return &InvoiceWithDetails{
		ID:          inv.ID,
		Number:      inv.Number,
		ClientID:    inv.ClientID,
		Client:      &clientData,
		Status:      InvoiceStatus(inv.Status),
		IssueDate:   inv.IssueDate,
		DueDate:     inv.DueDate,
		LineItems:   items,
		Subtotal:    inv.Subtotal,
		TaxRate:     inv.TaxRate,
		TaxAmount:   inv.TaxAmount,
		Total:       inv.Total,
		Notes:       nullStringToString(inv.Notes),
		PublicToken: nullStringToString(inv.PublicToken),
		CreatedAt:   nullTimeToTime(inv.CreatedAt),
		UpdatedAt:   nullTimeToTime(inv.UpdatedAt),
	}, nil
}

// List returns all invoices with details
func (s *InvoiceService) List(ctx context.Context) ([]InvoiceWithDetails, error) {
	invoices, err := s.queries.ListInvoices(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]InvoiceWithDetails, 0, len(invoices))
	for _, inv := range invoices {
		details, err := s.ToInvoiceWithDetails(ctx, inv)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// ListByStatus returns invoices filtered by status
func (s *InvoiceService) ListByStatus(ctx context.Context, status string) ([]InvoiceWithDetails, error) {
	invoices, err := s.queries.ListInvoicesByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	result := make([]InvoiceWithDetails, 0, len(invoices))
	for _, inv := range invoices {
		details, err := s.ToInvoiceWithDetails(ctx, inv)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// ListByClient returns invoices for a specific client
func (s *InvoiceService) ListByClient(ctx context.Context, clientID string) ([]InvoiceWithDetails, error) {
	invoices, err := s.queries.ListInvoicesByClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	result := make([]InvoiceWithDetails, 0, len(invoices))
	for _, inv := range invoices {
		details, err := s.ToInvoiceWithDetails(ctx, inv)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// GetRecent returns recent invoices
func (s *InvoiceService) GetRecent(ctx context.Context, limit int64) ([]InvoiceWithDetails, error) {
	invoices, err := s.queries.GetRecentInvoices(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]InvoiceWithDetails, 0, len(invoices))
	for _, inv := range invoices {
		details, err := s.ToInvoiceWithDetails(ctx, inv)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// Get returns a single invoice by ID
func (s *InvoiceService) Get(ctx context.Context, id string) (*InvoiceWithDetails, error) {
	invoice, err := s.queries.GetInvoice(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.ToInvoiceWithDetails(ctx, invoice)
}

// GetByToken returns an invoice by its public token
func (s *InvoiceService) GetByToken(ctx context.Context, token string) (*InvoiceWithDetails, error) {
	invoice, err := s.queries.GetInvoiceByToken(ctx, stringToNullString(token))
	if err != nil {
		return nil, err
	}
	return s.ToInvoiceWithDetails(ctx, invoice)
}

// LineItemInput represents input for creating/updating a line item
type LineItemInput struct {
	Description string
	Quantity    float64
	Rate        float64
}

// CreateInvoiceParams holds parameters for creating an invoice
type CreateInvoiceParams struct {
	ClientID  string
	IssueDate time.Time
	DueDate   time.Time
	TaxRate   float64
	Notes     string
	LineItems []LineItemInput
}

// GenerateInvoiceNumber generates a unique invoice number
func (s *InvoiceService) GenerateInvoiceNumber(ctx context.Context, clientID string) (string, error) {
	// Get client to derive prefix from company name
	client, err := s.queries.GetClient(ctx, clientID)
	if err != nil {
		return "", fmt.Errorf("failed to get client: %w", err)
	}

	// Generate prefix from company name (first 4 chars uppercase)
	prefix := generatePrefix(nullStringToString(client.Company), client.Name)
	year := int64(time.Now().Year())

	// Get next sequence number atomically
	seq, err := s.queries.GetNextInvoiceSequence(ctx, db.GetNextInvoiceSequenceParams{
		ClientPrefix: prefix,
		Year:         year,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get sequence: %w", err)
	}

	// Format: PREFIX-YYYY-NNNN
	return fmt.Sprintf("%s-%d-%04d", prefix, year, seq), nil
}

// generatePrefix creates a 4-letter prefix from company or name
func generatePrefix(company, name string) string {
	source := strings.TrimSpace(company)
	if source == "" {
		source = strings.TrimSpace(name)
	}
	if source == "" {
		return "INV"
	}

	// Remove special characters and get first 4 uppercase letters
	source = strings.ToUpper(source)
	var letters []rune
	for _, r := range source {
		if r >= 'A' && r <= 'Z' {
			letters = append(letters, r)
			if len(letters) >= 4 {
				break
			}
		}
	}

	if len(letters) == 0 {
		return "INV"
	}
	if len(letters) < 4 {
		return string(letters)
	}
	return string(letters[:4])
}

// Create creates a new invoice with line items
func (s *InvoiceService) Create(ctx context.Context, params CreateInvoiceParams) (*InvoiceWithDetails, error) {
	// Generate invoice number
	number, err := s.GenerateInvoiceNumber(ctx, params.ClientID)
	if err != nil {
		return nil, err
	}

	// Calculate totals from line items
	var subtotal float64
	for _, item := range params.LineItems {
		subtotal += item.Quantity * item.Rate
	}
	taxAmount := subtotal * (params.TaxRate / 100)
	total := subtotal + taxAmount

	// Generate IDs
	invoiceID := uuid.New().String()
	publicToken := uuid.New().String()

	// Create invoice
	invoice, err := s.queries.CreateInvoice(ctx, db.CreateInvoiceParams{
		ID:          invoiceID,
		Number:      number,
		ClientID:    params.ClientID,
		Status:      string(StatusDraft),
		IssueDate:   params.IssueDate,
		DueDate:     params.DueDate,
		Subtotal:    subtotal,
		TaxRate:     params.TaxRate,
		TaxAmount:   taxAmount,
		Total:       total,
		Notes:       stringToNullString(params.Notes),
		PublicToken: stringToNullString(publicToken),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create line items
	for i, item := range params.LineItems {
		_, err := s.queries.CreateLineItem(ctx, db.CreateLineItemParams{
			ID:          uuid.New().String(),
			InvoiceID:   invoiceID,
			Description: item.Description,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      item.Quantity * item.Rate,
			SortOrder:   int64(i),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create line item: %w", err)
		}
	}

	return s.ToInvoiceWithDetails(ctx, invoice)
}

// UpdateInvoiceParams holds parameters for updating an invoice
type UpdateInvoiceParams struct {
	ClientID  string
	IssueDate time.Time
	DueDate   time.Time
	TaxRate   float64
	Notes     string
	LineItems []LineItemInput
}

// Update updates an existing invoice
func (s *InvoiceService) Update(ctx context.Context, id string, params UpdateInvoiceParams) (*InvoiceWithDetails, error) {
	// Get existing invoice to preserve certain fields
	existing, err := s.queries.GetInvoice(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete existing line items
	if err := s.queries.DeleteLineItemsByInvoice(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to delete line items: %w", err)
	}

	// Calculate totals from new line items
	var subtotal float64
	for _, item := range params.LineItems {
		subtotal += item.Quantity * item.Rate
	}
	taxAmount := subtotal * (params.TaxRate / 100)
	total := subtotal + taxAmount

	// Update invoice
	invoice, err := s.queries.UpdateInvoice(ctx, db.UpdateInvoiceParams{
		ID:        id,
		ClientID:  params.ClientID,
		Status:    existing.Status,
		IssueDate: params.IssueDate,
		DueDate:   params.DueDate,
		Subtotal:  subtotal,
		TaxRate:   params.TaxRate,
		TaxAmount: taxAmount,
		Total:     total,
		Notes:     stringToNullString(params.Notes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	// Create new line items
	for i, item := range params.LineItems {
		_, err := s.queries.CreateLineItem(ctx, db.CreateLineItemParams{
			ID:          uuid.New().String(),
			InvoiceID:   id,
			Description: item.Description,
			Quantity:    item.Quantity,
			Rate:        item.Rate,
			Amount:      item.Quantity * item.Rate,
			SortOrder:   int64(i),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create line item: %w", err)
		}
	}

	return s.ToInvoiceWithDetails(ctx, invoice)
}

// UpdateStatus updates an invoice's status
func (s *InvoiceService) UpdateStatus(ctx context.Context, id string, status InvoiceStatus) (*InvoiceWithDetails, error) {
	invoice, err := s.queries.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		return nil, err
	}
	return s.ToInvoiceWithDetails(ctx, invoice)
}

// MarkAsPaid marks an invoice as paid
func (s *InvoiceService) MarkAsPaid(ctx context.Context, id string) (*InvoiceWithDetails, error) {
	return s.UpdateStatus(ctx, id, StatusPaid)
}

// Send marks an invoice as sent
func (s *InvoiceService) Send(ctx context.Context, id string) (*InvoiceWithDetails, error) {
	return s.UpdateStatus(ctx, id, StatusSent)
}

// Delete deletes an invoice and its line items
func (s *InvoiceService) Delete(ctx context.Context, id string) error {
	// Line items will be deleted by CASCADE
	return s.queries.DeleteInvoice(ctx, id)
}

// CheckOverdue marks overdue invoices
func (s *InvoiceService) CheckOverdue(ctx context.Context) error {
	return s.queries.MarkOverdueInvoices(ctx)
}

// AddLineItem adds a line item to an invoice
func (s *InvoiceService) AddLineItem(ctx context.Context, invoiceID string, item LineItemInput) (*InvoiceWithDetails, error) {
	// Get current item count for sort order
	count, _ := s.queries.CountLineItemsByInvoice(ctx, invoiceID)

	// Create line item
	_, err := s.queries.CreateLineItem(ctx, db.CreateLineItemParams{
		ID:          uuid.New().String(),
		InvoiceID:   invoiceID,
		Description: item.Description,
		Quantity:    item.Quantity,
		Rate:        item.Rate,
		Amount:      item.Quantity * item.Rate,
		SortOrder:   count,
	})
	if err != nil {
		return nil, err
	}

	// Recalculate totals
	return s.recalculateTotals(ctx, invoiceID)
}

// RemoveLineItem removes a line item from an invoice
func (s *InvoiceService) RemoveLineItem(ctx context.Context, invoiceID, lineItemID string) (*InvoiceWithDetails, error) {
	if err := s.queries.DeleteLineItem(ctx, lineItemID); err != nil {
		return nil, err
	}
	return s.recalculateTotals(ctx, invoiceID)
}

// recalculateTotals recalculates invoice totals based on line items
func (s *InvoiceService) recalculateTotals(ctx context.Context, invoiceID string) (*InvoiceWithDetails, error) {
	invoice, err := s.queries.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	lineItems, err := s.queries.ListLineItemsByInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	var subtotal float64
	for _, li := range lineItems {
		subtotal += li.Amount
	}

	taxAmount := subtotal * (invoice.TaxRate / 100)
	total := subtotal + taxAmount

	updatedInvoice, err := s.queries.UpdateInvoiceTotals(ctx, db.UpdateInvoiceTotalsParams{
		ID:        invoiceID,
		Subtotal:  subtotal,
		TaxAmount: taxAmount,
		Total:     total,
	})
	if err != nil {
		return nil, err
	}

	return s.ToInvoiceWithDetails(ctx, updatedInvoice)
}
