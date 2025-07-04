package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victorgiudicissi/your-diet/internal/constants"
	"github.com/victorgiudicissi/your-diet/internal/handler"
	"github.com/victorgiudicissi/your-diet/internal/middleware"
	"github.com/victorgiudicissi/your-diet/internal/repository"
	"github.com/victorgiudicissi/your-diet/internal/usecase"
	"github.com/victorgiudicissi/your-diet/internal/utils"
)

func main() {
	cfg := utils.LoadEnvConfig()

	dietRepo, err := repository.NewDietRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for diets: %v", err)
	}

	userRepo, err := repository.NewMongoUserRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for users: %v", err)
	}

	createDietUseCase := usecase.NewCreateDiet(dietRepo)
	updateDietUseCase := usecase.NewUpdateDiet(dietRepo)
	createUserUseCase := usecase.NewCreateUser(userRepo)
	loginUseCase := usecase.NewLogin(userRepo)
	listDietsUseCase := usecase.NewListDiets(dietRepo, userRepo)
	hub := usecase.NewNotificationHub()

	dietHandler := handler.NewCreateDietHandler(createDietUseCase)
	updateDietHandler := handler.NewUpdateDietHandler(updateDietUseCase)
	registerUserHandler := handler.NewRegisterUserHandler(createUserUseCase)
	userLoginHandler := handler.NewLoginHandler(loginUseCase)
	listDietsHandler := handler.NewListDietsHandler(listDietsUseCase)
	sseHandler := handler.NewSSEHandler(hub)
	notificationHandler := handler.NewNotificationHandler(hub)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:5173": true,
			"http://localhost:3000": true,
			"https://your-diet-frontend-prod-26110891251.us-central1.run.app": true,
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

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.RemoveExtraSlash = true

	r.GET("/ping", handler.Ping)

	apiGroup := r.Group("/v1")

	userGroup := apiGroup.Group("/users")
	{
		userGroup.POST("", registerUserHandler.Handle)
		userGroup.POST("/login", userLoginHandler.HandleLogin)
	}

	dietGroup := apiGroup.Group("/diets")
	dietGroup.Use(middleware.AuthMiddleware([]byte(usecase.JWTSecretKey)))
	{
		dietGroup.POST("", middleware.HasPermission(constants.PermissionCreateDiet), dietHandler.Handle)
		dietGroup.PUT("/:id", middleware.HasPermission(constants.PermissionUpdateDiet), updateDietHandler.Handle)
		dietGroup.GET("", middleware.HasPermission(constants.PermissionListDiet), listDietsHandler.Handle)
	}

	sseGroup := apiGroup.Group("/sse")
	sseGroup.Use(middleware.AuthMiddleware([]byte(usecase.JWTSecretKey)))
	{
		sseGroup.GET("/events", sseHandler.Handle)
		sseGroup.POST("/notify", notificationHandler.Handle)
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
