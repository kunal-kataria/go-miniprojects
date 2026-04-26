package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Read connection string from environment.
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {  // Stop if the required env var is missing.
		log.Fatal("DATABASE_URL variable is not set")
	}

	// Prepare connection state and retry limits.
	var conn *pgx.Conn
	var err error
	maxRetries := 5

	// Try to connect with exponential backoff.
	for i := 0; i < maxRetries; i++ {
		conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}
		// Log failures and wait before retrying.
		log.Printf("Connection failed: %v. Retrying...", err)
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
	if err != nil {
		// Abort if all attempts fail.
		log.Fatalf("Unable to connect to database after %d attempts: %v\n", maxRetries, err)
	}

	// Confirm the connection succeeded.
	fmt.Println("Successfully connected to database")
	// Ensure the connection is closed on exit.
	defer conn.Close(context.Background())
	// Signal that the program is finishing; deferred close runs next.
	fmt.Println("Exiting; connection will close now")
}