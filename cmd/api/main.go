package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/handler"
	"github.com/victorgiudicissi/your-diet/internal/repository"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
)

func main() {
	// Initialize MongoDB repository
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "your-diet"
	}

	dbCollection := os.Getenv("MONGO_COLLECTION")
	if dbCollection == "" {
		dbCollection = "diets"
	}

	// Create MongoDB repository
	dietRepo, err := repository.NewDietRepository(mongoURI, dbName, dbCollection)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Create use case with repository
	createDietUseCase := usecase.NewCreateDietUseCase(dietRepo)

	// Create handler with use case
	dietHandler := handler.NewCreateDietHandler(createDietUseCase)

	// Set up router
	r := gin.Default()

	// Set trusted proxies
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Routes
	r.GET("/ping", handler.Ping)
	r.POST("/diet", dietHandler.Handle)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
