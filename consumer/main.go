package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
)

type Config struct {
	AuthToken string
	RedisAddr string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var redisClient *redis.Client
var db *sql.DB
var ctx = context.Background()

func main() {
	if os.Getenv("RUNNING_IN_DOCKER") == "" {
		_ = godotenv.Load("../.env")
	}

	config := Config{
		AuthToken: getEnvWithDefault("AUTH_TOKEN", "changeme"),
		RedisAddr: getEnvWithDefault("REDIS_ADDR", "redis:6379"),
		DBHost:    getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:    getEnvWithDefault("DB_PORT", "5432"),
		DBUser:    getEnvWithDefault("DB_USER", "user"),
		DBPass:    getEnvWithDefault("DB_PASS", "password"),
		DBName:    getEnvWithDefault("DB_NAME", "numbersdb"),
	}

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)

	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	go consumeStream()

	log.Println("Consumer started and consuming from Redis stream...")
	select {} // Keep the main goroutine alive
}

func consumeStream() {
	log.Println("🔄 Starting to consume from Redis stream...")
	lastID := "$" // Only new messages after startup

	for {
		streams, err := redisClient.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"numbers-stream", lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Printf("❌ Stream read error: %v", err)
			continue
		}

		if len(streams) == 0 || len(streams[0].Messages) == 0 {
			continue // no new messages, just wait
		}

		for _, stream := range streams {
			for i, message := range stream.Messages {
				numStr := fmt.Sprintf("%v", message.Values["number"])
				pubID := fmt.Sprintf("%v", message.Values["publisher_id"])

				number, err := strconv.Atoi(numStr)
				if err != nil {
					log.Printf("❌ Failed to parse number: %v", err)
					continue
				}

				log.Printf("✅ Consumed: %d from %s", number, pubID)

				_, err = db.Exec(
					`INSERT INTO published_numbers (number, publisher_id, received_at)
					 VALUES ($1, $2, $3)`,
					number, pubID, time.Now(),
				)
				if err != nil {
					log.Printf("❌ Insert failed: %v", err)
				}

				if i == len(stream.Messages)-1 {
					lastID = message.ID // Update only after processing full batch
				}
			}
		}
	}
}

func getEnvWithDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
