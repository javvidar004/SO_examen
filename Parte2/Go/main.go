package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	mysql "github.com/go-sql-driver/mysql"
)

const (
	dbUser = "root"
	dbPass = "root"
	dbHost = "db"
	dbPort = "3306"
	dbName = "classicmodels"
)

// Item struct represents the data to be retrieved from the database
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

func getDBConfig() string {
	// Construir el DSN usando mysql.Config para mayor legibilidad y seguridad
	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:               dbName,
		AllowNativePasswords: true,
		ParseTime:            true, // ¡Importante para mapear a time.Time!
		Loc:                  time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	return cfg.FormatDSN()
}

// getItems handles the GET request for all items
func getItems(c *gin.Context) {
	db, err := sql.Open("mysql", getDBConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT customerNumber, customerName, city FROM customers") // Replace "items" with your actual table name
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to execute query", "error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Model); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to scan data", "error": err.Error()})
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error during rows iteration", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// getItemByID handles the GET request for a specific item by ID
func getItemByID(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("mysql", getDBConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var item Item
	row := db.QueryRow("SELECT customerNumber, customerName, city FROM customers WHERE customerNumber =?", id) // Replace "items" with your actual table name
	err = row.Scan(&item.ID, &item.Name, &item.Model)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to scan data", "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

// errorHandler middleware for handling errors
func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		lastErr := c.Errors.Last()
		log.Printf("API Error: %v", lastErr.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": lastErr.Error()})
	}
}

func main() {
	// Load environment variables from.env file if it exists
	// This is useful for local development. In production, you might set
	// the DB_CONNECTION_STRING environment variable directly.
	// Note: You might need to install a library like 'github.com/joho/godotenv'
	// and call godotenv.Load() here if you want to use a.env file.
	// For this example, we'll assume the DB_CONNECTION_STRING is already set
	// as an environment variable.

	// Example of how to set the environment variable (for testing purposes only,
	// do not hardcode credentials in production):
	// os.Setenv("DB_CONNECTION_STRING", "your_user:your_password@tcp(127.0.0.1:3306)/your_database?charset=utf8mb4&parseTime=True&loc=Local")

	/*if os.Getenv("DB_CONNECTION_STRING") == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable not set")
		return
	}*/

	router := gin.Default()

	// Use the error handler middleware
	router.Use(errorHandler)

	// Define API endpoints
	router.GET("/items", getItems)
	router.GET("/items/:id", getItemByID)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in environment
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
