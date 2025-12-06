-- +goose Up

-- Seed expense categories
INSERT INTO expense_categories (id, name, slug) VALUES
    ('cat_001', 'Office Supplies', 'office'),
    ('cat_002', 'Travel', 'travel'),
    ('cat_003', 'Software & Tools', 'software'),
    ('cat_004', 'Marketing', 'marketing'),
    ('cat_005', 'Utilities', 'utilities'),
    ('cat_006', 'Other', 'other');

-- Seed merchants
INSERT INTO merchants (id, name) VALUES
    ('mer_001', 'Adobe'),
    ('mer_002', 'Amazon'),
    ('mer_003', 'American Airlines'),
    ('mer_004', 'Apple'),
    ('mer_005', 'AT&T'),
    ('mer_006', 'Comcast'),
    ('mer_007', 'Delta Airlines'),
    ('mer_008', 'Figma'),
    ('mer_009', 'GitHub'),
    ('mer_010', 'Google'),
    ('mer_011', 'Hilton Hotels'),
    ('mer_012', 'LinkedIn'),
    ('mer_013', 'Marriott'),
    ('mer_014', 'Microsoft'),
    ('mer_015', 'Notion'),
    ('mer_016', 'Slack'),
    ('mer_017', 'Southwest Airlines'),
    ('mer_018', 'Staples'),
    ('mer_019', 'Stripe'),
    ('mer_020', 'United Airlines'),
    ('mer_021', 'Verizon'),
    ('mer_022', 'Zoom'),
    ('mer_023', 'The Capital Grille');

-- Seed clients
INSERT INTO clients (id, name, email, phone, company, address_street, address_city, address_state, address_zip, address_country, created_at) VALUES
    ('cli_001', 'Sarah Johnson', 'sarah@techcorp.io', '(555) 123-4567', 'TechCorp Industries', '123 Innovation Way', 'San Francisco', 'CA', '94102', 'USA', datetime('now', '-6 months')),
    ('cli_002', 'Michael Chen', 'm.chen@startuplab.co', '(555) 234-5678', 'StartupLab', '456 Venture Blvd', 'Austin', 'TX', '78701', 'USA', datetime('now', '-8 months')),
    ('cli_003', 'Emma Rodriguez', 'emma@designstudio.com', '(555) 345-6789', 'Creative Design Studio', '789 Art District', 'Brooklyn', 'NY', '11201', 'USA', datetime('now', '-3 months')),
    ('cli_004', 'James Wilson', 'jwilson@financeplus.com', '(555) 456-7890', 'Finance Plus LLC', '321 Wall Street', 'New York', 'NY', '10005', 'USA', datetime('now', '-1 year')),
    ('cli_005', 'Olivia Martinez', 'olivia@greentech.eco', '(555) 567-8901', 'GreenTech Solutions', '555 Eco Park', 'Portland', 'OR', '97201', 'USA', datetime('now', '-4 months')),
    ('cli_006', 'David Park', 'david@cloudnine.dev', '(555) 678-9012', 'CloudNine Development', '888 Cloud Ave', 'Seattle', 'WA', '98101', 'USA', datetime('now', '-5 months'));

-- Seed invoice sequences (for existing invoices)
INSERT INTO invoice_sequences (client_prefix, year, last_number) VALUES
    ('INV', 2024, 7);

-- Seed invoices
INSERT INTO invoices (id, number, client_id, status, issue_date, due_date, subtotal, tax_rate, tax_amount, total, notes, public_token) VALUES
    ('inv_001', 'INV-2024-001', 'cli_001', 'paid', date('now', '-1 month', '-15 days'), date('now', '-1 month', '+15 days'), 6500.00, 8.25, 536.25, 7036.25, 'Thank you for your business!', 'abc123'),
    ('inv_002', 'INV-2024-002', 'cli_002', 'sent', date('now', '-10 days'), date('now', '+20 days'), 19000.00, 8.25, 1567.50, 20567.50, 'Payment due within 30 days.', 'def456'),
    ('inv_003', 'INV-2024-003', 'cli_003', 'draft', date('now'), date('now', '+30 days'), 4250.00, 8.25, 350.63, 4600.63, '', 'ghi789'),
    ('inv_004', 'INV-2024-004', 'cli_004', 'overdue', date('now', '-2 months'), date('now', '-1 month'), 21600.00, 8.25, 1782.00, 23382.00, 'OVERDUE - Please remit payment immediately.', 'jkl012'),
    ('inv_005', 'INV-2024-005', 'cli_005', 'paid', date('now', '-20 days'), date('now', '+10 days'), 4500.00, 8.25, 371.25, 4871.25, 'Thank you for choosing eco-friendly solutions!', 'mno345'),
    ('inv_006', 'INV-2024-006', 'cli_006', 'sent', date('now', '-5 days'), date('now', '+25 days'), 13500.00, 8.25, 1113.75, 14613.75, '', 'pqr678'),
    ('inv_007', 'INV-2024-007', 'cli_002', 'draft', date('now'), date('now', '+30 days'), 15000.00, 8.25, 1237.50, 16237.50, '', 'stu901');

