// Package services provides OAuth2 and OpenID Connect service implementations.
// This package contains the core business logic for OAuth2 authorization flows,
// client management, and token operations.
package services

import (
	"crypto/rand"
	"database/sql"
	"flash-oauth2/models"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// OAuthService provides OAuth2 and OpenID Connect operations including
// client management, authorization code generation, and token lifecycle management.
// It implements the OAuth2 Authorization Code Flow with PKCE support.
type OAuthService struct {
	db *sql.DB // Database connection for persistent OAuth2 data storage
}

// NewOAuthService creates a new OAuthService instance with database connection.
//
// Parameters:
//   - db: Database connection for OAuth2 data persistence
//
// Returns:
//   - *OAuthService: Configured OAuth service instance
func NewOAuthService(db *sql.DB) *OAuthService {
	return &OAuthService{
		db: db,
	}
}

// GetClient retrieves an OAuth2 client configuration by client ID.
// This method validates client credentials and returns client metadata
// including redirect URIs, allowed grant types, and scopes.
//
// Parameters:
//   - clientID: The unique identifier of the OAuth2 client
//
// Returns:
//   - *models.OAuthClient: The client configuration if found
//   - error: An error if the client is not found or database operations fail
//
// Example:
//
//	client, err := oauthService.GetClient("my-app-client-id")
func (s *OAuthService) GetClient(clientID string) (*models.OAuthClient, error) {
	client := &models.OAuthClient{}
	err := s.db.QueryRow(`
		SELECT id, secret, name, redirect_uris, grant_types, response_types, scope, created_at 
		FROM oauth_clients WHERE id = $1
	`, clientID).Scan(
		&client.ID,
		&client.Secret,
		&client.Name,
		pq.Array(&client.RedirectURIs),
		pq.Array(&client.GrantTypes),
		pq.Array(&client.ResponseTypes),
		&client.Scope,
		&client.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// ValidateClient verifies OAuth2 client credentials (client ID and secret).
// This method is used during token exchange to authenticate the client application.
//
// Parameters:
//   - clientID: The client identifier to validate
//   - clientSecret: The client secret for authentication
//
// Returns:
//   - *models.OAuthClient: The validated client if credentials are correct
//   - error: An error if validation fails or client is not found
//
// Example:
//
//	client, err := oauthService.ValidateClient("my-app", "secret123")
func (s *OAuthService) ValidateClient(clientID, clientSecret string) (*models.OAuthClient, error) {
	client, err := s.GetClient(clientID)
	if err != nil {
		return nil, err
	}

	if client.Secret != clientSecret {
		return nil, fmt.Errorf("invalid client secret")
	}

	return client, nil
}

// CreateAuthCode generates a new authorization code for the OAuth2 Authorization Code Flow.
// The authorization code is used to exchange for access tokens and has a 10-minute expiration.
//
// Parameters:
//   - clientID: The OAuth2 client identifier requesting the authorization
//   - userID: The authenticated user's unique identifier
//   - redirectURI: The URI to redirect to after authorization
//   - scope: The requested OAuth2 scopes (space-separated)
//
// Returns:
//   - *models.AuthCode: The generated authorization code with metadata
//   - error: An error if database operations fail
//
// Example:
//
//	authCode, err := oauthService.CreateAuthCode("my-app", 123, "https://app.com/callback", "openid profile")
func (s *OAuthService) CreateAuthCode(clientID string, userID int, redirectURI, scope string) (*models.AuthCode, error) {
	code := generateRandomString(32)
	expiresAt := time.Now().Add(10 * time.Minute) // 授权码10分钟有效期

	authCode := &models.AuthCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   expiresAt,
	}

	_, err := s.db.Exec(`
		INSERT INTO auth_codes (code, client_id, user_id, redirect_uri, scope, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, authCode.Code, authCode.ClientID, authCode.UserID, authCode.RedirectURI, authCode.Scope, authCode.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return authCode, nil
}

// ExchangeAuthCode validates and exchanges an authorization code for user information.
// This method verifies the authorization code, client ID, and redirect URI match,
// and ensures the code hasn't expired. Used in the token exchange step of OAuth2 flow.
//
// Parameters:
//   - code: The authorization code to exchange
//   - clientID: The client ID that originally requested the code
//   - redirectURI: The redirect URI that must match the original request
//
// Returns:
//   - *models.AuthCode: The valid authorization code with user and scope information
//   - error: An error if the code is invalid, expired, or doesn't match parameters
//
// Example:
//
//	authCode, err := oauthService.ExchangeAuthCode("abc123", "my-app", "https://app.com/callback")
func (s *OAuthService) ExchangeAuthCode(code, clientID, redirectURI string) (*models.AuthCode, error) {
	authCode := &models.AuthCode{}
	err := s.db.QueryRow(`
		SELECT code, client_id, user_id, redirect_uri, scope, expires_at, created_at
		FROM auth_codes 
		WHERE code = $1 AND client_id = $2 AND redirect_uri = $3
	`, code, clientID, redirectURI).Scan(
		&authCode.Code,
		&authCode.ClientID,
		&authCode.UserID,
		&authCode.RedirectURI,
		&authCode.Scope,
		&authCode.ExpiresAt,
		&authCode.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 检查授权码是否过期
	if time.Now().After(authCode.ExpiresAt) {
		// 删除过期的授权码
		s.db.Exec("DELETE FROM auth_codes WHERE code = $1", code)
		return nil, fmt.Errorf("authorization code expired")
	}

	// 使用后删除授权码
	_, err = s.db.Exec("DELETE FROM auth_codes WHERE code = $1", code)
	if err != nil {
		return nil, err
	}

	return authCode, nil
}

// CreateAccessToken generates a new OAuth2 access token for authenticated API access.
// Access tokens have a 1-hour expiration and are used to authorize API requests.
//
// Parameters:
//   - clientID: The OAuth2 client identifier that requested the token
//   - userID: The authenticated user's unique identifier
//   - scope: The granted OAuth2 scopes for this token
//
// Returns:
//   - *models.AccessToken: The generated access token with metadata
//   - error: An error if database operations fail
//
// Example:
//
//	token, err := oauthService.CreateAccessToken("my-app", 123, "openid profile")
func (s *OAuthService) CreateAccessToken(clientID string, userID int, scope string) (*models.AccessToken, error) {
	token := generateRandomString(64)
	expiresAt := time.Now().Add(1 * time.Hour) // 访问令牌1小时有效期

	accessToken := &models.AccessToken{
		Token:     token,
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: expiresAt,
	}

	_, err := s.db.Exec(`
		INSERT INTO access_tokens (token, client_id, user_id, scope, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, accessToken.Token, accessToken.ClientID, accessToken.UserID, accessToken.Scope, accessToken.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// CreateRefreshToken generates a new OAuth2 refresh token for token renewal.
// Refresh tokens have a 30-day expiration and are used to obtain new access tokens
// without requiring user re-authentication.
//
// Parameters:
//   - clientID: The OAuth2 client identifier that requested the token
//   - userID: The authenticated user's unique identifier
//   - scope: The granted OAuth2 scopes for this token
//
// Returns:
//   - *models.RefreshToken: The generated refresh token with metadata
//   - error: An error if database operations fail
//
// Example:
//
//	refreshToken, err := oauthService.CreateRefreshToken("my-app", 123, "openid profile")
func (s *OAuthService) CreateRefreshToken(clientID string, userID int, scope string) (*models.RefreshToken, error) {
	token := generateRandomString(64)
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 刷新令牌30天有效期

	refreshToken := &models.RefreshToken{
		Token:     token,
		ClientID:  clientID,
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: expiresAt,
	}

	_, err := s.db.Exec(`
		INSERT INTO refresh_tokens (token, client_id, user_id, scope, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, refreshToken.Token, refreshToken.ClientID, refreshToken.UserID, refreshToken.Scope, refreshToken.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

// ValidateAccessToken verifies an OAuth2 access token and returns its metadata.
// This method checks if the token exists and hasn't expired, used for API authentication.
//
// Parameters:
//   - token: The access token string to validate
//
// Returns:
//   - *models.AccessToken: The token metadata if valid
//   - error: An error if the token is invalid, expired, or not found
//
// Example:
//
//	accessToken, err := oauthService.ValidateAccessToken("abc123token")
func (s *OAuthService) ValidateAccessToken(token string) (*models.AccessToken, error) {
	accessToken := &models.AccessToken{}
	err := s.db.QueryRow(`
		SELECT token, client_id, user_id, scope, expires_at, created_at
		FROM access_tokens 
		WHERE token = $1
	`, token).Scan(
		&accessToken.Token,
		&accessToken.ClientID,
		&accessToken.UserID,
		&accessToken.Scope,
		&accessToken.ExpiresAt,
		&accessToken.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 检查令牌是否过期
	if time.Now().After(accessToken.ExpiresAt) {
		// 删除过期的令牌
		s.db.Exec("DELETE FROM access_tokens WHERE token = $1", token)
		return nil, fmt.Errorf("access token expired")
	}

	return accessToken, nil
}

// RefreshAccessToken validates a refresh token and returns its metadata for token renewal.
// This method is used in the refresh token flow to obtain new access tokens.
//
// Parameters:
//   - refreshToken: The refresh token string to validate
//
// Returns:
//   - *models.RefreshToken: The refresh token metadata if valid
//   - error: An error if the token is invalid, expired, or not found
//
// Example:
//
//	refreshToken, err := oauthService.RefreshAccessToken("refresh123token")
func (s *OAuthService) RefreshAccessToken(refreshToken string) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	err := s.db.QueryRow(`
		SELECT token, client_id, user_id, scope, expires_at, created_at
		FROM refresh_tokens 
		WHERE token = $1
	`, refreshToken).Scan(
		&token.Token,
		&token.ClientID,
		&token.UserID,
		&token.Scope,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// 检查刷新令牌是否过期
	if time.Now().After(token.ExpiresAt) {
		// 删除过期的刷新令牌
		s.db.Exec("DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	return token, nil
}

// generateRandomString creates a cryptographically secure random string of specified length.
// This helper function is used to generate authorization codes, access tokens, and refresh tokens.
//
// Parameters:
//   - length: The desired length of the random string
//
// Returns:
//   - string: A random hexadecimal string of the specified length
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)

	// 使用base64编码并移除填充字符
	encoded := fmt.Sprintf("%x", bytes)
	if len(encoded) > length {
		return encoded[:length]
	}
	return encoded
}
