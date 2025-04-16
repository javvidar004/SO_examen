package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Country string `json:"country"`
}

func getDBConfig() string {
	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:               dbName,
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	return cfg.FormatDSN()
}

// getItems for request to database
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
		if err := rows.Scan(&item.ID, &item.Name, &item.Country); err != nil {
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

// errorHandler
func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		lastErr := c.Errors.Last()
		log.Printf("API Error: %v", lastErr.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": lastErr.Error()})
	}
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Use the error handler
	router.Use(errorHandler)

	// API endpoints
	router.GET("/items", getItems)

	router.Run(":8080")
}
