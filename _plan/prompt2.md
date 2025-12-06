Let's start with the frontend design. We won't build the full API
and backend functionality just yet. We'll focus on the design and
UI first using Go, Templ, HTMX, and Alpine.js.

Start by creating:
- A base layout component with navigation
- Reusable UI components (buttons, cards, tables, modals, form inputs)
- Stub routes and Templ views for:
  - Dashboard
  - Clients list
  - Invoices list
  - Invoice editor
  - Expenses list
  - Shareable invoice (public-facing)

With your frontend-design kill to design a clean, professional, frontend design, UI and UX for each of those views.  Use hardcoded mock data for now. Add HTMX attributes (hx-get, hx-post, etc.)
as placeholders where server interactions will eventually occur. Use Alpine.js
for client-side UI state (modals, dropdowns, mobile menu toggle).

Requirements:
- All views are mobile responsive (mobile, tablet, desktop)
- Light/dark mode based on system preference (Tailwind's `dark:` variant)
- TailwindCSS v4 utility classes only, no custom CSS
- TailwindCSS v4 color palette only, no custom colors
- Google Fonts
- Lucide Icons (via templ-lucide or inline SVG)

When finished, provide localhost URLs for each view.

Start by asking clarifying questions to inform the design.
dscott@panthro:~/Documents/code/projects/invoicey-opus$ cat _plan/prompt-2.md
Now that we've designed the frontends for:

- Dashboard (Metrics, recent activity)
- Clients list (Table view)
- Invoices list (Table view with status badges)
- Invoice editor (Form with mock inputs for line items)
- Expense list (Table view)
- Shareable invoice (The public-facing view for the client)
ake a plan to complete the full backend implementation for this app.  Refer to @product-overview.md for an overview of functionality that we're going for.

All buttons, navigation links, and functionality should be functional and ready for me to test in a browser.

Ask me clarifying questions to dial in the details for your plan.
