// Package middleware provides HTTP middleware functions for request processing.
// This includes authentication, authorization, and request filtering middleware.
package middleware

import (
	"flash-oauth2/models"
	"flash-oauth2/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminAuth creates a middleware that requires admin authentication for management platform access.
// It checks for a valid admin session and redirects to login if not authenticated.
type AdminAuth struct {
	userService *services.UserService
}

// NewAdminAuth creates a new admin authentication middleware
func NewAdminAuth(userService *services.UserService) *AdminAuth {
	return &AdminAuth{
		userService: userService,
	}
}

// RequireAdmin middleware function that checks if the user is authenticated and has admin role
func (a *AdminAuth) RequireAdmin() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check if user is already authenticated via session
		session := getSession(c)
		userIDStr, exists := session["user_id"]
		if !exists {
			// Not logged in, redirect to admin login
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr.(string))
		if err != nil {
			// Invalid user ID, clear session and redirect to login
			clearSession(c)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// Get user from database
		user, err := a.userService.GetUserByID(userID)
		if err != nil {
			// User not found, clear session and redirect to login
			clearSession(c)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// Check if user has admin role
		if user.Role != "admin" {
			// Not an admin, return forbidden
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. Admin privileges required.",
			})
			c.Abort()
			return
		}

		// User is authenticated and is an admin, set user context
		c.Set("user", user)
		c.Next()
	})
}

// Simple session management using cookies (in production, use proper session store like Redis)
func getSession(c *gin.Context) map[string]interface{} {
	session := make(map[string]interface{})

	// Get user_id from cookie
	if userID, err := c.Cookie("admin_user_id"); err == nil && userID != "" {
		session["user_id"] = userID
	}

	return session
}

func setSession(c *gin.Context, key string, value interface{}) {
	switch key {
	case "user_id":
		c.SetCookie("admin_user_id", value.(string), 86400, "/admin", "", false, true) // 24 hours, HttpOnly
	}
}

func clearSession(c *gin.Context) {
	c.SetCookie("admin_user_id", "", -1, "/admin", "", false, true)
}

// Helper function to set user session
func SetUserSession(c *gin.Context, userID int) {
	setSession(c, "user_id", strconv.Itoa(userID))
}

// Helper function to get current admin user
func GetCurrentAdminUser(c *gin.Context) *models.User {
	if user, exists := c.Get("user"); exists {
		if adminUser, ok := user.(*models.User); ok {
			return adminUser
		}
	}
	return nil
}
