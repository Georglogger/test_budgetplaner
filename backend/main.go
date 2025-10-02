package main

import (
	"log"
	"os"

	"budgetplan/database"
	"budgetplan/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Datenbank-Konfiguration (aus Umgebungsvariablen oder Defaults)
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "budgetuser"),
		Password: getEnv("DB_PASSWORD", "budgetpass"),
		DBName:   getEnv("DB_NAME", "budgetplan"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Fehler bei der Datenbankverbindung: %v", err)
	}
	defer db.Close()

	// Handlers initialisieren
	budgetHandler := &handlers.BudgetHandler{DB: db}
	importHandler := &handlers.ImportHandler{DB: db}
	reportHandler := &handlers.ReportHandler{DB: db}

	// Router konfigurieren
	router := gin.Default()

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API v1 Routen
	v1 := router.Group("/api/v1")
	{
		// Budget Routen
		v1.POST("/budgets", budgetHandler.CreateBudget)
		v1.GET("/budgets", budgetHandler.GetBudgets)
		v1.GET("/budgets/:id", budgetHandler.GetBudget)
		v1.PUT("/budgets/:id", budgetHandler.UpdateBudget)
		v1.DELETE("/budgets/:id", budgetHandler.DeleteBudget)

		// Budget Lines Routen
		v1.POST("/budget-lines", budgetHandler.CreateBudgetLine)
		v1.GET("/budgets/:id/lines", budgetHandler.GetBudgetLines)

		// Import Routen
		v1.POST("/budgets/:id/import/actuals", importHandler.ImportActualsCSV)
		v1.GET("/budgets/:id/actuals", importHandler.GetActuals)

		// Report Routen
		v1.GET("/budgets/:id/reports/plan-actual", reportHandler.GetPlanActualReport)
		v1.GET("/budgets/:id/reports/category-summary", reportHandler.GetCategorySummary)
	}

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server startet auf Port %s...", port)
	router.Run(":" + port)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
