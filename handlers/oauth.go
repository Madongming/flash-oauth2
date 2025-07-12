// Package handlers provides HTTP request handlers for OAuth2 and OpenID Connect endpoints.
package handlers

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthorizeRequest represents the parameters for an OAuth2 authorization request.
// It follows RFC 6749 (OAuth 2.0) specification for authorization endpoint parameters.
type AuthorizeRequest struct {
	ResponseType string `form:"response_type" binding:"required"` // Must be "code" for authorization code flow
	ClientID     string `form:"client_id" binding:"required"`     // Client identifier
	RedirectURI  string `form:"redirect_uri" binding:"required"`  // Client redirect URI
	Scope        string `form:"scope"`                            // Requested scopes (optional)
	State        string `form:"state"`                            // Opaque value to prevent CSRF attacks
}

// TokenRequest represents the parameters for an OAuth2 token request.
// It supports both authorization_code and refresh_token grant types.
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"` // "authorization_code" or "refresh_token"
	Code         string `form:"code"`                          // Authorization code (for authorization_code grant)
	RedirectURI  string `form:"redirect_uri"`                  // Must match authorization request
	ClientID     string `form:"client_id"`                     // Client identifier
	ClientSecret string `form:"client_secret"`                 // Client secret for authentication
	RefreshToken string `form:"refresh_token"`                 // Refresh token (for refresh_token grant)
}

// LoginRequest represents the parameters for user authentication.
// Users authenticate using phone number and verification code.
type LoginRequest struct {
	Phone string `form:"phone" binding:"required"` // User's phone number
	Code  string `form:"code" binding:"required"`  // 6-digit verification code
}

// SendCodeRequest represents the parameters for sending verification codes.
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"` // Phone number to send verification code to
}

// TokenResponse represents the OAuth2 token response format.
// It follows RFC 6749 specification for token endpoint responses.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`            // JWT access token
	TokenType    string `json:"token_type"`              // Always "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // Token lifetime in seconds
	RefreshToken string `json:"refresh_token,omitempty"` // Long-lived refresh token
	IDToken      string `json:"id_token,omitempty"`      // OpenID Connect ID token
	Scope        string `json:"scope,omitempty"`         // Granted scopes
}

// UserInfoResponse represents the OpenID Connect UserInfo endpoint response.
// It contains user profile information that can be accessed with an access token.
type UserInfoResponse struct {
	Sub   string `json:"sub"`   // Subject identifier (user ID)
	Phone string `json:"phone"` // User's phone number
}

