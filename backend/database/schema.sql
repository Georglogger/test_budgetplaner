-- Budgets Tabelle
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)
);

-- Budget Lines Tabelle (flexible Struktur)
CREATE TABLE budget_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    category VARCHAR(255) NOT NULL,
    subcategory VARCHAR(255),
    amount DECIMAL(15, 2) NOT NULL,
    driver VARCHAR(255),
    driver_value DECIMAL(15, 4),
    attributes JSONB, -- Flexible Attribute als JSON
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Actuals Tabelle (Ist-Daten)
CREATE TABLE actuals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    category VARCHAR(255) NOT NULL,
    subcategory VARCHAR(255),
    amount DECIMAL(15, 2) NOT NULL,
    date DATE NOT NULL,
    source VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scenarios Tabelle
CREATE TABLE scenarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parameters JSONB, -- Szenario-Parameter als JSON
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rules Tabelle (treiberbasierte Berechnungsregeln)
CREATE TABLE rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    formula TEXT NOT NULL,
    parameters JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indizes für Performance
CREATE INDEX idx_budget_lines_budget_id ON budget_lines(budget_id);
CREATE INDEX idx_budget_lines_category ON budget_lines(category);
CREATE INDEX idx_actuals_budget_id ON actuals(budget_id);
CREATE INDEX idx_actuals_date ON actuals(date);
CREATE INDEX idx_scenarios_budget_id ON scenarios(budget_id);
CREATE INDEX idx_rules_category ON rules(category);
