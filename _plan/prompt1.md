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
  - Export invoice as PDF

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
