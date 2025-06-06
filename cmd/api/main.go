package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/constants"
	"github.com/victorgiudicissi/your-diet/internal/handler"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
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
	updateDietUseCase := usecase.NewUpdateDietUseCase(dietRepo)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)
	listDietsUseCase := usecase.NewListDietsUseCase(dietRepo)

	// Create handlers with use cases
	dietHandler := handler.NewCreateDietHandler(createDietUseCase)
	updateDietHandler := handler.NewUpdateDietHandler(updateDietUseCase)
	registerUserHandler := handler.NewRegisterUserHandler(createUserUseCase)
	userLoginHandler := handler.NewLoginHandler(loginUseCase)
	listDietsHandler := handler.NewListDietsHandler(listDietsUseCase)

	// Set up router
	r := gin.Default()

	// Set trusted proxies
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Public routes
	r.GET("/ping", handler.Ping)

	apiGroup := r.Group("/v1")

	// User routes (public)
	userGroup := apiGroup.Group("/users")
	{
		userGroup.POST("/", registerUserHandler.Handle)
		userGroup.POST("/login", userLoginHandler.HandleLogin)
	}

	{
		// Diet routes
		dietGroup := apiGroup.Group("/diets")
		{
			dietGroup.Use(middleware.AuthMiddleware([]byte(usecase.JWTSecretKey)))

			dietGroup.POST("/", middleware.HasPermission(constants.PermissionCreateDiet), dietHandler.Handle)
			dietGroup.PUT("/:id", middleware.HasPermission(constants.PermissionUpdateDiet), updateDietHandler.Handle)
			dietGroup.GET("/", middleware.HasPermission(constants.PermissionListDiet), listDietsHandler.Handle)
		}
	}

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
