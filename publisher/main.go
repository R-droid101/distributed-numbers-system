package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	PublisherID string
	StartNumber int
	EndNumber   int
	AuthToken   string
	RedisAddr   string
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type NumbersResponse struct {
	PublisherID string `json:"publisher_id"`
	Numbers     []int  `json:"numbers"`
}

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	publisherID := getEnvWithDefault("PUBLISHER_ID", "publisher-1")
	startNumber := getEnvAsIntWithDefault("START_NUMBER", 1)
	endNumber := getEnvAsIntWithDefault("END_NUMBER", 10)
	authToken := getEnvWithDefault("AUTH_TOKEN", "changeme")
	redisAddr := getEnvWithDefault("REDIS_ADDR", "redis:6379")
	port := getEnvWithDefault("PORT", "8080")

	config := Config{
		PublisherID: publisherID,
		StartNumber: startNumber,
		EndNumber:   endNumber,
		AuthToken:   authToken,
		RedisAddr:   redisAddr,
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", config.RedisAddr, err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", createHealthHandler(config)).Methods("GET")
	router.Handle("/publish", authMiddleware(config, http.HandlerFunc(createPublishHandler(config)))).Methods("POST")

	log.Printf("Publisher %s starting. Range: %d-%d", publisherID, startNumber, endNumber)
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func createPublishHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numbers := make([]int, 0, config.EndNumber-config.StartNumber+1)
		for num := config.StartNumber; num <= config.EndNumber; num++ {
			numbers = append(numbers, num)

			msg := map[string]interface{}{
				"number":       num,
				"publisher_id": config.PublisherID,
				"timestamp":    time.Now().Format(time.RFC3339),
			}

			if err := redisClient.XAdd(ctx, &redis.XAddArgs{
				Stream: "numbers-stream",
				Values: msg,
			}).Err(); err != nil {
				log.Printf("❌ Failed to push number %d to stream: %v", num, err)
			}
		}

		data := NumbersResponse{
			PublisherID: config.PublisherID,
			Numbers:     numbers,
		}

		respondWithJSON(w, http.StatusOK, Response{
			Success: true,
			Message: fmt.Sprintf("Pushed %d numbers to stream", len(numbers)),
			Data:    data,
		})
	}
}

func createHealthHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, Response{
			Success: true,
			Message: fmt.Sprintf("Publisher %s is healthy", config.PublisherID),
			Data: map[string]interface{}{
				"publisher_id": config.PublisherID,
				"range": map[string]int{
					"start": config.StartNumber,
					"end":   config.EndNumber,
				},
			},
		})
	}
}

func authMiddleware(config Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") || strings.TrimPrefix(header, "Bearer ") != config.AuthToken {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized: Invalid or missing token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{
		Success: false,
		Message: message,
	})
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
