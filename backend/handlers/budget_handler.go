package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"budgetplan/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BudgetHandler struct {
	DB *sql.DB
}

// CreateBudget erstellt ein neues Budget
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		// Log the error for debugging
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "JSON parsing failed"})
		return
	}

	budget.ID = uuid.New().String()
	budget.CreatedAt = time.Now()
	budget.UpdatedAt = time.Now()

	query := `INSERT INTO budgets (id, name, description, period_start, period_end, status, created_at, updated_at, created_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := h.DB.Exec(query, budget.ID, budget.Name, budget.Description, budget.PeriodStart,
		budget.PeriodEnd, budget.Status, budget.CreatedAt, budget.UpdatedAt, budget.CreatedBy)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen des Budgets"})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// GetBudgets gibt alle Budgets zurück
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	query := `SELECT id, name, description, period_start, period_end, status, created_at, updated_at, created_by
			  FROM budgets ORDER BY created_at DESC`

	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Budgets"})
		return
	}
	defer rows.Close()

	budgets := []models.Budget{}
	for rows.Next() {
		var budget models.Budget
		err := rows.Scan(&budget.ID, &budget.Name, &budget.Description, &budget.PeriodStart,
			&budget.PeriodEnd, &budget.Status, &budget.CreatedAt, &budget.UpdatedAt, &budget.CreatedBy)
		if err != nil {
			continue
		}
		budgets = append(budgets, budget)
	}

	c.JSON(http.StatusOK, budgets)
}

// GetBudget gibt ein einzelnes Budget zurück
func (h *BudgetHandler) GetBudget(c *gin.Context) {
	id := c.Param("id")

	query := `SELECT id, name, description, period_start, period_end, status, created_at, updated_at, created_by
			  FROM budgets WHERE id = $1`

	var budget models.Budget
	err := h.DB.QueryRow(query, id).Scan(&budget.ID, &budget.Name, &budget.Description,
		&budget.PeriodStart, &budget.PeriodEnd, &budget.Status, &budget.CreatedAt, &budget.UpdatedAt, &budget.CreatedBy)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget nicht gefunden"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen des Budgets"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// UpdateBudget aktualisiert ein Budget
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget.UpdatedAt = time.Now()

	query := `UPDATE budgets SET name = $1, description = $2, period_start = $3, period_end = $4,
			  status = $5, updated_at = $6 WHERE id = $7`

	result, err := h.DB.Exec(query, budget.Name, budget.Description, budget.PeriodStart,
		budget.PeriodEnd, budget.Status, budget.UpdatedAt, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Aktualisieren des Budgets"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget nicht gefunden"})
		return
	}

	budget.ID = id
	c.JSON(http.StatusOK, budget)
}

// DeleteBudget löscht ein Budget
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM budgets WHERE id = $1`
	result, err := h.DB.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Löschen des Budgets"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget nicht gefunden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget erfolgreich gelöscht"})
}

// CreateBudgetLine erstellt eine neue Budget-Zeile
func (h *BudgetHandler) CreateBudgetLine(c *gin.Context) {
	var line models.BudgetLine
	if err := c.ShouldBindJSON(&line); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	line.ID = uuid.New().String()
	line.CreatedAt = time.Now()
	line.UpdatedAt = time.Now()

	attributesJSON, _ := json.Marshal(line.Attributes)

	query := `INSERT INTO budget_lines (id, budget_id, category, subcategory, amount, driver, driver_value, attributes, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := h.DB.Exec(query, line.ID, line.BudgetID, line.Category, line.Subcategory,
		line.Amount, line.Driver, line.DriverValue, attributesJSON, line.CreatedAt, line.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Erstellen der Budget-Zeile"})
		return
	}

	c.JSON(http.StatusCreated, line)
}

// GetBudgetLines gibt alle Budget-Zeilen für ein Budget zurück
func (h *BudgetHandler) GetBudgetLines(c *gin.Context) {
	budgetID := c.Param("id")

	query := `SELECT id, budget_id, category, subcategory, amount, driver, driver_value, attributes, created_at, updated_at
			  FROM budget_lines WHERE budget_id = $1 ORDER BY created_at`

	rows, err := h.DB.Query(query, budgetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Budget-Zeilen"})
		return
	}
	defer rows.Close()

	lines := []models.BudgetLine{}
	for rows.Next() {
		var line models.BudgetLine
		var attributesJSON []byte

		err := rows.Scan(&line.ID, &line.BudgetID, &line.Category, &line.Subcategory,
			&line.Amount, &line.Driver, &line.DriverValue, &attributesJSON, &line.CreatedAt, &line.UpdatedAt)
		if err != nil {
			continue
		}

		if len(attributesJSON) > 0 {
			json.Unmarshal(attributesJSON, &line.Attributes)
		}

		lines = append(lines, line)
	}

	c.JSON(http.StatusOK, lines)
}
