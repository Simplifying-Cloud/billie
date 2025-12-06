# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Development (runs templ watch, CSS watch, and server in parallel)
make dev

# One-time builds
make templ          # Generate templ files
make css            # Build CSS
make build          # Production build (templ + minified CSS + Go binary)

# Database
make sqlc           # Regenerate SQLC code after modifying queries

# Run server
make run            # Generate + run server
go run ./cmd/server # Direct run (requires templ/css already generated)

# Clean generated files
make clean
```

Server runs on port 3000. Templ proxy runs on port 8080 during development.

## Architecture

This is a Go invoicing application using Echo, Templ, HTMX, Alpine.js, and SQLite.

### Code Structure

- `cmd/server/main.go` - Entry point, route definitions, dependency injection
- `internal/database/` - Database layer
  - `database.go` - SQLite connection with embedded Goose migrations
  - `migrations/` - SQL migration files (schema + seed data)
  - `queries/` - SQLC query definitions
  - `db/` - Generated SQLC code (do not edit)
- `internal/services/` - Business logic layer
  - `clients.go` - Client CRUD + stats
  - `invoices.go` - Invoice CRUD + number generation + status transitions
  - `expenses.go` - Expense CRUD + category/merchant management
  - `metrics.go` - Dashboard aggregations
  - `uploads.go` - Receipt file handling
- `internal/handlers/` - HTTP handlers
  - `handlers.go` - Handler container with DI
  - `api.go` - JSON API endpoints
  - `convert.go` - Service-to-template type converters
  - `dashboard.go`, `clients.go`, `invoices.go`, `expenses.go`, `public.go` - Page handlers
- `internal/render/` - Helper to render Templ components via Echo
- `internal/mock/` - Mock data types (used by templates)
- `views/layouts/` - Base layout templates
- `views/pages/` - Page templates
- `views/components/ui/` - Reusable UI components
- `views/components/icons/` - Lucide icon components
- `views/shared/` - Shared formatting utilities
- `static/css/` - TailwindCSS v4 (input.css -> output.css)
- `static/js/` - Vendored JS (htmx.min.js, alpine.min.js, chart.min.js)

### Database

- **SQLite** with WAL mode for concurrent access
- **SQLC** for type-safe query generation
- **Goose** migrations embedded in binary (auto-run on startup)

To modify the database schema:
1. Create a new migration in `internal/database/migrations/`
2. Update queries in `internal/database/queries/`
3. Run `make sqlc` to regenerate Go code

### Invoice Number Format

Invoices use customer-prefix format: `PREFIX-YYYY-NNNN`
- PREFIX: First 4 characters of company name (uppercase)
- YYYY: Current year
- NNNN: Sequential number per prefix/year

Example: TechCorp Industries -> `TECH-2025-0001`

### Pattern: Handler -> Service -> Database

```
HTTP Request -> Handler -> Service -> SQLC Queries -> SQLite
                  |
                  v
            Templ Template (HTML) or JSON Response
```

### API Endpoints

All JSON APIs are prefixed with `/api/v1/`:

**Clients**
- `GET /api/v1/clients` - List all clients
- `GET /api/v1/clients/:id` - Get client by ID
- `POST /api/v1/clients` - Create client
- `PUT /api/v1/clients/:id` - Update client
- `DELETE /api/v1/clients/:id` - Delete client

**Invoices**
- `GET /api/v1/invoices` - List invoices (optional `?status=` filter)
- `GET /api/v1/invoices/:id` - Get invoice by ID
- `POST /api/v1/invoices` - Create invoice
- `PUT /api/v1/invoices/:id` - Update invoice
- `DELETE /api/v1/invoices/:id` - Delete invoice
- `POST /api/v1/invoices/:id/send` - Mark invoice as sent
- `POST /api/v1/invoices/:id/mark-paid` - Mark invoice as paid

**Expenses**
- `GET /api/v1/expenses` - List expenses (optional `?category=` filter)
- `GET /api/v1/expenses/:id` - Get expense by ID
- `POST /api/v1/expenses` - Create expense
- `PUT /api/v1/expenses/:id` - Update expense
- `DELETE /api/v1/expenses/:id` - Delete expense

**Lookups**
- `GET /api/v1/categories` - List expense categories
- `GET /api/v1/merchants` - List merchants
- `GET /api/v1/merchants/search?q=` - Search merchants

**Metrics**
- `GET /api/v1/metrics/dashboard` - Dashboard metrics and chart data

### HTML Routes

- `/` - Dashboard
- `/clients` - Clients list
- `/clients/:id` - Client detail
- `/clients/search` - Client search (HTMX)
- `/invoices` - Invoices list
- `/invoices/new` - Invoice editor
- `/invoices/:id` - Invoice detail/editor
- `/expenses` - Expenses list
- `/p/:token` - Public shareable invoice view

### Frontend Stack

- **HTMX**: Dynamic interactions (hx-get, hx-post, hx-target, hx-swap)
- **Alpine.js**: Client-side state (sidebar, theme, dropdowns)
- **TailwindCSS v4**: Utility-first CSS with custom fonts
- **Dark mode**: Automatic via prefers-color-scheme with manual toggle

### Templ Components

UI components in `views/components/ui/` accept props structs:

```go
templ Button(props ButtonProps) {
    <button class={ buttonClasses(props) }>
        { children... }
    </button>
}
```

Button variants: `ButtonPrimary`, `ButtonSecondary`, `ButtonGhost`, `ButtonDestructive`, `ButtonOutline`
