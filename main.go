// Package main implements a complete OAuth2 + OpenID Connect authorization server.
//
// This server provides:
//   - OAuth2 Authorization Code Flow
//   - OpenID Connect (OIDC) support
//   - JWT access tokens and ID tokens with RSA signature
//   - Phone number + verification code authentication
//   - Automatic user registration/login (idempotent operations)
//   - PostgreSQL for persistent data storage
//   - Redis for session and temporary data storage
//
// The server implements all standard OAuth2 endpoints including authorization,
// token exchange, user info, token refresh, and token introspection.
//
// Usage:
//
//	go run main.go
//
// Environment variables:
//   - PORT: Server port (default: 8080)
//   - DATABASE_URL: PostgreSQL connection string
//   - REDIS_URL: Redis connection string
package main

import (
	"flash-oauth2/config"
	"flash-oauth2/database"
	"flash-oauth2/handlers"
	"flash-oauth2/middleware"
	"flash-oauth2/redis_client"
	"flash-oauth2/routes"
	"log"

	"github.com/gin-gonic/gin"
)

// main initializes and starts the OAuth2 authorization server.
// It sets up database connections, Redis client, middleware, routes,
// and starts the HTTP server on the configured port.
func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 运行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Redis
	redisClient := redis_client.Init(cfg.RedisURL)

	// 创建Gin引擎
	r := gin.Default()

	// 设置中间件
	r.Use(middleware.CORS())

	// 创建处理器
	handler := handlers.New(db, redisClient, cfg)

	// 设置路由
	routes.Setup(r, handler)

	// 设置应用管理路由
	routes.SetupAppManagement(r, db, redisClient, cfg)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Management Dashboard: http://localhost:%s/admin/dashboard", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
