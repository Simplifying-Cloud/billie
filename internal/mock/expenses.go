package mock

import "time"

type ExpenseCategory string

const (
	CategoryOffice    ExpenseCategory = "office"
	CategoryTravel    ExpenseCategory = "travel"
	CategorySoftware  ExpenseCategory = "software"
	CategoryMarketing ExpenseCategory = "marketing"
	CategoryUtilities ExpenseCategory = "utilities"
	CategoryOther     ExpenseCategory = "other"
)

type Expense struct {
	ID          string
	Description string
	Category    ExpenseCategory
	Merchant    string
	Amount      float64
	Date        time.Time
	ReceiptURL  string
	Notes       string
}

// Merchants is a list of known merchants
var Merchants = []string{
	"Adobe",
	"Amazon",
	"American Airlines",
	"Apple",
	"AT&T",
	"Comcast",
	"Delta Airlines",
	"Figma",
	"GitHub",
	"Google",
	"Hilton Hotels",
	"LinkedIn",
	"Marriott",
	"Microsoft",
	"Notion",
	"Slack",
	"Southwest Airlines",
	"Staples",
	"Stripe",
	"United Airlines",
	"Verizon",
	"Zoom",
}

var Expenses = []Expense{
	{
		ID:          "exp_001",
		Description: "Adobe Creative Cloud Annual",
		Category:    CategorySoftware,
		Merchant:    "Adobe",
		Amount:      599.88,
		Date:        time.Now().AddDate(0, 0, -15),
		ReceiptURL:  "/static/mock/receipt-001.jpg",
		Notes:       "Annual subscription renewal",
	},
	{
		ID:          "exp_002",
		Description: "Client Lunch Meeting - TechCorp",
		Category:    CategoryMarketing,
		Merchant:    "The Capital Grille",
		Amount:      127.50,
		Date:        time.Now().AddDate(0, 0, -12),
		ReceiptURL:  "/static/mock/receipt-002.jpg",
		Notes:       "Business development lunch with Sarah Johnson",
	},
	{
		ID:          "exp_003",
		Description: "Figma Pro Subscription",
		Category:    CategorySoftware,
		Merchant:    "Figma",
		Amount:      15.00,
		Date:        time.Now().AddDate(0, 0, -10),
		ReceiptURL:  "",
		Notes:       "Monthly design tool",
	},
	{
		ID:          "exp_004",
		Description: "Office Supplies - Amazon",
		Category:    CategoryOffice,
		Merchant:    "Amazon",
		Amount:      89.43,
		Date:        time.Now().AddDate(0, 0, -8),
		ReceiptURL:  "/static/mock/receipt-004.jpg",
		Notes:       "Notebooks, pens, sticky notes",
	},
	{
		ID:          "exp_005",
		Description: "Flight to Austin - Client Visit",
		Category:    CategoryTravel,
		Merchant:    "Southwest Airlines",
		Amount:      342.00,
		Date:        time.Now().AddDate(0, 0, -5),
		ReceiptURL:  "/static/mock/receipt-005.jpg",
		Notes:       "Round trip for StartupLab project kickoff",
	},
	{
		ID:          "exp_006",
		Description: "Hotel - Austin",
		Category:    CategoryTravel,
		Merchant:    "Hilton Hotels",
		Amount:      189.00,
		Date:        time.Now().AddDate(0, 0, -4),
		ReceiptURL:  "/static/mock/receipt-006.jpg",
		Notes:       "One night stay",
	},
	{
		ID:          "exp_007",
		Description: "Internet Service - Monthly",
		Category:    CategoryUtilities,
		Merchant:    "Comcast",
		Amount:      79.99,
		Date:        time.Now().AddDate(0, 0, -3),
		ReceiptURL:  "",
		Notes:       "High-speed fiber connection",
	},
	{
		ID:          "exp_008",
		Description: "GitHub Pro",
		Category:    CategorySoftware,
		Merchant:    "GitHub",
		Amount:      7.00,
		Date:        time.Now().AddDate(0, 0, -1),
		ReceiptURL:  "",
		Notes:       "Monthly subscription",
	},
}

var CategoryLabels = map[ExpenseCategory]string{
	CategoryOffice:    "Office Supplies",
	CategoryTravel:    "Travel",
	CategorySoftware:  "Software & Tools",
	CategoryMarketing: "Marketing",
	CategoryUtilities: "Utilities",
	CategoryOther:     "Other",
}

func GetExpense(id string) *Expense {
	for _, e := range Expenses {
		if e.ID == id {
			return &e
		}
	}
	return nil
}

func GetExpensesByCategory(category ExpenseCategory) []Expense {
	if category == "" {
		return Expenses
	}
	var results []Expense
	for _, e := range Expenses {
		if e.Category == category {
			results = append(results, e)
		}
	}
	return results
}

func GetTotalExpenses() float64 {
	var total float64
	for _, e := range Expenses {
		total += e.Amount
	}
	return total
}

func GetRecentExpenses(limit int) []Expense {
	if limit >= len(Expenses) {
		return Expenses
	}
	return Expenses[:limit]
}