-- Seed line items
INSERT INTO line_items (id, invoice_id, description, quantity, rate, amount, sort_order) VALUES
    -- Invoice 1
    ('li_001', 'inv_001', 'Website Redesign', 1, 5000.00, 5000.00, 0),
    ('li_002', 'inv_001', 'SEO Optimization', 10, 150.00, 1500.00, 1),
    -- Invoice 2
    ('li_003', 'inv_002', 'Mobile App Development - Phase 1', 1, 12000.00, 12000.00, 0),
    ('li_004', 'inv_002', 'UI/UX Design', 40, 125.00, 5000.00, 1),
    ('li_005', 'inv_002', 'Project Management', 20, 100.00, 2000.00, 2),
    -- Invoice 3
    ('li_006', 'inv_003', 'Brand Identity Package', 1, 3500.00, 3500.00, 0),
    ('li_007', 'inv_003', 'Logo Design Revisions', 3, 250.00, 750.00, 1),
    -- Invoice 4
    ('li_008', 'inv_004', 'Financial Dashboard Development', 1, 15000.00, 15000.00, 0),
    ('li_009', 'inv_004', 'Data Integration Services', 1, 5000.00, 5000.00, 1),
    ('li_010', 'inv_004', 'Training Sessions', 8, 200.00, 1600.00, 2),
    -- Invoice 5
    ('li_011', 'inv_005', 'Sustainability Report Design', 1, 2500.00, 2500.00, 0),
    ('li_012', 'inv_005', 'Infographic Creation', 5, 400.00, 2000.00, 1),
    -- Invoice 6
    ('li_013', 'inv_006', 'Cloud Migration Consulting', 40, 200.00, 8000.00, 0),
    ('li_014', 'inv_006', 'AWS Setup & Configuration', 1, 3000.00, 3000.00, 1),
    ('li_015', 'inv_006', 'Security Audit', 1, 2500.00, 2500.00, 2),
    -- Invoice 7
    ('li_016', 'inv_007', 'Mobile App Development - Phase 2', 1, 15000.00, 15000.00, 0);

-- Seed expenses
INSERT INTO expenses (id, description, category_id, merchant_id, amount, date, receipt_path, notes) VALUES
    ('exp_001', 'Adobe Creative Cloud Annual', 'cat_003', 'mer_001', 599.88, date('now', '-15 days'), '', 'Annual subscription renewal'),
    ('exp_002', 'Client Lunch Meeting - TechCorp', 'cat_004', 'mer_023', 127.50, date('now', '-12 days'), '', 'Business development lunch with Sarah Johnson'),
    ('exp_003', 'Figma Pro Subscription', 'cat_003', 'mer_008', 15.00, date('now', '-10 days'), '', 'Monthly design tool'),
    ('exp_004', 'Office Supplies - Amazon', 'cat_001', 'mer_002', 89.43, date('now', '-8 days'), '', 'Notebooks, pens, sticky notes'),
    ('exp_005', 'Flight to Austin - Client Visit', 'cat_002', 'mer_017', 342.00, date('now', '-5 days'), '', 'Round trip for StartupLab project kickoff'),
    ('exp_006', 'Hotel - Austin', 'cat_002', 'mer_011', 189.00, date('now', '-4 days'), '', 'One night stay'),
    ('exp_007', 'Internet Service - Monthly', 'cat_005', 'mer_006', 79.99, date('now', '-3 days'), '', 'High-speed fiber connection'),
    ('exp_008', 'GitHub Pro', 'cat_003', 'mer_009', 7.00, date('now', '-1 day'), '', 'Monthly subscription');

-- +goose Down
DELETE FROM expenses;
DELETE FROM line_items;
DELETE FROM invoices;
DELETE FROM invoice_sequences;
DELETE FROM clients;
DELETE FROM merchants;
DELETE FROM expense_categories;
