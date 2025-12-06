package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/dscott/invoicey/internal/database/db"
	"github.com/google/uuid"
)

// ExpenseService handles expense business logic
type ExpenseService struct {
	queries *db.Queries
}

// NewExpenseService creates a new expense service
func NewExpenseService(queries *db.Queries) *ExpenseService {
	return &ExpenseService{queries: queries}
}

// ExpenseWithDetails combines expense with category and merchant names
type ExpenseWithDetails struct {
	ID           string
	Description  string
	CategoryID   string
	CategoryName string
	CategorySlug string
	MerchantID   string
	MerchantName string
	Amount       float64
	Date         time.Time
	ReceiptPath  string
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CategoryData represents a category for display
type CategoryData struct {
	ID   string
	Name string
	Slug string
}

// MerchantData represents a merchant for display
type MerchantData struct {
	ID   string
	Name string
}

// ToExpenseWithDetails converts db.Expense to ExpenseWithDetails
func (s *ExpenseService) ToExpenseWithDetails(ctx context.Context, e db.Expense) (*ExpenseWithDetails, error) {
	var categoryName, categorySlug, merchantName string

	if e.CategoryID.Valid {
		cat, err := s.queries.GetCategory(ctx, e.CategoryID.String)
		if err == nil {
			categoryName = cat.Name
			categorySlug = cat.Slug
		}
	}

	if e.MerchantID.Valid {
		merchant, err := s.queries.GetMerchant(ctx, e.MerchantID.String)
		if err == nil {
			merchantName = merchant.Name
		}
	}

	return &ExpenseWithDetails{
		ID:           e.ID,
		Description:  e.Description,
		CategoryID:   nullStringToString(e.CategoryID),
		CategoryName: categoryName,
		CategorySlug: categorySlug,
		MerchantID:   nullStringToString(e.MerchantID),
		MerchantName: merchantName,
		Amount:       e.Amount,
		Date:         e.Date,
		ReceiptPath:  nullStringToString(e.ReceiptPath),
		Notes:        nullStringToString(e.Notes),
		CreatedAt:    nullTimeToTime(e.CreatedAt),
		UpdatedAt:    nullTimeToTime(e.UpdatedAt),
	}, nil
}

// List returns all expenses with details
func (s *ExpenseService) List(ctx context.Context) ([]ExpenseWithDetails, error) {
	expenses, err := s.queries.ListExpenses(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ExpenseWithDetails, 0, len(expenses))
	for _, e := range expenses {
		details, err := s.ToExpenseWithDetails(ctx, e)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// ListByCategory returns expenses filtered by category slug
func (s *ExpenseService) ListByCategory(ctx context.Context, categorySlug string) ([]ExpenseWithDetails, error) {
	// Get category by slug
	cat, err := s.queries.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	expenses, err := s.queries.ListExpensesByCategory(ctx, sql.NullString{String: cat.ID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]ExpenseWithDetails, 0, len(expenses))
	for _, e := range expenses {
		details, err := s.ToExpenseWithDetails(ctx, e)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// GetRecent returns recent expenses
func (s *ExpenseService) GetRecent(ctx context.Context, limit int64) ([]ExpenseWithDetails, error) {
	expenses, err := s.queries.GetRecentExpenses(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]ExpenseWithDetails, 0, len(expenses))
	for _, e := range expenses {
		details, err := s.ToExpenseWithDetails(ctx, e)
		if err != nil {
			continue
		}
		result = append(result, *details)
	}
	return result, nil
}

// Get returns a single expense by ID
func (s *ExpenseService) Get(ctx context.Context, id string) (*ExpenseWithDetails, error) {
	expense, err := s.queries.GetExpense(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.ToExpenseWithDetails(ctx, expense)
}

// CreateExpenseParams holds parameters for creating an expense
type CreateExpenseParams struct {
	Description string
	CategoryID  string
	MerchantID  string
	Amount      float64
	Date        time.Time
	ReceiptPath string
	Notes       string
}

// Create creates a new expense
func (s *ExpenseService) Create(ctx context.Context, params CreateExpenseParams) (*ExpenseWithDetails, error) {
	id := uuid.New().String()

	expense, err := s.queries.CreateExpense(ctx, db.CreateExpenseParams{
		ID:          id,
		Description: params.Description,
		CategoryID:  stringToNullString(params.CategoryID),
		MerchantID:  stringToNullString(params.MerchantID),
		Amount:      params.Amount,
		Date:        params.Date,
		ReceiptPath: stringToNullString(params.ReceiptPath),
		Notes:       stringToNullString(params.Notes),
	})
	if err != nil {
		return nil, err
	}

	return s.ToExpenseWithDetails(ctx, expense)
}

// UpdateExpenseParams holds parameters for updating an expense
type UpdateExpenseParams struct {
	Description string
	CategoryID  string
	MerchantID  string
	Amount      float64
	Date        time.Time
	ReceiptPath string
	Notes       string
}

// Update updates an existing expense
func (s *ExpenseService) Update(ctx context.Context, id string, params UpdateExpenseParams) (*ExpenseWithDetails, error) {
	expense, err := s.queries.UpdateExpense(ctx, db.UpdateExpenseParams{
		ID:          id,
		Description: params.Description,
		CategoryID:  stringToNullString(params.CategoryID),
		MerchantID:  stringToNullString(params.MerchantID),
		Amount:      params.Amount,
		Date:        params.Date,
		ReceiptPath: stringToNullString(params.ReceiptPath),
		Notes:       stringToNullString(params.Notes),
	})
	if err != nil {
		return nil, err
	}

	return s.ToExpenseWithDetails(ctx, expense)
}

// UpdateReceipt updates only the receipt path for an expense
func (s *ExpenseService) UpdateReceipt(ctx context.Context, id, receiptPath string) (*ExpenseWithDetails, error) {
	expense, err := s.queries.UpdateExpenseReceipt(ctx, db.UpdateExpenseReceiptParams{
		ID:          id,
		ReceiptPath: stringToNullString(receiptPath),
	})
	if err != nil {
		return nil, err
	}

	return s.ToExpenseWithDetails(ctx, expense)
}

// Delete deletes an expense
func (s *ExpenseService) Delete(ctx context.Context, id string) error {
	return s.queries.DeleteExpense(ctx, id)
}

// GetTotal returns total expenses
func (s *ExpenseService) GetTotal(ctx context.Context) (float64, error) {
	val, err := s.queries.GetTotalExpenses(ctx)
	if err != nil {
		return 0, err
	}
	return interfaceToFloat64(val, nil), nil
}

// ListCategories returns all expense categories
func (s *ExpenseService) ListCategories(ctx context.Context) ([]CategoryData, error) {
	categories, err := s.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryData, len(categories))
	for i, c := range categories {
		result[i] = CategoryData{
			ID:   c.ID,
			Name: c.Name,
			Slug: c.Slug,
		}
	}
	return result, nil
}

// CreateCategory creates a new expense category
func (s *ExpenseService) CreateCategory(ctx context.Context, name, slug string) (*CategoryData, error) {
	id := uuid.New().String()
	cat, err := s.queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   id,
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return nil, err
	}
	return &CategoryData{
		ID:   cat.ID,
		Name: cat.Name,
		Slug: cat.Slug,
	}, nil
}

// ListMerchants returns all merchants
func (s *ExpenseService) ListMerchants(ctx context.Context) ([]MerchantData, error) {
	merchants, err := s.queries.ListMerchants(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]MerchantData, len(merchants))
	for i, m := range merchants {
		result[i] = MerchantData{
			ID:   m.ID,
			Name: m.Name,
		}
	}
	return result, nil
}

// SearchMerchants searches merchants by name
func (s *ExpenseService) SearchMerchants(ctx context.Context, query string) ([]MerchantData, error) {
	merchants, err := s.queries.SearchMerchants(ctx, "%"+query+"%")
	if err != nil {
		return nil, err
	}

	result := make([]MerchantData, len(merchants))
	for i, m := range merchants {
		result[i] = MerchantData{
			ID:   m.ID,
			Name: m.Name,
		}
	}
	return result, nil
}

// CreateMerchant creates a new merchant
func (s *ExpenseService) CreateMerchant(ctx context.Context, name string) (*MerchantData, error) {
	id := uuid.New().String()
	merchant, err := s.queries.CreateMerchant(ctx, db.CreateMerchantParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &MerchantData{
		ID:   merchant.ID,
		Name: merchant.Name,
	}, nil
}

// GetOrCreateMerchant gets or creates a merchant by name
func (s *ExpenseService) GetOrCreateMerchant(ctx context.Context, name string) (*MerchantData, error) {
	merchant, err := s.queries.GetMerchantByName(ctx, name)
	if err == nil {
		return &MerchantData{
			ID:   merchant.ID,
			Name: merchant.Name,
		}, nil
	}

	// Create new merchant
	return s.CreateMerchant(ctx, name)
}
