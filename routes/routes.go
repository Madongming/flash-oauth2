package routes

import (
	"database/sql"
	"html/template"
	"os"

	"flash-oauth2/config"
	"flash-oauth2/handlers"
	"flash-oauth2/middleware"
	"flash-oauth2/models"
	"flash-oauth2/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Setup(r *gin.Engine, handler *handlers.Handler) {
	// API文档端点
	r.GET("/docs", handler.Documentation)

	// OAuth2端点
	r.GET("/authorize", handler.Authorize)
	r.POST("/token", handler.Token)
	r.POST("/introspect", handler.Introspect)

	// OpenID Connect端点
	r.GET("/userinfo", handler.UserInfo)
	r.GET("/.well-known/jwks.json", handler.JWKs)

	// 用户认证端点
	r.POST("/login", handler.Login)
	r.POST("/send-code", handler.SendVerificationCode)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// SetupAppManagement adds application management routes with admin authentication
func SetupAppManagement(r *gin.Engine, db *sql.DB, redis *redis.Client, cfg *config.Config) {
	appService := services.NewAppManagementService(db)
	appHandler := handlers.NewAppManagementHandler(appService)

	// Create SMS service
	smsService := services.NewSMSService(cfg)

	// Create user service for admin authentication with SMS service
	userService := services.NewUserService(db, redis, smsService)

	// Create main handler for admin auth endpoints
	// We need to get the config for this, but for now we'll create a minimal handler
	// Create admin auth middleware
	adminAuth := middleware.NewAdminAuth(userService)

	// Add custom template functions
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"countActive": func(apps any) int {
			count := 0
			if appSlice, ok := apps.([]*models.ExternalApp); ok {
				for _, app := range appSlice {
					if app.Status == "active" {
						count++
					}
				}
			}
			return count
		},
	})

	// Load HTML templates
	// Try different paths for development and test environments
	if _, err := os.Stat("templates"); err == nil {
		r.LoadHTMLGlob("templates/*")
	} else if _, err := os.Stat("../templates"); err == nil {
		r.LoadHTMLGlob("../templates/*")
	} else {
		// For testing, create minimal templates in memory
		r.SetHTMLTemplate(template.Must(template.New("").Parse(`
			{{define "dashboard.gohtml"}}<!DOCTYPE html><html><head><title>Dashboard</title></head><body><h1>Dashboard</h1></body></html>{{end}}
			{{define "app_details.gohtml"}}<!DOCTYPE html><html><head><title>App Details</title></head><body><h1>App Details</h1></body></html>{{end}}
			{{define "login.gohtml"}}<!DOCTYPE html><html><head><title>Login</title></head><body><h1>Login</h1></body></html>{{end}}
			{{define "admin_login.gohtml"}}<!DOCTYPE html><html><head><title>Admin Login</title></head><body><h1>Admin Login</h1></body></html>{{end}}
			{{define "register_developer.gohtml"}}<!DOCTYPE html><html><head><title>Register Developer</title></head><body><h1>Register Developer</h1></body></html>{{end}}
		`)))
	}

	// We need to create a handler with proper config for admin auth
	// For now, let's create a minimal handler setup
	handler := handlers.New(db, redis, cfg)

	// Admin authentication routes (no auth required)
	r.GET("/admin/login", handler.AdminLogin)
	r.POST("/admin/login", handler.AdminLoginPost)
	r.POST("/admin/logout", handler.AdminLogout)

	// Admin web interface (auth required)
	admin := r.Group("/admin")
	admin.Use(adminAuth.RequireAdmin())
	{
		admin.GET("/dashboard", appHandler.ShowManagementDashboard)
		admin.GET("/apps/:app_id", appHandler.ShowAppDetails)
		admin.GET("/apps/:app_id/keys", appHandler.ShowAppDetails) // Alias for key management

		// Developer management pages
		admin.GET("/developers/new", appHandler.ShowRegisterDeveloper)
		admin.POST("/developers", appHandler.RegisterDeveloperForm)

		// Application management pages
		admin.GET("/apps/new", appHandler.ShowRegisterApp)
	}

	// Admin API endpoints (auth required)
	api := r.Group("/api/admin")
	api.Use(adminAuth.RequireAdmin())
	{
		// Developer management
		api.POST("/developers", appHandler.RegisterDeveloper)
		api.GET("/developers/:developer_id/apps", appHandler.GetDeveloperApps)

		// Application management
		api.POST("/apps", appHandler.RegisterApp)
		api.GET("/apps", appHandler.GetAllApps)

		// Key management
		api.POST("/apps/:app_id/keys", appHandler.GenerateKeyPair)
		api.GET("/apps/:app_id/keys", appHandler.GetAppKeys)
		api.POST("/keys/:key_id/revoke", appHandler.RevokeKey)
	}
}
