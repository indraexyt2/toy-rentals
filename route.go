package main

import (
	"final-project/config"
	"github.com/gin-gonic/gin"
)

func setupRoutes(cfg *config.Config, db *config.Database) *gin.Engine {
	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	return r
}
