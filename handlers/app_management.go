// Package handlers provides HTTP handlers for application management endpoints.
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"flash-oauth2/services"

	"github.com/gin-gonic/gin"
)

// AppManagementHandler handles requests for application and developer management
type AppManagementHandler struct {
	appService *services.AppManagementService
}

// NewAppManagementHandler creates a new AppManagementHandler
func NewAppManagementHandler(appService *services.AppManagementService) *AppManagementHandler {
	return &AppManagementHandler{
		appService: appService,
	}
}

// RegisterDeveloper handles developer registration
func (h *AppManagementHandler) RegisterDeveloper(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	developer, err := h.appService.RegisterDeveloper(req.Name, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register developer", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Developer registered successfully",
		"developer": developer,
	})
}

// ShowRegisterDeveloper renders the developer registration form
func (h *AppManagementHandler) ShowRegisterDeveloper(c *gin.Context) {
	c.HTML(http.StatusOK, "register_developer.gohtml", gin.H{
		"title": "Register Developer",
	})
}

// RegisterDeveloperForm handles developer registration from web form
func (h *AppManagementHandler) RegisterDeveloperForm(c *gin.Context) {
	var req struct {
		Name  string `form:"name" binding:"required"`
		Email string `form:"email" binding:"required,email"`
		Phone string `form:"phone"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data: %v", err.Error())
		return
	}

	developer, err := h.appService.RegisterDeveloper(req.Name, req.Email, req.Phone)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to register developer: %v", err.Error())
		return
	}

	// Return success response for AJAX form submission
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Developer registered successfully",
		"developer": developer,
	})
}

// RegisterApp handles external application registration
func (h *AppManagementHandler) RegisterApp(c *gin.Context) {
	var req struct {
		DeveloperID string `json:"developer_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CallbackURL string `json:"callback_url" binding:"required,url"`
		Scopes      string `json:"scopes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if req.Scopes == "" {
		req.Scopes = "openid profile"
	}

	app, err := h.appService.RegisterExternalApp(req.DeveloperID, req.Name, req.Description, req.CallbackURL, req.Scopes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register application", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Application registered successfully",
		"app":     app,
	})
}

// ShowRegisterApp renders the application registration form
func (h *AppManagementHandler) ShowRegisterApp(c *gin.Context) {
	// Get all developers for the dropdown
	developers, err := h.appService.GetAllDevelopers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load developers",
		})
		return
	}

	c.HTML(http.StatusOK, "register_app.gohtml", gin.H{
		"title":      "Register Application",
		"developers": developers,
	})
}

// GenerateKeyPair handles key pair generation for an application
func (h *AppManagementHandler) GenerateKeyPair(c *gin.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App ID is required"})
		return
	}

	var req struct {
		Algorithm string `json:"algorithm"`
		ExpiresIn string `json:"expires_in"` // Duration like "30d", "1y", etc.
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// Use defaults if no body provided
		req.Algorithm = "RS256"
	}

	if req.Algorithm == "" {
		req.Algorithm = "RS256"
	}

	var expiresAt *time.Time
	if req.ExpiresIn != "" {
		duration, err := parseDuration(req.ExpiresIn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expires_in format", "details": err.Error()})
			return
		}
		expTime := time.Now().Add(duration)
		expiresAt = &expTime
	}

	keyPair, err := h.appService.GenerateKeyPair(appID, req.Algorithm, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate key pair", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Key pair generated successfully",
		"key_pair": keyPair,
	})
}

// GetAppKeys retrieves all key pairs for an application
func (h *AppManagementHandler) GetAppKeys(c *gin.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App ID is required"})
		return
	}

	keys, err := h.appService.GetAppKeyPairs(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve keys", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

// RevokeKey revokes a key pair
func (h *AppManagementHandler) RevokeKey(c *gin.Context) {
	keyID := c.Param("key_id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key ID is required"})
		return
	}

	err := h.appService.RevokeKeyPair(keyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke key", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Key revoked successfully",
	})
}

// GetDeveloperApps retrieves all applications for a developer
func (h *AppManagementHandler) GetDeveloperApps(c *gin.Context) {
	developerID := c.Param("developer_id")
	if developerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Developer ID is required"})
		return
	}

	apps, err := h.appService.GetDeveloperApps(developerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apps": apps,
	})
}

// GetAllApps retrieves all applications (admin endpoint)
func (h *AppManagementHandler) GetAllApps(c *gin.Context) {
	apps, err := h.appService.GetAllApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apps": apps,
	})
}

// ShowManagementDashboard displays the key management dashboard
func (h *AppManagementHandler) ShowManagementDashboard(c *gin.Context) {
	apps, err := h.appService.GetAllApps()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load applications",
		})
		return
	}

	developers, err := h.appService.GetAllDevelopers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load developers",
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.gohtml", gin.H{
		"title":      "Application Management Dashboard",
		"apps":       apps,
		"developers": developers,
	})
}

// ShowAppDetails displays detailed information about an application and its keys
func (h *AppManagementHandler) ShowAppDetails(c *gin.Context) {
	appID := c.Param("app_id")

	keys, err := h.appService.GetAppKeyPairs(appID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load application keys",
		})
		return
	}

	c.HTML(http.StatusOK, "app_details.gohtml", gin.H{
		"title":  "Application Details",
		"app_id": appID,
		"keys":   keys,
	})
}

// Helper function to parse duration strings like "30d", "1y", etc.
func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, gin.Error{Err: gin.Error{}}
	}

	unit := s[len(s)-1:]
	value := s[:len(s)-1]

	var multiplier time.Duration
	switch unit {
	case "d":
		multiplier = 24 * time.Hour
	case "w":
		multiplier = 7 * 24 * time.Hour
	case "m":
		multiplier = 30 * 24 * time.Hour
	case "y":
		multiplier = 365 * 24 * time.Hour
	default:
		// Try to parse as standard duration
		return time.ParseDuration(s)
	}

	// Parse the numeric part
	var num int
	if _, err := fmt.Sscanf(value, "%d", &num); err != nil {
		return 0, err
	}

	return time.Duration(num) * multiplier, nil
}
