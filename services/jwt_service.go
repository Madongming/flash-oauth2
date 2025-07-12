// Package services provides JWT token generation and validation services.
// This package implements JWT-based OAuth2 tokens using RSA asymmetric encryption
// for secure token signing and verification.
package services

import (
	"crypto/rsa"
	"flash-oauth2/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService provides JWT token generation and validation using RSA asymmetric encryption.
// It supports creating OAuth2 access tokens and OpenID Connect ID tokens with
// proper claims and cryptographic signatures.
type JWTService struct {
	privateKey *rsa.PrivateKey // RSA private key for token signing
	publicKey  *rsa.PublicKey  // RSA public key for token verification
	issuer     string          // Token issuer identifier (typically server URL)
}

// NewJWTService creates a new JWTService instance with RSA key pair and issuer configuration.
//
// Parameters:
//   - privateKey: RSA private key for signing tokens
//   - publicKey: RSA public key for verifying tokens
//   - issuer: The issuer identifier for generated tokens (e.g., "https://auth.example.com")
//
// Returns:
//   - *JWTService: Configured JWT service instance
func NewJWTService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, issuer string) *JWTService {
	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
	}
}

// GenerateAccessToken creates a signed JWT access token for OAuth2 authentication.
// The token includes user ID, client ID, scope, and standard JWT claims with 1-hour expiration.
//
// Parameters:
//   - userID: The authenticated user's unique identifier
//   - clientID: The OAuth2 client that requested the token
//   - scope: The granted OAuth2 scopes (space-separated)
//
// Returns:
//   - string: The signed JWT access token
//   - error: An error if token generation or signing fails
//
// Example:
//
//	token, err := jwtService.GenerateAccessToken(123, "my-app", "openid profile")
func (s *JWTService) GenerateAccessToken(userID int, clientID, scope string) (string, error) {
	now := time.Now()
	claims := &models.AccessTokenClaims{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
		Exp:      now.Add(1 * time.Hour).Unix(),
		Iat:      now.Unix(),
		Iss:      s.issuer,
		Aud:      clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":        claims.UserID,
		"client_id":  claims.ClientID,
		"scope":      claims.Scope,
		"exp":        claims.Exp,
		"iat":        claims.Iat,
		"iss":        claims.Iss,
		"aud":        claims.Aud,
		"token_type": "access_token",
		"jti":        fmt.Sprintf("%d-%d", now.UnixNano(), userID), // Add unique JWT ID
	})

	return token.SignedString(s.privateKey)
}

// GenerateIDToken creates a signed JWT ID token for OpenID Connect authentication.
// ID tokens contain user identity information and are used by clients to verify user authentication.
//
// Parameters:
//   - user: The authenticated user object containing identity information
//   - clientID: The OAuth2 client that requested the token
//
// Returns:
//   - string: The signed JWT ID token
//   - error: An error if token generation or signing fails
//
// Example:
//
//	idToken, err := jwtService.GenerateIDToken(user, "my-app")
func (s *JWTService) GenerateIDToken(user *models.User, clientID string) (string, error) {
	now := time.Now()
	claims := &models.IDTokenClaims{
		UserID:   user.ID,
		Phone:    user.Phone,
		Exp:      now.Add(1 * time.Hour).Unix(),
		Iat:      now.Unix(),
		Iss:      s.issuer,
		Aud:      clientID,
		AuthTime: now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":        claims.UserID,
		"phone":      claims.Phone,
		"exp":        claims.Exp,
		"iat":        claims.Iat,
		"iss":        claims.Iss,
		"aud":        claims.Aud,
		"auth_time":  claims.AuthTime,
		"token_type": "id_token",
	})

	return token.SignedString(s.privateKey)
}

// ValidateToken verifies the signature and validity of a JWT token using the RSA public key.
// This method ensures the token was signed by this service and hasn't been tampered with.
//
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - *jwt.Token: The parsed and validated token object
//   - error: An error if validation fails or signature is invalid
//
// Example:
//
//	token, err := jwtService.ValidateToken("eyJhbGciOiJSUzI1NiI...")
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用的是RSA签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.publicKey, nil
	})
}

// ParseAccessToken validates and parses an access token, extracting its claims.
// This method is used to authenticate API requests and extract user/client information.
//
// Parameters:
//   - tokenString: The JWT access token string to parse
//
// Returns:
//   - *models.AccessTokenClaims: The parsed token claims containing user and scope information
//   - error: An error if parsing fails or token is invalid
//
// Example:
//
//	claims, err := jwtService.ParseAccessToken("eyJhbGciOiJSUzI1NiI...")
func (s *JWTService) ParseAccessToken(tokenString string) (*models.AccessTokenClaims, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// 检查token类型
	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "access_token" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	accessTokenClaims := &models.AccessTokenClaims{
		UserID:   int(claims["sub"].(float64)),
		ClientID: claims["client_id"].(string),
		Scope:    claims["scope"].(string),
		Exp:      int64(claims["exp"].(float64)),
		Iat:      int64(claims["iat"].(float64)),
		Iss:      claims["iss"].(string),
		Aud:      claims["aud"].(string),
	}

	return accessTokenClaims, nil
}

// GetPublicKey returns the RSA public key used for token verification
func (s *JWTService) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}
