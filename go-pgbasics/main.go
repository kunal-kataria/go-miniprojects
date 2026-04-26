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

var preparedFindBooksByAuthor bool

// Function to demonstrate full column selection
func getAllBookInfo(ctx context.Context, conn *pgx.Conn) error {
    books, err := conn.Query(ctx, "Select * from books")
	if err != nil {
		return err
	}
	defer books.Close()

	for books.Next() {
		var id, book_id int
		var title, author, publications string

		err = books.Scan(&id, &book_id, &title, &author, &publications)
		if err != nil {
			return err
		}

		fmt.Printf("Book: %d, %d, %s, %s, %s\n", id, book_id, title, author, publications)
	}

	return nil
}

// Function to demonstrate optimized column selection
func getBookBasics(ctx context.Context, conn *pgx.Conn) error {
    books, err := conn.Query(ctx, "Select title,author from books")
	if err != nil {
		return err
	}
	defer books.Close()

	for books.Next() {
		var title, author string

		err = books.Scan(&title, &author)
		if err != nil {
			return err
		}

		fmt.Printf("Book: %s, %s\n", title, author)
	}
    return books.Err()
}

func findBooksByAuthor(ctx context.Context, conn *pgx.Conn, author string) error {
	stmt := "findBooksByAuthor"
    
	if !preparedFindBooksByAuthor {
		_, err := conn.Prepare(ctx, stmt, "Select title, publications from books where author=$1")
		if err != nil {
			return err
		}
		preparedFindBooksByAuthor = true
	}

	var title, publications string

	book, err := conn.Query(ctx, stmt, author)
	if err != nil {
		return err
	}

	for book.Next() {
		err = book.Scan(&title, &publications)
		if err != nil {
			return err
		}
		fmt.Printf("Book: %s, %s, %s\n", title, author, publications)
	}
	
    return nil
}

func main() {
	connString := os.Getenv("DATABASE_URL")
	// set DATABASE_URL=postgresql://username:password@localhost:5432/myDatabse
	if connString == "" {
		log.Fatal("DATABASE_URL variable is not set")
	}
	var conn *pgx.Conn
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}
		log.Printf("Connection failed: %v. Retrying...", err)
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
	if err != nil {
		log.Fatalf("Unable to connect to database after %d attempts: %v\n", maxRetries, err)
	}
	fmt.Printf("Successfully connected to the database...")

	// Add these tests before the connection is closed:
    fmt.Println("\nTesting full column selection:")
    if err := getAllBookInfo(context.Background(), conn); err != nil {
        log.Printf("Error getting all book info: %v\n", err)
    }

	fmt.Println("\nTesting optimized column selection:")
    if err := getBookBasics(context.Background(), conn); err != nil {
        log.Printf("Error getting all book info: %v\n", err)
    }

	fmt.Println("\nTesting prepared statement:")
    if err := findBooksByAuthor(context.Background(), conn, "GRR"); err != nil {
        log.Printf("Error finding books by author: %v\n", err)
    }

	defer conn.Close(context.Background())
}