// Authorize handles OAuth2 authorization requests (RFC 6749 Section 4.1.1).
// This endpoint initiates the authorization code flow by:
//  1. Validating the client and redirect URI
//  2. Checking if the user is authenticated
//  3. Displaying login form if not authenticated
//  4. Creating authorization code if authenticated
//  5. Redirecting back to client with authorization code
//
// Supported parameters:
//   - response_type: Must be "code"
//   - client_id: Registered client identifier
//   - redirect_uri: Must match registered URI
//   - scope: Requested permissions (optional)
//   - state: CSRF protection token (recommended)
//
// Example:
//
//	GET /authorize?response_type=code&client_id=123&redirect_uri=https://client.com/callback&state=xyz
func (h *Handler) Authorize(c *gin.Context) {
	var req AuthorizeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	// 验证客户端
	client, err := h.oauthService.GetClient(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	// 验证重定向URI
	validRedirectURI := false
	for _, uri := range client.RedirectURIs {
		if uri == req.RedirectURI {
			validRedirectURI = true
			break
		}
	}
	if !validRedirectURI {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_redirect_uri"})
		return
	}

	// 验证响应类型
	if req.ResponseType != "code" {
		redirectURL := fmt.Sprintf("%s?error=unsupported_response_type&state=%s", req.RedirectURI, req.State)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// 检查用户是否已登录
	userIDStr, exists := c.Get("user_id")
	if !exists {
		// 用户未登录，显示登录页面
		c.HTML(http.StatusOK, "login.gohtml", gin.H{
			"client_id":     req.ClientID,
			"redirect_uri":  req.RedirectURI,
			"scope":         req.Scope,
			"state":         req.State,
			"response_type": req.ResponseType,
		})
		return
	}

	userID := userIDStr.(int)

	// 创建授权码
	scope := req.Scope
	if scope == "" {
		scope = client.Scope
	}

	authCode, err := h.oauthService.CreateAuthCode(req.ClientID, userID, req.RedirectURI, scope)
	if err != nil {
		redirectURL := fmt.Sprintf("%s?error=server_error&state=%s", req.RedirectURI, req.State)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// 重定向到客户端
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, authCode.Code, req.State)
	c.Redirect(http.StatusFound, redirectURL)
}

// Login handles user authentication using phone number and verification code.
// This endpoint provides idempotent user operations - if a user doesn't exist,
// they are automatically registered; if they exist, they are logged in.
//
// The process:
//  1. Validates phone number and verification code
//  2. Creates user account if it doesn't exist
//  3. Sets user session
//  4. If OAuth2 parameters are present, creates authorization code and redirects
//  5. Otherwise returns login success response
//
// Parameters:
//   - phone: User's phone number
//   - code: 6-digit verification code
//   - client_id: OAuth2 client ID (optional, for OAuth2 flow)
//   - redirect_uri: OAuth2 redirect URI (optional, for OAuth2 flow)
//   - scope: Requested scopes (optional, for OAuth2 flow)
//   - state: CSRF protection (optional, for OAuth2 flow)
//
// Example:
//
//	POST /login
//	Content-Type: application/x-www-form-urlencoded
//	phone=13800138000&code=123456
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	// 验证验证码并获取用户
	user, err := h.userService.VerifyCode(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_verification_code", "error_description": err.Error()})
		return
	}

	// 设置用户会话
	c.Set("user_id", user.ID)

	// 检查是否有OAuth2参数
	clientID := c.PostForm("client_id")
	redirectURI := c.PostForm("redirect_uri")
	scope := c.PostForm("scope")
	state := c.PostForm("state")

	if clientID != "" && redirectURI != "" {
		// 创建授权码
		authCode, err := h.oauthService.CreateAuthCode(clientID, user.ID, redirectURI, scope)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}

		// 重定向到客户端
		redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, authCode.Code, state)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user_id": user.ID,
	})
}

// SendVerificationCode sends a verification code to the specified phone number.
// The verification code is stored in Redis with a 5-minute expiration.
//
// In a production environment, this would integrate with an SMS service
// to send the code via text message. For development, the code is printed
// to the server console.
//
// Parameters:
//   - phone: Phone number to send verification code to
//
// Example:
//
//	POST /send-code
//	Content-Type: application/json
//	{"phone": "13800138000"}
//
// Response:
//
//	{"message": "verification code sent"}
func (h *Handler) SendVerificationCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	err := h.userService.SendVerificationCode(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent"})
}

