package main

import (
	"fmt"
	"log"
	"net/http"
	"rewardpage/router"
	"rewardpage/service"

	"github.com/rs/cors"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Initialize database connection
	err = service.InitializeDB()
	if err != nil {
		log.Panic("Failed to initialize database:", err)
	}

	// CHANGE: Initialize daily task service with MongoDB connection
	// This sets up the daily task checklist system and TTL indexes
	service.InitDailyTaskService(service.GetDB())

	// Application entry point
	fmt.Println("MongoDB Api")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router.Router())

	// Start the server
	fmt.Println("Starting server on :4000...")
	log.Fatal(http.ListenAndServe(":4000", handler))
	fmt.Println("Server started on port 4000...")
}
