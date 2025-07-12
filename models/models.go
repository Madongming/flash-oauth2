// Package models defines the data structures used throughout the OAuth2 server.
// It includes database models for users, OAuth2 clients, authorization codes,
// access tokens, refresh tokens, and JWT claims structures.
package models

import (
	"time"
)

// User represents a user in the system.
// Users are identified by their phone number and can authenticate
// using phone number + verification code.
type User struct {
	ID        int       `json:"id" db:"id"`                 // Unique user identifier
	Phone     string    `json:"phone" db:"phone"`           // User's phone number (unique)
	Role      string    `json:"role" db:"role"`             // User role (user/admin)
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Account creation time
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update time
}

// OAuthClient represents an OAuth2 client application.
// Clients must be registered before they can request authorization.
type OAuthClient struct {
	ID            string    `json:"id" db:"id"`                         // Client identifier
	Secret        string    `json:"secret" db:"secret"`                 // Client secret
	Name          string    `json:"name" db:"name"`                     // Human-readable client name
	RedirectURIs  []string  `json:"redirect_uris" db:"redirect_uris"`   // Allowed redirect URIs
	GrantTypes    []string  `json:"grant_types" db:"grant_types"`       // Supported grant types
	ResponseTypes []string  `json:"response_types" db:"response_types"` // Supported response types
	Scope         string    `json:"scope" db:"scope"`                   // Default scopes
	CreatedAt     time.Time `json:"created_at" db:"created_at"`         // Client registration time
}

// AuthCode represents an OAuth2 authorization code.
// Authorization codes are short-lived tokens that can be exchanged for access tokens.
type AuthCode struct {
	Code        string    `json:"code" db:"code"`                 // The authorization code
	ClientID    string    `json:"client_id" db:"client_id"`       // Client that requested the code
	UserID      int       `json:"user_id" db:"user_id"`           // User who authorized the client
	RedirectURI string    `json:"redirect_uri" db:"redirect_uri"` // URI to redirect after authorization
	Scope       string    `json:"scope" db:"scope"`               // Requested scopes
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`     // Code expiration time
	CreatedAt   time.Time `json:"created_at" db:"created_at"`     // Code creation time
}

// AccessToken represents an OAuth2 access token stored in the database.
// Note: In this implementation, JWT tokens are self-contained and don't need database storage,
// but this table can be used for audit purposes and token tracking.
type AccessToken struct {
	Token     string    `json:"token" db:"token"`           // The access token
	ClientID  string    `json:"client_id" db:"client_id"`   // Client that owns the token
	UserID    int       `json:"user_id" db:"user_id"`       // User the token represents
	Scope     string    `json:"scope" db:"scope"`           // Token scopes
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"` // Token expiration time
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Token creation time
}

// RefreshToken represents an OAuth2 refresh token.
// Refresh tokens are long-lived tokens used to obtain new access tokens.
type RefreshToken struct {
	Token     string    `json:"token" db:"token"`           // The refresh token
	ClientID  string    `json:"client_id" db:"client_id"`   // Client that owns the token
	UserID    int       `json:"user_id" db:"user_id"`       // User the token represents
	Scope     string    `json:"scope" db:"scope"`           // Token scopes
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"` // Token expiration time
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Token creation time
}

// AccessTokenClaims represents the claims contained in a JWT access token.
// These claims follow OAuth2 and JWT standards.
type AccessTokenClaims struct {
	UserID   int    `json:"sub"`       // Subject (user ID)
	ClientID string `json:"client_id"` // Client identifier
	Scope    string `json:"scope"`     // Token scopes
	Exp      int64  `json:"exp"`       // Expiration time (Unix timestamp)
	Iat      int64  `json:"iat"`       // Issued at time (Unix timestamp)
	Iss      string `json:"iss"`       // Issuer
	Aud      string `json:"aud"`       // Audience (client ID)
}

// IDTokenClaims represents the claims contained in an OpenID Connect ID token.
// These claims follow OpenID Connect specifications.
type IDTokenClaims struct {
	UserID   int    `json:"sub"`       // Subject (user ID)
	Phone    string `json:"phone"`     // User's phone number
	Exp      int64  `json:"exp"`       // Expiration time (Unix timestamp)
	Iat      int64  `json:"iat"`       // Issued at time (Unix timestamp)
	Iss      string `json:"iss"`       // Issuer
	Aud      string `json:"aud"`       // Audience (client ID)
	AuthTime int64  `json:"auth_time"` // Authentication time (Unix timestamp)
}

// ExternalApp represents an external application registered on the platform
type ExternalApp struct {
	ID          string     `json:"id" db:"id"`                     // Unique application identifier
	Name        string     `json:"name" db:"name"`                 // Application name
	Description string     `json:"description" db:"description"`   // Application description
	DeveloperID string     `json:"developer_id" db:"developer_id"` // Developer/company identifier
	Status      string     `json:"status" db:"status"`             // active, suspended, revoked
	CallbackURL string     `json:"callback_url" db:"callback_url"` // OAuth callback URL
	Scopes      string     `json:"scopes" db:"scopes"`             // Allowed scopes (space-separated)
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`     // App registration time
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`     // Last update time
	RevokedAt   *time.Time `json:"revoked_at" db:"revoked_at"`     // Revocation time (if revoked)
}

// AppKeyPair represents an RSA key pair issued to an external application
type AppKeyPair struct {
	ID         string     `json:"id" db:"id"`                     // Unique key pair identifier
	AppID      string     `json:"app_id" db:"app_id"`             // Associated application ID
	KeyID      string     `json:"key_id" db:"key_id"`             // Key identifier (kid)
	PrivateKey string     `json:"private_key" db:"private_key"`   // RSA private key (PEM format)
	PublicKey  string     `json:"public_key" db:"public_key"`     // RSA public key (PEM format)
	Algorithm  string     `json:"algorithm" db:"algorithm"`       // Signing algorithm (RS256, RS384, RS512)
	Status     string     `json:"status" db:"status"`             // active, expired, revoked
	ExpiresAt  *time.Time `json:"expires_at" db:"expires_at"`     // Key expiration time (null = no expiry)
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`     // Key generation time
	RevokedAt  *time.Time `json:"revoked_at" db:"revoked_at"`     // Key revocation time
	LastUsedAt *time.Time `json:"last_used_at" db:"last_used_at"` // Last time key was used
}

// Developer represents a developer/company that can register applications
type Developer struct {
	ID        string    `json:"id" db:"id"`                 // Unique developer identifier
	Name      string    `json:"name" db:"name"`             // Developer/company name
	Email     string    `json:"email" db:"email"`           // Contact email
	Phone     string    `json:"phone" db:"phone"`           // Contact phone
	Status    string    `json:"status" db:"status"`         // active, suspended
	APIQuota  int       `json:"api_quota" db:"api_quota"`   // Daily API call quota
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Registration time
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update time
}
