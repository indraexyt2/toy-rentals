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

	// Uaers
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, userTokenRepo, *jwtHelper)
	userController := controller.NewUserController(userSvc, userTokenSvc)

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
	}

	return r
}
