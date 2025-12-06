We're going to build a simple invoicing app for managing a list of clients and creating invoices with nice clean shareable invoice pages that show a client the line items they're being billed for, along with a pay now button to pay.

The app is for my own use and not others, so we'll skip authentication for now.  The application needs to be able to run locally in a container using Docker, and in production at Cloudflare using workers, D1 and R2.

## App shell

- Simple logo (icon + text)
- Navigation: Dashboard, Clients, Invoices, Expenses

## Dashboard

- Metrics
- Graph of receivables and expenses as line, bar and pie chart
- Recent invoices
- Recent clients
- Recent expenses
- Actions (create invoice, create client, create expense)

## Clients

- List clients, amount billed
- Create new clients

## Invoices

- List all invoices
- Filter by status
- Create invoice
  - Create line items

## Shareable invoice view

- Clean professional design
- Show client's info
- Line items and total
- Pay Now button launches a payment modal
  - Just mockup the payment form for now.  We wont' wire up payment processing yet.
- Export to PDF option

## Expenses
- List all invoices
- Filter by category
- Create expense
  - Attach receipt
  - Dynamic lists with ability to search and create

# Tech stack

This is a brand new application.  It will use the following:

- Language: Go (Golang) 1.23+
- Web Framework: Echo (v4)
- Templating: Templ (Type-safe HTML generation for Go)
- Interactivity: HTMX (For modals, dynamic table rows, and SPA-like navigation)
- Database: SQLite
- ORM/Queries: SQLC (for type-safe SQL)
- Migrations: Goose (to manage database schema versions)
- Styling: TailwindCSS v4 (using the standalone CLI)
- Typography: Google Fonts (Inter for UI, JetBrains Mono for numbers/code)
- UI Components: FrankenUI (ShadCN compatible CSS-only components)
- PDF Generation: Maroto (Go PDF library)
- Icons: Lucide Icons (SVGs embeded via Templ)
- Dev Tooling: Air (for live reloading during development)
