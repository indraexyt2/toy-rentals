package main

import (
	"final-project/config"
	"final-project/controller"
	"final-project/middleware"
	"final-project/repository"
	"final-project/service"
	"final-project/utils/helpers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func setupRoutes(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// JWT Konfigurasi
	jwtHelper := helpers.NewJWTHelper(cfg.JWTSecret, cfg.AccessTokenExp, cfg.RefreshTokenExp, cfg.Issuer)

	// User token
	userTokenRepo := repository.NewUserTokenRepository(db)
	userTokenSvc := service.NewTokenService(userTokenRepo, *jwtHelper)

	// Users
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, userTokenRepo, *jwtHelper)
	userController := controller.NewUserController(userSvc, userTokenSvc)

	// Toy category
	toyCategoryRepo := repository.NewToyCategoryRepository(db)
	toyCategorySvc := service.NewToyCategoryService(toyCategoryRepo)
	toyCategoryController := controller.NewToyCategoryController(toyCategorySvc)

	// Toy
	toyRepo := repository.NewToyRepository(db)
	toySvc := service.NewToyService(toyRepo)
	toyController := controller.NewToyController(toySvc)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(*jwtHelper, userTokenSvc)

	// Public routes
	public := r.Group("/api")
	{
		// User routes
		auth := public.Group("/user")
		{
			auth.POST("/auth/register", userController.Insert)
			auth.POST("/auth/login", userController.Login)
		}

		// Toy category routes
		toyCategory := public.Group("/toy")
		{
			toyCategory.GET("/category", toyCategoryController.FindAll)
			toyCategory.GET("/category/:id", toyCategoryController.FinById)
		}

		// Toy routes
		toy := public.Group("/toy")
		{
			toy.GET("", toyController.FindAll)
			toy.GET("/:id", toyController.FinById)
		}
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthMiddleware())
	{
		// User routes
		auth := protected.Group("/user")
		{
			auth.PUT("/auth/:id", userController.UpdateById)
			auth.DELETE("/auth/:id", userController.DeleteById)
			auth.DELETE("/auth/logout", userController.Logout)
		}
	}

	// Admin routes
	admin := r.Group("/api")
	admin.Use(authMiddleware.AdminMiddleware())
	{
		// Admin user routes
		auth := admin.Group("/admin")
		{
			auth.GET("/users", userController.FindAll)
			auth.GET("/user/:id", userController.FinById)
		}

		// Admin toy category routes
		toyCategory := admin.Group("/toy")
		{
			toyCategory.POST("/category", toyCategoryController.Insert)
			toyCategory.PUT("/category/:id", toyCategoryController.UpdateById)
			toyCategory.DELETE("/category/:id", toyCategoryController.DeleteById)
		}

		// Admin toy routes
		toy := admin.Group("/toy")
		{
			toy.POST("", toyController.Insert)
			toy.PUT("/:id", toyController.UpdateById)
			toy.DELETE("/:id", toyController.DeleteById)
		}
	}

	return r
}
