package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskmaster/internal/config"
	"taskmaster/internal/handler"
	"taskmaster/internal/middleware"
	"taskmaster/internal/repository/postgres"
	"taskmaster/internal/repository/redis"
	"taskmaster/internal/service"
	"taskmaster/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := config.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to postgres: ", err)
	}
	defer db.Close()

	// Redis
	redisClient, err := config.NewRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal("failed to connect to redis: ", err)
	}
	defer redisClient.Close()

	// Repositories (Dependency Injection)
	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	cacheRepo := redis.NewCacheRepository(redisClient)

	// WebSocket Hub
	wsHub := websocket.NewHub(redisClient)
	go wsHub.Run()

	// Services
	authService := service.NewAuthService(userRepo, cacheRepo, cfg.JWTSecret)
	taskService := service.NewTaskService(taskRepo, teamRepo, cacheRepo, wsHub)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	wsHandler := handler.NewWSHandler(wsHub)

	// Router
	r := gin.Default()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(cacheRepo))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
	})

	// Public
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)

	// Protected
	api := r.Group("/api/v1")
	api.Use(middleware.JWT(cfg.JWTSecret))
	api.Use(middleware.RBAC())
	{
		api.POST("/tasks", taskHandler.Create)
		api.GET("/tasks", taskHandler.List)
		api.GET("/tasks/:id", taskHandler.Get)
		api.PUT("/tasks/:id", taskHandler.Update)
		api.DELETE("/tasks/:id", taskHandler.Delete)
		api.POST("/tasks/:id/assign", taskHandler.Assign)
	}

	// WebSocket
	r.GET("/ws", wsHandler.Handle)

	// Server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed: ", err)
		}
	}()

	log.Printf("🚀 TaskMaster running on http://localhost:%s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("forced shutdown: ", err)
	}
	log.Println("✅ Server stopped")
}
