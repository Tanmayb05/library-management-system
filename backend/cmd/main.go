package main

import (
	"encoding/json"
	"library-management/internal/database"
	"library-management/internal/handlers"
	"library-management/internal/logger"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize structured logging
	log := logger.GetLogger()

	log.WithFields(map[string]interface{}{
		"component": "main",
	}).Info("Starting Library Management System")

	// Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	log.Info("Database connection established successfully")

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(db)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Book routes
	api.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	api.HandleFunc("/books", bookHandler.GetAllBooks).Methods("GET")
	api.HandleFunc("/books/{id:[0-9]+}", bookHandler.GetBook).Methods("GET")
	api.HandleFunc("/books/{id:[0-9]+}", bookHandler.UpdateBook).Methods("PUT")
	api.HandleFunc("/books/{id:[0-9]+}", bookHandler.DeleteBook).Methods("DELETE")

	// Enhanced health check endpoint
	router.HandleFunc("/health", createHealthHandler(db)).Methods("GET")

	// Setup CORS with configurable origins
	allowedOrigins := getAllowedOrigins()
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})

	log.WithField("origins", allowedOrigins).Info("CORS configuration set")

	handler := c.Handler(router)

	// Get server port from environment
	port := getEnv("PORT", "8080")

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.WithFields(map[string]interface{}{
		"port":         port,
		"read_timeout": "15s",
		"write_timeout": "15s",
		"idle_timeout": "60s",
	}).Info("Starting HTTP server")

	log.WithFields(map[string]interface{}{
		"health_check": "http://localhost:" + port + "/health",
		"api_base":     "http://localhost:" + port + "/api/v1",
	}).Info("Server endpoints configured")

	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("Server failed to start")
	}
}

// getAllowedOrigins returns the list of allowed CORS origins from environment variables
func getAllowedOrigins() []string {
	// Get allowed origins from environment variable
	originsEnv := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")

	// Split by comma and trim whitespace
	origins := strings.Split(originsEnv, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}

// createHealthHandler creates an enhanced health check handler that includes database connectivity
func createHealthHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()

		log.WithField("endpoint", "/health").Debug("Health check requested")

		// Basic health status
		health := map[string]interface{}{
			"status":    "healthy",
			"service":   "library-management",
			"timestamp": time.Now().UTC(),
		}

		// Check database connectivity
		dbStatus := "healthy"
		dbError := ""

		if err := db.Ping(); err != nil {
			dbStatus = "unhealthy"
			dbError = err.Error()
			health["status"] = "degraded"

			log.WithError(err).Error("Database health check failed")
		}

		health["database"] = map[string]interface{}{
			"status": dbStatus,
			"error":  dbError,
		}

		// Set appropriate HTTP status
		statusCode := http.StatusOK
		if health["status"] == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		} else if health["status"] == "degraded" {
			statusCode = http.StatusOK // Still available but degraded
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(health); err != nil {
			log.WithError(err).Error("Failed to encode health response")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.WithFields(map[string]interface{}{
			"status":     health["status"],
			"db_status":  dbStatus,
			"status_code": statusCode,
		}).Info("Health check completed")
	}
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
