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

	createDietUseCase := usecase.NewCreateDietUseCase(dietRepo)
	updateDietUseCase := usecase.NewUpdateDietUseCase(dietRepo)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)
	listDietsUseCase := usecase.NewListDietsUseCase(dietRepo, userRepo)

	dietHandler := handler.NewCreateDietHandler(createDietUseCase)
	updateDietHandler := handler.NewUpdateDietHandler(updateDietUseCase)
	registerUserHandler := handler.NewRegisterUserHandler(createUserUseCase)
	userLoginHandler := handler.NewLoginHandler(loginUseCase)
	listDietsHandler := handler.NewListDietsHandler(listDietsUseCase)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
