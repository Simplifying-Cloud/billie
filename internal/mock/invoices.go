package mock

import "time"

type InvoiceStatus string

const (
	StatusDraft   InvoiceStatus = "draft"
	StatusSent    InvoiceStatus = "sent"
	StatusPaid    InvoiceStatus = "paid"
	StatusOverdue InvoiceStatus = "overdue"
)

type LineItem struct {
	Description string
	Quantity    float64
	Rate        float64
	Amount      float64
}

type Invoice struct {
	ID          string
	Number      string
	ClientID    string
	Client      *Client
	Status      InvoiceStatus
	IssueDate   time.Time
	DueDate     time.Time
	LineItems   []LineItem
	Subtotal    float64
	TaxRate     float64
	TaxAmount   float64
	Total       float64
	Notes       string
	PublicToken string
}

var Invoices = []Invoice{
	{
		ID:        "inv_001",
		Number:    "INV-2024-001",
		ClientID:  "cli_001",
		Status:    StatusPaid,
		IssueDate: time.Now().AddDate(0, -1, -15),
		DueDate:   time.Now().AddDate(0, -1, 15),
		LineItems: []LineItem{
			{Description: "Website Redesign", Quantity: 1, Rate: 5000.00, Amount: 5000.00},
			{Description: "SEO Optimization", Quantity: 10, Rate: 150.00, Amount: 1500.00},
		},
		Subtotal:    6500.00,
		TaxRate:     8.25,
		TaxAmount:   536.25,
		Total:       7036.25,
		Notes:       "Thank you for your business!",
		PublicToken: "abc123",
	},
	{
		ID:        "inv_002",
		Number:    "INV-2024-002",
		ClientID:  "cli_002",
		Status:    StatusSent,
		IssueDate: time.Now().AddDate(0, 0, -10),
		DueDate:   time.Now().AddDate(0, 0, 20),
		LineItems: []LineItem{
			{Description: "Mobile App Development - Phase 1", Quantity: 1, Rate: 12000.00, Amount: 12000.00},
			{Description: "UI/UX Design", Quantity: 40, Rate: 125.00, Amount: 5000.00},
			{Description: "Project Management", Quantity: 20, Rate: 100.00, Amount: 2000.00},
		},
		Subtotal:    19000.00,
		TaxRate:     8.25,
		TaxAmount:   1567.50,
		Total:       20567.50,
		Notes:       "Payment due within 30 days.",
		PublicToken: "def456",
	},
	{
		ID:        "inv_003",
		Number:    "INV-2024-003",
		ClientID:  "cli_003",
		Status:    StatusDraft,
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		LineItems: []LineItem{
			{Description: "Brand Identity Package", Quantity: 1, Rate: 3500.00, Amount: 3500.00},
			{Description: "Logo Design Revisions", Quantity: 3, Rate: 250.00, Amount: 750.00},
		},
		Subtotal:    4250.00,
		TaxRate:     8.25,
		TaxAmount:   350.63,
		Total:       4600.63,
		Notes:       "",
		PublicToken: "ghi789",
	},
	{
		ID:        "inv_004",
		Number:    "INV-2024-004",
		ClientID:  "cli_004",
		Status:    StatusOverdue,
		IssueDate: time.Now().AddDate(0, -2, 0),
		DueDate:   time.Now().AddDate(0, -1, 0),
		LineItems: []LineItem{
			{Description: "Financial Dashboard Development", Quantity: 1, Rate: 15000.00, Amount: 15000.00},
			{Description: "Data Integration Services", Quantity: 1, Rate: 5000.00, Amount: 5000.00},
			{Description: "Training Sessions", Quantity: 8, Rate: 200.00, Amount: 1600.00},
		},
		Subtotal:    21600.00,
		TaxRate:     8.25,
		TaxAmount:   1782.00,
		Total:       23382.00,
		Notes:       "OVERDUE - Please remit payment immediately.",
		PublicToken: "jkl012",
	},
	{
		ID:        "inv_005",
		Number:    "INV-2024-005",
		ClientID:  "cli_005",
		Status:    StatusPaid,
		IssueDate: time.Now().AddDate(0, 0, -20),
		DueDate:   time.Now().AddDate(0, 0, 10),
		LineItems: []LineItem{
			{Description: "Sustainability Report Design", Quantity: 1, Rate: 2500.00, Amount: 2500.00},
			{Description: "Infographic Creation", Quantity: 5, Rate: 400.00, Amount: 2000.00},
		},
		Subtotal:    4500.00,
		TaxRate:     8.25,
		TaxAmount:   371.25,
		Total:       4871.25,
		Notes:       "Thank you for choosing eco-friendly solutions!",
		PublicToken: "mno345",
	},
	{
		ID:        "inv_006",
		Number:    "INV-2024-006",
		ClientID:  "cli_006",
		Status:    StatusSent,
		IssueDate: time.Now().AddDate(0, 0, -5),
		DueDate:   time.Now().AddDate(0, 0, 25),
		LineItems: []LineItem{
			{Description: "Cloud Migration Consulting", Quantity: 40, Rate: 200.00, Amount: 8000.00},
			{Description: "AWS Setup & Configuration", Quantity: 1, Rate: 3000.00, Amount: 3000.00},
			{Description: "Security Audit", Quantity: 1, Rate: 2500.00, Amount: 2500.00},
		},
		Subtotal:    13500.00,
		TaxRate:     8.25,
		TaxAmount:   1113.75,
		Total:       14613.75,
		Notes:       "",
		PublicToken: "pqr678",
	},
	{
		ID:        "inv_007",
		Number:    "INV-2024-007",
		ClientID:  "cli_002",
		Status:    StatusDraft,
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		LineItems: []LineItem{
			{Description: "Mobile App Development - Phase 2", Quantity: 1, Rate: 15000.00, Amount: 15000.00},
		},
		Subtotal:    15000.00,
		TaxRate:     8.25,
		TaxAmount:   1237.50,
		Total:       16237.50,
		Notes:       "",
		PublicToken: "stu901",
	},
}

func init() {
	// Link clients to invoices
	for i := range Invoices {
		Invoices[i].Client = GetClient(Invoices[i].ClientID)
	}
}

func GetInvoice(id string) *Invoice {
	for _, inv := range Invoices {
		if inv.ID == id {
			return &inv
		}
	}
	return nil
}

func GetInvoiceByToken(token string) *Invoice {
	for _, inv := range Invoices {
		if inv.PublicToken == token {
			return &inv
		}
	}
	return nil
}

func GetInvoicesByStatus(status InvoiceStatus) []Invoice {
	if status == "" {
		return Invoices
	}
	var results []Invoice
	for _, inv := range Invoices {
		if inv.Status == status {
			results = append(results, inv)
		}
	}
	return results
}

func GetRecentInvoices(limit int) []Invoice {
	if limit >= len(Invoices) {
		return Invoices
	}
	return Invoices[:limit]
}
