// Package handlers provides HTTP request handlers for OAuth2 and OpenID Connect endpoints.
// It implements the complete OAuth2 Authorization Server specification including
// authorization, token exchange, user authentication, and OpenID Connect features.
package handlers

import (
	"database/sql"
	"flash-oauth2/config"
	"flash-oauth2/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Handler contains all the service dependencies needed for OAuth2 operations.
// It acts as a container for business logic services and configuration.
type Handler struct {
	userService  *services.UserService  // User management and authentication
	oauthService *services.OAuthService // OAuth2 core operations
	jwtService   *services.JWTService   // JWT token operations
	smsService   services.SMSService    // SMS service for testing access
	config       *config.Config         // Server configuration
}

// New creates a new Handler instance with all required dependencies.
// It initializes all service layers and returns a configured handler
// ready to process OAuth2 requests.
//
// Parameters:
//   - db: Database connection for persistent storage
//   - redis: Redis client for temporary data (verification codes, sessions)
//   - cfg: Server configuration including cryptographic keys
//
// Returns:
//   - *Handler: Configured handler with all services initialized
func New(db *sql.DB, redis *redis.Client, cfg *config.Config) *Handler {
	// 创建SMS服务
	smsService := services.NewSMSService(cfg)

	// 创建用户服务，传入SMS服务
	userService := services.NewUserService(db, redis, smsService)
	oauthService := services.NewOAuthService(db)
	jwtService := services.NewJWTService(cfg.JWTPrivateKey, cfg.JWTPublicKey, "flash-oauth2")

	return &Handler{
		userService:  userService,
		oauthService: oauthService,
		jwtService:   jwtService,
		smsService:   smsService,
		config:       cfg,
	}
}

// Documentation serves the API documentation as HTML
// This endpoint provides a web-accessible version of the complete API documentation
//
// GET /docs
//
// Returns:
//   - HTML page with comprehensive API documentation
//   - Includes all endpoints, models, and usage examples
//
// Example:
//
//	curl http://localhost:8080/docs
func (h *Handler) Documentation(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, getAPIDocumentationHTML())
}

