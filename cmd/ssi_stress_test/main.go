package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	numThreads = 16
	connStr    = "postgres://postgres@localhost:5432/postgres"
	duration   = 5 * time.Minute
)

var tuple_padding = strings.Repeat("a", 100)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for i := 0; i < numThreads; i++ {
		go func() {
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					tx, _ := db.BeginTx(context.Background(), &sql.TxOptions{
						Isolation: sql.LevelSerializable,
					})

					var prev_entry_id, balance int64
					tx.QueryRow("SELECT entry_id, balance FROM entries WHERE account_id = 0 ORDER BY entry_id DESC LIMIT 1").Scan(&prev_entry_id, &balance)
					tx.Exec("INSERT INTO entries (account_id, previous_entry_id, balance, data) VALUES ($1, $2, $3, $4)",
						0, prev_entry_id, balance+1, tuple_padding)
					tx.Commit()
				}
			}
		}()
	}

	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)

	db, _ := sql.Open("postgres", connStr)
	rows, _ := db.Query("SELECT account_id, balance, COUNT(*) FROM entries GROUP BY account_id, balance HAVING COUNT(*) > 1")
	defer rows.Close()
	for rows.Next() {
		var accountID, balance, count int64
		rows.Scan(&accountID, &balance, &count)
		fmt.Printf("Account %d: balance %d appears %d times\n", accountID, balance, count)
	}
}
