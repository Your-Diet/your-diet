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

	// Create MongoDB repository for Diets
	dietRepo, err := repository.NewDietRepository(mongoURI, dbName) // dbCollection is 'diets'
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for diets: %v", err)
	}

	// Create MongoDB repository for Users (uses 'users' collection by default)
	userRepo, err := repository.NewMongoUserRepository(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for users: %v", err)
	}

	// Create use cases with repositories
	createDietUseCase := usecase.NewCreateDietUseCase(dietRepo)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)

	// Create handlers with use cases
	dietHandler := handler.NewCreateDietHandler(createDietUseCase)
	registerUserHandler := handler.NewRegisterUserHandler(createUserUseCase)
	userLoginHandler := handler.NewLoginHandler(loginUseCase)

	// Set up router
	r := gin.Default()

	// Set trusted proxies
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Routes
	r.GET("/ping", handler.Ping)

	dietGroup := r.Group("/v1/diets")

	dietGroup.POST("/", dietHandler.Handle)

	userGroup := r.Group("/v1/users")

	userGroup.POST("/", registerUserHandler.Handle)      
	userGroup.POST("/login", userLoginHandler.HandleLogin) 

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
