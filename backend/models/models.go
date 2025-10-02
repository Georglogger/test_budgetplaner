package models

import (
	"time"
)

// Budget repräsentiert einen Budgetplan für eine Periode
type Budget struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Status      string    `json:"status"` // draft, approved, closed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
}

// BudgetLine repräsentiert eine flexible Budget-Zeile als Schlüssel-Wert-Paare
type BudgetLine struct {
	ID         string                 `json:"id"`
	BudgetID   string                 `json:"budget_id"`
	Category   string                 `json:"category"`    // z.B. "Personal", "Marketing", "IT"
	Subcategory string                 `json:"subcategory"` // z.B. "Gehälter", "Online-Werbung"
	Amount     float64                `json:"amount"`
	Driver     string                 `json:"driver"`      // Treiber-Name (z.B. "anzahl_mitarbeiter")
	DriverValue float64               `json:"driver_value"` // Wert des Treibers
	Attributes map[string]interface{} `json:"attributes"`  // Flexible zusätzliche Attribute
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// Actual repräsentiert Ist-Daten
type Actual struct {
	ID          string    `json:"id"`
	BudgetID    string    `json:"budget_id"`
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Source      string    `json:"source"` // z.B. "ERP", "Manual", "CSV"
	CreatedAt   time.Time `json:"created_at"`
}

// Scenario ermöglicht What-If-Analysen
type Scenario struct {
	ID          string                 `json:"id"`
	BudgetID    string                 `json:"budget_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"` // Szenario-Parameter (z.B. {"price_increase": 1.05})
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Rule definiert eine treiberbasierte Berechnungsregel
type Rule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Category    string                 `json:"category"`
	Formula     string                 `json:"formula"`     // z.B. "driver_value * unit_cost"
	Parameters  map[string]interface{} `json:"parameters"`  // Regel-Parameter
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
