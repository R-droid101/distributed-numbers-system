package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type IncomingNumber struct {
	Number      int    `json:"number"`
	PublisherID string `json:"publisher_id"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var db *sql.DB

func main() {
	var err error

	if os.Getenv("RUNNING_IN_DOCKER") == "" {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Failed to load env vars from root folder")
		}
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "testuser"),
		getEnv("DB_PASS", "postgres"),
		getEnv("DB_NAME", "numbersdb"),
	)

	if err := waitForDB(connStr); err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Final DB open failed:", err)
	}

	if err := ensureTableExists(db); err != nil {
		log.Fatal("Failed to ensure table exists:", err)
	}

	http.HandleFunc("/consume", consumeHandler)
	log.Println("Consumer listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func consumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data IncomingNumber
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Received: %d from %s", data.Number, data.PublisherID)

	_, err = db.Exec(
		`INSERT INTO published_numbers (number, publisher_id, received_at)
		 VALUES ($1, $2, $3)`,
		data.Number, data.PublisherID, time.Now(),
	)
	if err != nil {
		log.Println("Insert failed:", err)
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, Response{true, "Number received and stored"})
}

func respondJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func ensureTableExists(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS published_numbers (
			id SERIAL PRIMARY KEY,
			number INTEGER NOT NULL,
			publisher_id TEXT NOT NULL,
			received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func waitForDB(connStr string) error {
	const maxRetries = 10
	var err error
	var tryDB *sql.DB

	for i := 1; i <= maxRetries; i++ {
		tryDB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = tryDB.Ping()
			if err == nil {
				log.Println("✅ Connected to DB")
				tryDB.Close()
				return nil
			}
		}
		log.Printf("⏳ Waiting for DB... (%d/%d)", i, maxRetries)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("❌ Failed to connect to DB after %d attempts: %w", maxRetries, err)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
