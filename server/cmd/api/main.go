package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorayuth/task-manager-go/server/internal/config"
	"github.com/sorayuth/task-manager-go/server/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Database ready!!")

	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Println("listening on http://localhost:" + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
