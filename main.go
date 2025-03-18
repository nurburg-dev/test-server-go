package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// dbPool is a global variable to hold the database connection pool.
var dbPool *sql.DB
var dbOnce sync.Once

// initDB initializes the database connection pool.
func initDB() {
	dbOnce.Do(func() {
		// Get database connection details from environment variables
		postgresUser := os.Getenv("POSTGRES_USER")
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		postgresDb := os.Getenv("POSTGRES_DB")
		postgresHost := os.Getenv("POSTGRES_HOST")
		postgresPort := os.Getenv("POSTGRES_PORT")

		// Check if all required environment variables are set
		if postgresUser == "" || postgresPassword == "" || postgresDb == "" || postgresHost == "" || postgresPort == "" {
			log.Fatalf("Missing one or more required environment variables: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST, POSTGRES_PORT")
		}

		// Construct the connection string
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDb)

		// Open a database connection
		var err error
		dbPool, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Error opening database connection: %v", err)
		}

		// Set connection pool parameters
		dbPool.SetMaxOpenConns(25)                 // Maximum number of open connections
		dbPool.SetMaxIdleConns(25)                 // Maximum number of idle connections
		dbPool.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

		// Test the database connection
		err = dbPool.Ping()
		if err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}
		log.Println("Database connection pool initialized successfully.")
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection from the pool
	db := dbPool

	// Query the database for all users
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the rows and build a list of users
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error reading data from database", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error reading data from database", http.StatusInternalServerError)
		return
	}

	// Build the response string
	var response strings.Builder
	for _, user := range users {
		response.WriteString(fmt.Sprintf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email))
	}

	// If no users are found, return a message
	if len(users) == 0 {
		response.WriteString("No users found in the database.\n")
	}

	// Write the response
	fmt.Fprintf(w, "%s", response.String())
}

func main() {
	// Initialize the database connection pool
	initDB()
	defer dbPool.Close()

	http.HandleFunc("/", helloHandler)

	fmt.Println("Starting server on port 9000...")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