// Token handles OAuth2 token requests (RFC 6749 Section 4.1.3).
// This endpoint supports multiple grant types:
//   - authorization_code: Exchange authorization code for access token
//   - refresh_token: Use refresh token to get new access token
//
// For authorization_code grant:
//   - Validates authorization code
//   - Issues JWT access token and refresh token
//   - Optionally issues OpenID Connect ID token
//
// For refresh_token grant:
//   - Validates refresh token
//   - Issues new JWT access token
//   - Optionally issues new ID token
//
// All tokens are signed with RSA private key and can be verified using
// the public key available at /.well-known/jwks.json
//
// Example:
//
//	POST /token
//	Content-Type: application/x-www-form-urlencoded
//	grant_type=authorization_code&code=ABC123&redirect_uri=https://client.com/callback&client_id=123&client_secret=secret
func (h *Handler) Token(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	switch req.GrantType {
	case "authorization_code":
		h.handleAuthorizationCodeGrant(c, req)
	case "refresh_token":
		h.handleRefreshTokenGrant(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
	}
}

func (h *Handler) handleAuthorizationCodeGrant(c *gin.Context, req TokenRequest) {
	// 验证客户端
	client, err := h.oauthService.ValidateClient(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}

	// 交换授权码
	authCode, err := h.oauthService.ExchangeAuthCode(req.Code, req.ClientID, req.RedirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": err.Error()})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(authCode.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 生成JWT访问令牌
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, client.ID, authCode.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 生成刷新令牌
	refreshToken, err := h.oauthService.CreateRefreshToken(client.ID, user.ID, authCode.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小时
		RefreshToken: refreshToken.Token,
		Scope:        authCode.Scope,
	}

	// 如果请求包含openid scope，生成ID令牌
	if strings.Contains(authCode.Scope, "openid") {
		idToken, err := h.jwtService.GenerateIDToken(user, client.ID)
		if err == nil {
			response.IDToken = idToken
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) handleRefreshTokenGrant(c *gin.Context, req TokenRequest) {
	// 验证客户端
	client, err := h.oauthService.ValidateClient(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}

	// 验证刷新令牌
	refreshToken, err := h.oauthService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": err.Error()})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(refreshToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 生成新的访问令牌
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, client.ID, refreshToken.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	response := TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1小时
		Scope:       refreshToken.Scope,
	}

	// 如果请求包含openid scope，生成新的ID令牌
	if strings.Contains(refreshToken.Scope, "openid") {
		idToken, err := h.jwtService.GenerateIDToken(user, client.ID)
		if err == nil {
			response.IDToken = idToken
		}
	}

	c.JSON(http.StatusOK, response)
}

// UserInfo handles OpenID Connect UserInfo requests (OpenID Connect Core 1.0 Section 5.3).
// This endpoint returns user profile information for the authenticated user.
// The access token must be provided in the Authorization header.
//
// The endpoint:
//  1. Extracts and validates the Bearer token from Authorization header
//  2. Verifies the JWT signature and claims
//  3. Retrieves user information from the database
//  4. Returns user profile data
//
// Authentication:
//
//	Authorization: Bearer <access_token>
//
// Example:
//
//	GET /userinfo
//	Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
//
// Response:
//
//	{
//	  "sub": "123",
//	  "phone": "13800138000"
//	}
func (h *Handler) UserInfo(c *gin.Context) {
	// 从Authorization header获取访问令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	// 提取Bearer令牌
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	accessToken := tokenParts[1]

	// 验证JWT令牌
	claims, err := h.jwtService.ParseAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token", "error_description": err.Error()})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	response := UserInfoResponse{
		Sub:   strconv.Itoa(user.ID),
		Phone: user.Phone,
	}

	c.JSON(http.StatusOK, response)
}

// JWKs handles JSON Web Key Set requests (RFC 7517).
// This endpoint provides the public keys used to verify JWT signatures.
// Clients can use these keys to verify access tokens and ID tokens
// without contacting the authorization server.
//
// The endpoint returns a JWKS (JSON Web Key Set) containing:
//   - Key type (RSA)
//   - Usage (signature)
//   - Algorithm (RS256)
//   - Key ID
//   - Public key material
//
// Example:
//
//	GET /.well-known/jwks.json
//
// Response:
//
//	{
//	  "keys": [
//	    {
//	      "kty": "RSA",
//	      "use": "sig",
//	      "alg": "RS256",
//	      "kid": "default",
//	      "n": "..."
//	    }
//	  ]
//	}
func (h *Handler) JWKs(c *gin.Context) {
	publicKey := h.jwtService.GetPublicKey()

	// Convert RSA public key to JWK format
	c.JSON(http.StatusOK, gin.H{
		"keys": []gin.H{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": "default",
				"n":   base64URLEncode(publicKey.N.Bytes()),
				"e":   base64URLEncode(big.NewInt(int64(publicKey.E)).Bytes()),
			},
		},
	})
}

// Introspect handles OAuth2 token introspection requests (RFC 7662).
// This endpoint allows clients to determine the active state and metadata
// of an access token. It accepts both active and inactive tokens.
//
// The endpoint:
//  1. Extracts token from form parameters
//  2. Validates and parses the JWT token
//  3. Returns token metadata if valid, or active=false if invalid
//
// Parameters:
//   - token: The access token to introspect
//
// Example:
//
//	POST /introspect
//	Content-Type: application/x-www-form-urlencoded
//	token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
//
// Response (active token):
//
//	{
//	  "active": true,
//	  "sub": 123,
//	  "client_id": "default-client",
//	  "scope": "openid profile",
//	  "exp": 1640995200,
//	  "iat": 1640991600,
//	  "iss": "flash-oauth2",
//	  "aud": "default-client"
//	}
//
// Response (inactive token):
//
//	{
//	  "active": false
//	}
func (h *Handler) Introspect(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	// 验证JWT令牌
	claims, err := h.jwtService.ParseAccessToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"active": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":    true,
		"sub":       claims.UserID,
		"client_id": claims.ClientID,
		"scope":     claims.Scope,
		"exp":       claims.Exp,
		"iat":       claims.Iat,
		"iss":       claims.Iss,
		"aud":       claims.Aud,
	})
}

// base64URLEncode encodes bytes to base64url format (RFC 4648)
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
