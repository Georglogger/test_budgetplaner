package handlers

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"budgetplan/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImportHandler struct {
	DB *sql.DB
}

// ImportActualsCSV importiert Ist-Daten aus einer CSV-Datei
// CSV Format: category,subcategory,amount,date
func (h *ImportHandler) ImportActualsCSV(c *gin.Context) {
	budgetID := c.Param("id")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keine Datei hochgeladen"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fehler beim Lesen der CSV-Datei"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV-Datei ist leer oder enthält keine Daten"})
		return
	}

	// Überspringe Header-Zeile
	records = records[1:]

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Starten der Transaktion"})
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO actuals (id, budget_id, category, subcategory, amount, date, source, created_at)
							 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Vorbereiten der Abfrage"})
		return
	}
	defer stmt.Close()

	importedCount := 0
	var errors []string

	for i, record := range records {
		if len(record) < 4 {
			errors = append(errors, "Zeile "+strconv.Itoa(i+2)+": Unvollständige Daten")
			continue
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			errors = append(errors, "Zeile "+strconv.Itoa(i+2)+": Ungültiger Betrag")
			continue
		}

		date, err := time.Parse("2006-01-02", record[3])
		if err != nil {
			errors = append(errors, "Zeile "+strconv.Itoa(i+2)+": Ungültiges Datum (Format: YYYY-MM-DD)")
			continue
		}

		actual := models.Actual{
			ID:          uuid.New().String(),
			BudgetID:    budgetID,
			Category:    record[0],
			Subcategory: record[1],
			Amount:      amount,
			Date:        date,
			Source:      "CSV",
			CreatedAt:   time.Now(),
		}

		_, err = stmt.Exec(actual.ID, actual.BudgetID, actual.Category, actual.Subcategory,
			actual.Amount, actual.Date, actual.Source, actual.CreatedAt)

		if err != nil {
			errors = append(errors, "Zeile "+strconv.Itoa(i+2)+": Fehler beim Einfügen")
			continue
		}

		importedCount++
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Speichern der Daten"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imported": importedCount,
		"total":    len(records),
		"errors":   errors,
	})
}

// GetActuals gibt alle Ist-Daten für ein Budget zurück
func (h *ImportHandler) GetActuals(c *gin.Context) {
	budgetID := c.Param("id")

	query := `SELECT id, budget_id, category, subcategory, amount, date, source, created_at
			  FROM actuals WHERE budget_id = $1 ORDER BY date DESC`

	rows, err := h.DB.Query(query, budgetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Ist-Daten"})
		return
	}
	defer rows.Close()

	actuals := []models.Actual{}
	for rows.Next() {
		var actual models.Actual
		err := rows.Scan(&actual.ID, &actual.BudgetID, &actual.Category, &actual.Subcategory,
			&actual.Amount, &actual.Date, &actual.Source, &actual.CreatedAt)
		if err != nil {
			continue
		}
		actuals = append(actuals, actual)
	}

	c.JSON(http.StatusOK, actuals)
}
