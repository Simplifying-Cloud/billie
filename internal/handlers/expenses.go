package handlers

import (
	"net/http"

	"github.com/dscott/invoicey/internal/mock"
	"github.com/dscott/invoicey/internal/render"
	"github.com/dscott/invoicey/views/pages"
	"github.com/labstack/echo/v4"
)

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