// getAPIDocumentationHTML returns the complete API documentation as HTML
func getAPIDocumentationHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flash OAuth2 Server - API Documentation</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1, h2, h3 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h1 { font-size: 2.5em; }
        h2 { font-size: 2em; margin-top: 40px; }
        h3 { font-size: 1.5em; margin-top: 30px; }
        code {
            background: #f8f9fa;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        }
        pre {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 6px;
            padding: 15px;
            overflow-x: auto;
            margin: 15px 0;
        }
        .endpoint {
            background: #e8f5e8;
            border-left: 4px solid #28a745;
            padding: 15px;
            margin: 15px 0;
            border-radius: 0 6px 6px 0;
        }
        .method {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: bold;
            color: white;
            margin-right: 10px;
        }
        .get { background: #007bff; }
        .post { background: #28a745; }
        .nav {
            background: #343a40;
            padding: 15px;
            margin: -30px -30px 30px -30px;
            border-radius: 8px 8px 0 0;
        }
        .nav a {
            color: #fff;
            text-decoration: none;
            margin-right: 20px;
            padding: 5px 10px;
            border-radius: 4px;
            transition: background 0.3s;
        }
        .nav a:hover {
            background: #495057;
        }
        .badge {
            display: inline-block;
            padding: 3px 8px;
            background: #6c757d;
            color: white;
            border-radius: 12px;
            font-size: 0.8em;
            margin-left: 5px;
        }
        .warning {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 6px;
            padding: 15px;
            margin: 15px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="nav">
            <a href="#overview">Overview</a>
            <a href="#endpoints">API Endpoints</a>
            <a href="#authentication">Authentication</a>
            <a href="#models">Data Models</a>
            <a href="#examples">Examples</a>
            <a href="/health">Health Check</a>
        </div>

        <h1 id="overview">Flash OAuth2 Server API Documentation</h1>
        
        <div class="warning">
            <strong>Live Documentation:</strong> This documentation is served by the running OAuth2 server instance.
            All endpoints listed below are available on this server.
        </div>

        <p>Flash OAuth2 is a complete OAuth2 Authorization Server with OpenID Connect support, featuring:</p>
        <ul>
            <li><strong>OAuth2 Authorization Code Flow</strong> with PKCE support</li>
            <li><strong>OpenID Connect (OIDC)</strong> implementation</li>
            <li><strong>JWT tokens</strong> with RSA asymmetric encryption</li>
            <li><strong>Phone-based authentication</strong> with verification codes</li>
            <li><strong>Automatic user registration</strong> (idempotent operations)</li>
        </ul>

        <h2 id="endpoints">API Endpoints</h2>

        <h3>OAuth2 Core Endpoints</h3>

        <div class="endpoint">
            <span class="method get">GET</span>
            <strong>/authorize</strong>
            <span class="badge">OAuth2</span>
            <p>Initiates the OAuth2 authorization flow. Redirects user to login if not authenticated.</p>
            <strong>Parameters:</strong>
            <ul>
                <li><code>client_id</code> (required): OAuth2 client identifier</li>
                <li><code>redirect_uri</code> (required): Callback URL after authorization</li>
                <li><code>response_type</code> (required): Must be "code"</li>
                <li><code>scope</code> (optional): Requested scopes (space-separated)</li>
                <li><code>state</code> (optional): Client state parameter</li>
            </ul>
        </div>

        <div class="endpoint">
            <span class="method post">POST</span>
            <strong>/token</strong>
            <span class="badge">OAuth2</span>
            <p>Exchanges authorization code for access tokens or refreshes existing tokens.</p>
            <strong>Content-Type:</strong> <code>application/x-www-form-urlencoded</code><br>
            <strong>Parameters:</strong>
            <ul>
                <li><code>grant_type</code>: "authorization_code" or "refresh_token"</li>
                <li><code>code</code>: Authorization code (for auth code flow)</li>
                <li><code>client_id</code>: OAuth2 client identifier</li>
                <li><code>client_secret</code>: OAuth2 client secret</li>
                <li><code>redirect_uri</code>: Must match original request</li>
            </ul>
        </div>

        <div class="endpoint">
            <span class="method post">POST</span>
            <strong>/introspect</strong>
            <span class="badge">OAuth2</span>
            <p>Validates and returns metadata about an access token.</p>
        </div>

        <h3>OpenID Connect Endpoints</h3>

        <div class="endpoint">
            <span class="method get">GET</span>
            <strong>/userinfo</strong>
            <span class="badge">OIDC</span>
            <p>Returns user profile information for authenticated requests.</p>
            <strong>Authorization:</strong> <code>Bearer {access_token}</code>
        </div>

        <div class="endpoint">
            <span class="method get">GET</span>
            <strong>/.well-known/jwks.json</strong>
            <span class="badge">OIDC</span>
            <p>Returns JSON Web Key Set for token verification.</p>
        </div>

        <h3 id="authentication">Authentication Endpoints</h3>

        <div class="endpoint">
            <span class="method post">POST</span>
            <strong>/send-code</strong>
            <span class="badge">Auth</span>
            <p>Sends SMS verification code to phone number.</p>
            <strong>Body (JSON):</strong>
            <pre>{"phone": "13800138000"}</pre>
        </div>

        <div class="endpoint">
            <span class="method post">POST</span>
            <strong>/login</strong>
            <span class="badge">Auth</span>
            <p>Verifies code and completes user authentication.</p>
            <strong>Body (JSON):</strong>
            <pre>{"phone": "13800138000", "code": "123456"}</pre>
        </div>

        <h2 id="models">Data Models</h2>

        <h3>Token Response</h3>
        <pre>{
  "access_token": "eyJhbGciOiJSUzI1NiI...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_token_string",
  "id_token": "eyJhbGciOiJSUzI1NiI...",
  "scope": "openid profile"
}</pre>

        <h3>UserInfo Response</h3>
        <pre>{
  "sub": "123",
  "phone": "13800138000"
}</pre>

        <h3>JWKS Response</h3>
        <pre>{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig", 
      "kid": "oauth2-rsa-key",
      "n": "...",
      "e": "AQAB"
    }
  ]
}</pre>

        <h2 id="examples">Quick Examples</h2>

        <h3>Complete OAuth2 Flow</h3>
        <ol>
            <li><strong>Redirect to Authorization:</strong>
                <pre>GET /authorize?client_id=my-app&redirect_uri=https://app.com/callback&response_type=code&scope=openid</pre>
            </li>
            <li><strong>User Login (if needed):</strong>
                <pre>POST /send-code
{"phone": "13800138000"}

POST /login  
{"phone": "13800138000", "code": "123456"}</pre>
            </li>
            <li><strong>Exchange Code for Tokens:</strong>
                <pre>POST /token
grant_type=authorization_code&code=ABC123&client_id=my-app&client_secret=secret&redirect_uri=https://app.com/callback</pre>
            </li>
            <li><strong>Access User Info:</strong>
                <pre>GET /userinfo
Authorization: Bearer eyJhbGciOiJSUzI1NiI...</pre>
            </li>
        </ol>

        <h3>Test with cURL</h3>
        <pre># Send verification code
curl -X POST http://localhost:8080/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone": "13800138000"}'

# Health check  
curl http://localhost:8080/health

# Get JWKS
curl http://localhost:8080/.well-known/jwks.json</pre>

        <div class="warning">
            <strong>Development Note:</strong> In development, verification codes are printed to the server console. 
            In production, integrate with an SMS service for real code delivery.
        </div>

        <footer style="margin-top: 50px; padding-top: 20px; border-top: 1px solid #dee2e6; color: #6c757d;">
            <p>Generated by Flash OAuth2 Server • 
            <a href="https://tools.ietf.org/html/rfc6749" target="_blank">OAuth2 RFC 6749</a> • 
            <a href="https://openid.net/connect/" target="_blank">OpenID Connect</a></p>
        </footer>
    </div>
</body>
</html>`
}

// AdminLogin shows the admin login page
// This endpoint provides a login form specifically for admin users
//
// GET /admin/login
//
// Returns:
//   - HTML login form for admin authentication
func (h *Handler) AdminLogin(c *gin.Context) {
	// Check if already logged in
	session := getAdminSession(c)
	if userID, exists := session["user_id"]; exists {
		if userID != "" {
			// Already logged in, redirect to dashboard
			c.Redirect(302, "/admin/dashboard")
			return
		}
	}

	// Show login form
	c.HTML(200, "admin_login.gohtml", gin.H{
		"error": c.Query("error"),
	})
}

// AdminLoginPost handles admin login form submission
// This endpoint processes admin authentication using phone + verification code
//
// POST /admin/login
//
// Form parameters:
//   - phone: Admin's phone number
//   - code: Verification code
//
// Returns:
//   - Redirect to dashboard on success
//   - Redirect back to login with error on failure
func (h *Handler) AdminLoginPost(c *gin.Context) {
	phone := c.PostForm("phone")
	code := c.PostForm("code")

	if phone == "" || code == "" {
		c.Redirect(302, "/admin/login?error=请填写完整的登录信息")
		return
	}

	// Verify the code using existing user service
	user, err := h.userService.VerifyCode(phone, code)
	if err != nil {
		c.Redirect(302, "/admin/login?error="+err.Error())
		return
	}

	// Check if user has admin role
	if user.Role != "admin" {
		c.Redirect(302, "/admin/login?error=您没有管理员权限，无法访问管理平台")
		return
	}

	// Set admin session
	setAdminSession(c, "user_id", fmt.Sprintf("%d", user.ID))

	// Redirect to dashboard
	c.Redirect(302, "/admin/dashboard")
}

// AdminLogout handles admin logout
// This endpoint clears the admin session and redirects to login
//
// POST /admin/logout
//
// Returns:
//   - Redirect to login page
func (h *Handler) AdminLogout(c *gin.Context) {
	clearAdminSession(c)
	c.Redirect(302, "/admin/login?message=已安全退出")
}

// GetSMSService returns the SMS service instance (for testing)
func (h *Handler) GetSMSService() services.SMSService {
	return h.smsService
}

// Helper functions for admin session management
func getAdminSession(c *gin.Context) map[string]interface{} {
	session := make(map[string]interface{})

	// Get user_id from cookie
	if userID, err := c.Cookie("admin_user_id"); err == nil && userID != "" {
		session["user_id"] = userID
	}

	return session
}

func setAdminSession(c *gin.Context, key string, value interface{}) {
	switch key {
	case "user_id":
		c.SetCookie("admin_user_id", value.(string), 86400, "/admin", "", false, true) // 24 hours, HttpOnly
	}
}

func clearAdminSession(c *gin.Context) {
	c.SetCookie("admin_user_id", "", -1, "/admin", "", false, true)
}
