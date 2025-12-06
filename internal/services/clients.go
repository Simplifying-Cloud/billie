package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/dscott/invoicey/internal/database/db"
	"github.com/google/uuid"
)

// ClientService handles client business logic
type ClientService struct {
	queries *db.Queries
}

// NewClientService creates a new client service
func NewClientService(queries *db.Queries) *ClientService {
	return &ClientService{queries: queries}
}

// ClientWithStats combines client data with computed stats
type ClientWithStats struct {
	ID             string
	Name           string
	Email          string
	Phone          string
	Company        string
	AddressStreet  string
	AddressCity    string
	AddressState   string
	AddressZip     string
	AddressCountry string
	TotalBilled    float64
	InvoiceCount   int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ToClientWithStats converts a db.Client to ClientWithStats
func ToClientWithStats(c db.Client, totalBilled float64, invoiceCount int64) ClientWithStats {
	return ClientWithStats{
		ID:             c.ID,
		Name:           c.Name,
		Email:          c.Email,
		Phone:          nullStringToString(c.Phone),
		Company:        nullStringToString(c.Company),
		AddressStreet:  nullStringToString(c.AddressStreet),
		AddressCity:    nullStringToString(c.AddressCity),
		AddressState:   nullStringToString(c.AddressState),
		AddressZip:     nullStringToString(c.AddressZip),
		AddressCountry: nullStringToString(c.AddressCountry),
		TotalBilled:    totalBilled,
		InvoiceCount:   invoiceCount,
		CreatedAt:      nullTimeToTime(c.CreatedAt),
		UpdatedAt:      nullTimeToTime(c.UpdatedAt),
	}
}

// List returns all clients
func (s *ClientService) List(ctx context.Context) ([]ClientWithStats, error) {
	clients, err := s.queries.ListClients(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ClientWithStats, len(clients))
	for i, c := range clients {
		totalBilled := interfaceToFloat64(s.queries.GetClientTotalBilled(ctx, c.ID))
		invoiceCount, _ := s.queries.GetClientInvoiceCount(ctx, c.ID)
		result[i] = ToClientWithStats(c, totalBilled, invoiceCount)
	}
	return result, nil
}

// Get returns a single client by ID with stats
func (s *ClientService) Get(ctx context.Context, id string) (*ClientWithStats, error) {
	client, err := s.queries.GetClient(ctx, id)
	if err != nil {
		return nil, err
	}

	totalBilled := interfaceToFloat64(s.queries.GetClientTotalBilled(ctx, client.ID))
	invoiceCount, _ := s.queries.GetClientInvoiceCount(ctx, client.ID)
	result := ToClientWithStats(client, totalBilled, invoiceCount)
	return &result, nil
}

// Search searches clients by name, company, or email
func (s *ClientService) Search(ctx context.Context, query string) ([]ClientWithStats, error) {
	searchPattern := "%" + query + "%"
	clients, err := s.queries.SearchClients(ctx, db.SearchClientsParams{
		Name:    searchPattern,
		Company: stringToNullString(searchPattern),
		Email:   searchPattern,
	})
	if err != nil {
		return nil, err
	}

	result := make([]ClientWithStats, len(clients))
	for i, c := range clients {
		totalBilled := interfaceToFloat64(s.queries.GetClientTotalBilled(ctx, c.ID))
		invoiceCount, _ := s.queries.GetClientInvoiceCount(ctx, c.ID)
		result[i] = ToClientWithStats(c, totalBilled, invoiceCount)
	}
	return result, nil
}

// CreateClientParams holds parameters for creating a client
type CreateClientParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Company        string `json:"company"`
	AddressStreet  string `json:"address_street"`
	AddressCity    string `json:"address_city"`
	AddressState   string `json:"address_state"`
	AddressZip     string `json:"address_zip"`
	AddressCountry string `json:"address_country"`
}

// Create creates a new client
func (s *ClientService) Create(ctx context.Context, params CreateClientParams) (*ClientWithStats, error) {
	id := uuid.New().String()

	country := params.AddressCountry
	if country == "" {
		country = "USA"
	}

	client, err := s.queries.CreateClient(ctx, db.CreateClientParams{
		ID:             id,
		Name:           params.Name,
		Email:          params.Email,
		Phone:          stringToNullString(params.Phone),
		Company:        stringToNullString(params.Company),
		AddressStreet:  stringToNullString(params.AddressStreet),
		AddressCity:    stringToNullString(params.AddressCity),
		AddressState:   stringToNullString(params.AddressState),
		AddressZip:     stringToNullString(params.AddressZip),
		AddressCountry: stringToNullString(country),
	})
	if err != nil {
		return nil, err
	}

	result := ToClientWithStats(client, 0, 0)
	return &result, nil
}

// UpdateClientParams holds parameters for updating a client
type UpdateClientParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Company        string `json:"company"`
	AddressStreet  string `json:"address_street"`
	AddressCity    string `json:"address_city"`
	AddressState   string `json:"address_state"`
	AddressZip     string `json:"address_zip"`
	AddressCountry string `json:"address_country"`
}

// Update updates an existing client
func (s *ClientService) Update(ctx context.Context, id string, params UpdateClientParams) (*ClientWithStats, error) {
	client, err := s.queries.UpdateClient(ctx, db.UpdateClientParams{
		ID:             id,
		Name:           params.Name,
		Email:          params.Email,
		Phone:          stringToNullString(params.Phone),
		Company:        stringToNullString(params.Company),
		AddressStreet:  stringToNullString(params.AddressStreet),
		AddressCity:    stringToNullString(params.AddressCity),
		AddressState:   stringToNullString(params.AddressState),
		AddressZip:     stringToNullString(params.AddressZip),
		AddressCountry: stringToNullString(params.AddressCountry),
	})
	if err != nil {
		return nil, err
	}

	totalBilled := interfaceToFloat64(s.queries.GetClientTotalBilled(ctx, client.ID))
	invoiceCount, _ := s.queries.GetClientInvoiceCount(ctx, client.ID)
	result := ToClientWithStats(client, totalBilled, invoiceCount)
	return &result, nil
}

// Delete deletes a client
func (s *ClientService) Delete(ctx context.Context, id string) error {
	return s.queries.DeleteClient(ctx, id)
}

// Helper functions
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullTimeToTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

// interfaceToFloat64 converts interface{} to float64 (for SQLC aggregate functions)
func interfaceToFloat64(val interface{}, err error) float64 {
	if err != nil || val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

// interfaceToString converts interface{} to string
func interfaceToString(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	default:
		return ""
	}
}
