package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/internal/services"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) ExpensesList(c echo.Context) error {
	ctx := c.Request().Context()
	category := c.QueryParam("category")

	var expenses []services.ExpenseWithDetails
	var err error

	if category == "" || category == "all" {
		expenses, err = h.Services.Expenses.List(ctx)
	} else {
		expenses, err = h.Services.Expenses.ListByCategory(ctx, category)
	}

	if err != nil {
		// Fallback to empty list if category not found
		expenses = []services.ExpenseWithDetails{}
	}

	mockExpenses := ConvertExpensesToMock(expenses)

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		return render.Render(c, http.StatusOK, pages.ExpensesTableBody(mockExpenses))
	}

	data := pages.GetExpensesPageData(mockExpenses, category)
	return render.Render(c, http.StatusOK, pages.ExpensesList(data))
}

// CreateExpense creates a new expense (POST /expenses)
func (h *Handlers) CreateExpense(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse date
	date, _ := time.Parse("2006-01-02", c.FormValue("date"))

	// Parse amount
	var amount float64
	if a := c.FormValue("amount"); a != "" {
		json.Unmarshal([]byte(a), &amount)
	}

	// Handle merchant - get or create
	merchantName := c.FormValue("merchant")
	var merchantID string
	if merchantName != "" {
		merchant, err := h.Services.Expenses.GetOrCreateMerchant(ctx, merchantName)
		if err == nil {
			merchantID = merchant.ID
		}
	}

	params := services.CreateExpenseParams{
		Description: c.FormValue("description"),
		CategoryID:  c.FormValue("category_id"),
		MerchantID:  merchantID,
		Amount:      amount,
		Date:        date,
		ReceiptPath: c.FormValue("receipt_path"),
		Notes:       c.FormValue("notes"),
	}

	expense, err := h.Services.Expenses.Create(ctx, params)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to create expense: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Trigger", "expense-created")
		// Return updated expense list
		expenses, _ := h.Services.Expenses.List(ctx)
		return render.Render(c, http.StatusOK, pages.ExpensesTableBody(ConvertExpensesToMock(expenses)))
	}

	return c.JSON(http.StatusCreated, ConvertExpenseToMock(*expense))
}

// UpdateExpense updates an existing expense (PUT /expenses/:id)
func (h *Handlers) UpdateExpense(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	// Parse date
	date, _ := time.Parse("2006-01-02", c.FormValue("date"))

	// Parse amount
	var amount float64
	if a := c.FormValue("amount"); a != "" {
		json.Unmarshal([]byte(a), &amount)
	}

	// Handle merchant - get or create
	merchantName := c.FormValue("merchant")
	var merchantID string
	if merchantName != "" {
		merchant, err := h.Services.Expenses.GetOrCreateMerchant(ctx, merchantName)
		if err == nil {
			merchantID = merchant.ID
		}
	}

	params := services.UpdateExpenseParams{
		Description: c.FormValue("description"),
		CategoryID:  c.FormValue("category_id"),
		MerchantID:  merchantID,
		Amount:      amount,
		Date:        date,
		ReceiptPath: c.FormValue("receipt_path"),
		Notes:       c.FormValue("notes"),
	}

	expense, err := h.Services.Expenses.Update(ctx, id, params)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to update expense: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Trigger", "expense-updated")
		expenses, _ := h.Services.Expenses.List(ctx)
		return render.Render(c, http.StatusOK, pages.ExpensesTableBody(ConvertExpensesToMock(expenses)))
	}

	return c.JSON(http.StatusOK, ConvertExpenseToMock(*expense))
}

// DeleteExpense deletes an expense (DELETE /expenses/:id)
func (h *Handlers) DeleteExpense(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.Services.Expenses.Delete(ctx, id); err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.String(http.StatusBadRequest, "Failed to delete expense: "+err.Error())
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Trigger", "expense-deleted")
		expenses, _ := h.Services.Expenses.List(ctx)
		return render.Render(c, http.StatusOK, pages.ExpensesTableBody(ConvertExpensesToMock(expenses)))
	}

	return c.NoContent(http.StatusNoContent)
}

// UploadReceipt handles receipt file upload (POST /expenses/upload-receipt)
func (h *Handlers) UploadReceipt(c echo.Context) error {
	file, err := c.FormFile("receipt")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Validate file
	if err := h.Services.Uploads.ValidateReceipt(file.Filename, file.Header.Get("Content-Type"), file.Size); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read file"})
	}
	defer src.Close()

	// Save file
	path, err := h.Services.Uploads.SaveReceipt(src, file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	return c.JSON(http.StatusOK, map[string]string{"path": path})
}

// ListCategories returns all expense categories
func (h *Handlers) ListCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.Services.Expenses.ListCategories(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}

// ListMerchants returns all merchants
func (h *Handlers) ListMerchants(c echo.Context) error {
	ctx := c.Request().Context()

	merchants, err := h.Services.Expenses.ListMerchants(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, merchants)
}

// SearchMerchants searches merchants by name
func (h *Handlers) SearchMerchants(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")

	merchants, err := h.Services.Expenses.SearchMerchants(ctx, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, merchants)
}

// CreateMerchant creates a new merchant
func (h *Handlers) CreateMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	name := c.FormValue("name")

	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}

	merchant, err := h.Services.Expenses.CreateMerchant(ctx, name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, merchant)
}

// Legacy function for backwards compatibility
func ExpensesList(c echo.Context) error {
	category := c.QueryParam("category")
	var expenses []mock.Expense

	if category == "" || category == "all" {
		expenses = mock.Expenses
	} else {
		expenses = mock.GetExpensesByCategory(mock.ExpenseCategory(category))
	}

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		return render.Render(c, http.StatusOK, pages.ExpensesTableBody(expenses))
	}

	data := pages.GetExpensesPageData(expenses, category)
	return render.Render(c, http.StatusOK, pages.ExpensesList(data))
}
