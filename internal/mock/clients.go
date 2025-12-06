package mock

import "time"

type Address struct {
	Street  string
	City    string
	State   string
	Zip     string
	Country string
}

type Client struct {
	ID           string
	Name         string
	Email        string
	Phone        string
	Company      string
	Address      Address
	TotalBilled  float64
	InvoiceCount int
	CreatedAt    time.Time
}

var Clients = []Client{
	{
		ID:           "cli_001",
		Name:         "Sarah Johnson",
		Email:        "sarah@techcorp.io",
		Phone:        "(555) 123-4567",
		Company:      "TechCorp Industries",
		Address:      Address{Street: "123 Innovation Way", City: "San Francisco", State: "CA", Zip: "94102", Country: "USA"},
		TotalBilled:  15750.00,
		InvoiceCount: 8,
		CreatedAt:    time.Now().AddDate(0, -6, 0),
	},
	{
		ID:           "cli_002",
		Name:         "Michael Chen",
		Email:        "m.chen@startuplab.co",
		Phone:        "(555) 234-5678",
		Company:      "StartupLab",
		Address:      Address{Street: "456 Venture Blvd", City: "Austin", State: "TX", Zip: "78701", Country: "USA"},
		TotalBilled:  32500.00,
		InvoiceCount: 12,
		CreatedAt:    time.Now().AddDate(0, -8, 0),
	},
	{
		ID:           "cli_003",
		Name:         "Emma Rodriguez",
		Email:        "emma@designstudio.com",
		Phone:        "(555) 345-6789",
		Company:      "Creative Design Studio",
		Address:      Address{Street: "789 Art District", City: "Brooklyn", State: "NY", Zip: "11201", Country: "USA"},
		TotalBilled:  8200.00,
		InvoiceCount: 5,
		CreatedAt:    time.Now().AddDate(0, -3, 0),
	},
	{
		ID:           "cli_004",
		Name:         "James Wilson",
		Email:        "jwilson@financeplus.com",
		Phone:        "(555) 456-7890",
		Company:      "Finance Plus LLC",
		Address:      Address{Street: "321 Wall Street", City: "New York", State: "NY", Zip: "10005", Country: "USA"},
		TotalBilled:  45000.00,
		InvoiceCount: 15,
		CreatedAt:    time.Now().AddDate(-1, 0, 0),
	},
	{
		ID:           "cli_005",
		Name:         "Olivia Martinez",
		Email:        "olivia@greentech.eco",
		Phone:        "(555) 567-8901",
		Company:      "GreenTech Solutions",
		Address:      Address{Street: "555 Eco Park", City: "Portland", State: "OR", Zip: "97201", Country: "USA"},
		TotalBilled:  12300.00,
		InvoiceCount: 6,
		CreatedAt:    time.Now().AddDate(0, -4, 0),
	},
	{
		ID:           "cli_006",
		Name:         "David Park",
		Email:        "david@cloudnine.dev",
		Phone:        "(555) 678-9012",
		Company:      "CloudNine Development",
		Address:      Address{Street: "888 Cloud Ave", City: "Seattle", State: "WA", Zip: "98101", Country: "USA"},
		TotalBilled:  28750.00,
		InvoiceCount: 10,
		CreatedAt:    time.Now().AddDate(0, -5, 0),
	},
}

func GetClient(id string) *Client {
	for _, c := range Clients {
		if c.ID == id {
			return &c
		}
	}
	return nil
}

func SearchClients(query string) []Client {
	if query == "" {
		return Clients
	}
	var results []Client
	for _, c := range Clients {
		// Simple case-insensitive search
		if containsIgnoreCase(c.Name, query) || containsIgnoreCase(c.Company, query) || containsIgnoreCase(c.Email, query) {
			results = append(results, c)
		}
	}
	return results
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && (s[0] == substr[0] || s[0]+32 == substr[0] || s[0] == substr[0]+32)))
}
