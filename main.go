package main

import (
	"log"

	"showcase_project/config"
	"showcase_project/db"
	"showcase_project/internal/handler"
	"showcase_project/internal/repository"
	"showcase_project/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	err = db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection established")

	// Dependency Injection Setup
	repo := repository.NewRepository(db.DB)
	services := service.NewService(repo, cfg)
	h := handler.NewHandler(services)

	// Set up Gin router
	router := h.InitRoutes()

	// Start the server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
