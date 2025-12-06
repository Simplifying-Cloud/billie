# Billie

A modern invoicing application built with Go, featuring a clean UI and real-time interactions.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)
![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-4-06B6D4?style=flat&logo=tailwindcss)

## Features

- **Client Management** - Track clients with contact info and billing history
- **Invoice Creation** - Create professional invoices with line items and automatic calculations
- **Expense Tracking** - Log expenses by category with receipt uploads
- **Dashboard** - Real-time metrics and charts for revenue and expenses
- **Public Invoice Links** - Share invoices via unique URLs for client viewing
- **Dark Mode** - Automatic theme detection with manual toggle
- **JSON API** - Full REST API for external integrations

## Tech Stack

- **Backend**: Go + Echo framework
- **Database**: SQLite with SQLC (type-safe queries)
- **Templates**: Templ (type-safe HTML)
- **Frontend**: HTMX + Alpine.js + TailwindCSS v4
- **Migrations**: Goose (embedded in binary)

## Quick Start

### Prerequisites

- Go 1.24+
- [Templ](https://templ.guide/) CLI (`go install github.com/a-h/templ/cmd/templ@latest`)
- [TailwindCSS](https://tailwindcss.com/) standalone CLI (download from releases)

### Installation

```bash
# Clone the repository
git clone https://github.com/dscott/billie.git
cd billie

# Install Go dependencies
go mod download

# Build and run
make build
./bin/server
```

The server starts at http://localhost:3000

### Development

```bash
# Run with live reload (watches templ, CSS, and Go files)
make dev
```

## Project Structure

```
billie/
├── cmd/server/          # Application entry point
├── internal/
│   ├── database/        # SQLite + SQLC + migrations
│   ├── services/        # Business logic layer
│   ├── handlers/        # HTTP handlers (HTML + JSON API)
│   ├── render/          # Templ rendering helper
│   └── mock/            # Data types for templates
├── views/
│   ├── layouts/         # Base page layouts
│   ├── pages/           # Page templates
│   └── components/      # Reusable UI components
├── static/
│   ├── css/             # TailwindCSS
│   └── js/              # HTMX, Alpine.js, Chart.js
└── uploads/             # Receipt file storage
```

## API Endpoints

All JSON APIs are prefixed with `/api/v1/`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/clients` | List all clients |
| POST | `/api/v1/clients` | Create a client |
| GET | `/api/v1/invoices` | List invoices |
| POST | `/api/v1/invoices` | Create an invoice |
| POST | `/api/v1/invoices/:id/send` | Mark invoice as sent |
| POST | `/api/v1/invoices/:id/mark-paid` | Mark invoice as paid |
| GET | `/api/v1/expenses` | List expenses |
| POST | `/api/v1/expenses` | Create an expense |
| GET | `/api/v1/metrics/dashboard` | Dashboard metrics |

See [CLAUDE.md](CLAUDE.md) for full API documentation.

## Invoice Numbers

Invoices automatically generate numbers using customer-prefix format:

```
PREFIX-YYYY-NNNN
```

- **PREFIX**: First 4 characters of company name (uppercase)
- **YYYY**: Current year
- **NNNN**: Sequential number per prefix/year

Example: TechCorp Industries → `TECH-2025-0001`

## Database

The application uses SQLite with embedded migrations. The database file (`billie.db`) is created automatically on first run with seed data including:

- 6 sample clients
- 7 sample invoices
- 8 sample expenses
- Expense categories and merchants

## Make Targets

```bash
make dev          # Development with live reload
make build        # Production build
make run          # Build and run
make templ        # Generate templ files
make css          # Build CSS
make sqlc         # Regenerate SQLC code
make clean        # Clean generated files
```

## License

MIT
