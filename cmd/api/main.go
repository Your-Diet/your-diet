package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/constants"
	"github.com/victorgiudicissi/your-diet/internal/handler"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	//"github.com/victorgiudicissi/your-diet/internal/repository"
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

	// // Create MongoDB repository for Diets
	// dietRepo, err := repository.NewDietRepository(mongoURI, dbName) // dbCollection is 'diets'
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB for diets: %v", err)
	// }

	// // Create MongoDB repository for Users (uses 'users' collection by default)
	// userRepo, err := repository.NewMongoUserRepository(mongoURI, dbName)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB for users: %v", err)
	// }

	// Create use cases with repositories
	createDietUseCase := usecase.NewCreateDietUseCase(nil)
	updateDietUseCase := usecase.NewUpdateDietUseCase(nil)
	createUserUseCase := usecase.NewCreateUserUseCase(nil)
	loginUseCase := usecase.NewLoginUseCase(nil)
	listDietsUseCase := usecase.NewListDietsUseCase(nil, nil)

	// Create handlers with use cases
	dietHandler := handler.NewCreateDietHandler(createDietUseCase)
	updateDietHandler := handler.NewUpdateDietHandler(updateDietUseCase)
	registerUserHandler := handler.NewRegisterUserHandler(createUserUseCase)
	userLoginHandler := handler.NewLoginHandler(loginUseCase)
	listDietsHandler := handler.NewListDietsHandler(listDietsUseCase)

	// Set up router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configure CORS
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			"http://localhost:3000":  true,
		}

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}

		c.Next()
	})

	// Set trusted proxies
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Remove trailing slashes from routes
	r.RemoveExtraSlash = true

	// Public routes
	r.GET("/ping", handler.Ping)

	apiGroup := r.Group("/v1")

	// User routes (public)
	userGroup := apiGroup.Group("/users")
	{
		userGroup.POST("", registerUserHandler.Handle)
		userGroup.POST("/login", userLoginHandler.HandleLogin)
	}

	// Diet routes
	dietGroup := apiGroup.Group("/diets")
	dietGroup.Use(middleware.AuthMiddleware([]byte(usecase.JWTSecretKey)))
	{
		dietGroup.POST("", middleware.HasPermission(constants.PermissionCreateDiet), dietHandler.Handle)
		dietGroup.PUT("/:id", middleware.HasPermission(constants.PermissionUpdateDiet), updateDietHandler.Handle)
		dietGroup.GET("", middleware.HasPermission(constants.PermissionListDiet), listDietsHandler.Handle)
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
