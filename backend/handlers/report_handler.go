package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	DB *sql.DB
}

type PlanActualComparison struct {
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory"`
	Planned     float64 `json:"planned"`
	Actual      float64 `json:"actual"`
	Variance    float64 `json:"variance"`
	VariancePct float64 `json:"variance_pct"`
}

// GetPlanActualReport gibt einen Plan-Ist-Vergleich zurück
func (h *ReportHandler) GetPlanActualReport(c *gin.Context) {
	budgetID := c.Param("id")

	query := `
		SELECT
			COALESCE(bl.category, a.category) as category,
			COALESCE(bl.subcategory, a.subcategory) as subcategory,
			COALESCE(SUM(bl.amount), 0) as planned,
			COALESCE(SUM(a.amount), 0) as actual
		FROM budget_lines bl
		FULL OUTER JOIN actuals a ON bl.budget_id = a.budget_id
			AND bl.category = a.category
			AND COALESCE(bl.subcategory, '') = COALESCE(a.subcategory, '')
		WHERE COALESCE(bl.budget_id, a.budget_id) = $1
		GROUP BY COALESCE(bl.category, a.category), COALESCE(bl.subcategory, a.subcategory)
		ORDER BY category, subcategory
	`

	rows, err := h.DB.Query(query, budgetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen des Berichts"})
		return
	}
	defer rows.Close()

	comparisons := []PlanActualComparison{}
	for rows.Next() {
		var comp PlanActualComparison
		err := rows.Scan(&comp.Category, &comp.Subcategory, &comp.Planned, &comp.Actual)
		if err != nil {
			continue
		}

		comp.Variance = comp.Actual - comp.Planned
		if comp.Planned != 0 {
			comp.VariancePct = (comp.Variance / comp.Planned) * 100
		}

		comparisons = append(comparisons, comp)
	}

	c.JSON(http.StatusOK, comparisons)
}

type CategorySummary struct {
	Category string  `json:"category"`
	Planned  float64 `json:"planned"`
	Actual   float64 `json:"actual"`
	Variance float64 `json:"variance"`
}

// GetCategorySummary gibt eine Zusammenfassung nach Kategorien zurück
func (h *ReportHandler) GetCategorySummary(c *gin.Context) {
	budgetID := c.Param("id")

	query := `
		SELECT
			COALESCE(bl.category, a.category) as category,
			COALESCE(SUM(bl.amount), 0) as planned,
			COALESCE(SUM(a.amount), 0) as actual
		FROM budget_lines bl
		FULL OUTER JOIN actuals a ON bl.budget_id = a.budget_id AND bl.category = a.category
		WHERE COALESCE(bl.budget_id, a.budget_id) = $1
		GROUP BY COALESCE(bl.category, a.category)
		ORDER BY category
	`

	rows, err := h.DB.Query(query, budgetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen der Zusammenfassung"})
		return
	}
	defer rows.Close()

	summaries := []CategorySummary{}
	for rows.Next() {
		var summary CategorySummary
		err := rows.Scan(&summary.Category, &summary.Planned, &summary.Actual)
		if err != nil {
			continue
		}

		summary.Variance = summary.Actual - summary.Planned
		summaries = append(summaries, summary)
	}

	c.JSON(http.StatusOK, summaries)
}
