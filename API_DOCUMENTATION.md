# Flash OAuth2 Server API Documentation

This document provides comprehensive API documentation for the Flash OAuth2 + OpenID Connect authorization server.

## Overview

Flash OAuth2 is a complete OAuth2 Authorization Server with OpenID Connect support, featuring:

- **OAuth2 Authorization Code Flow** with PKCE support
- **OpenID Connect (OIDC)** implementation
- **JWT tokens** with RSA asymmetric encryption
- **Phone-based authentication** with verification codes
- **PostgreSQL** persistence layer
- **Redis** for sessions and temporary data
- **Automatic user registration** (idempotent operations)

## Architecture

The server is organized into the following packages:

### Core Packages

1. **`main`** - Application entry point and server initialization
2. **`config`** - Configuration management and RSA key generation
3. **`models`** - Data structures and database models
4. **`database`** - PostgreSQL connection and schema management
5. **`services`** - Business logic services
6. **`handlers`** - HTTP request handlers and API endpoints
7. **`middleware`** - HTTP middleware (CORS, etc.)
8. **`routes`** - HTTP routing configuration

## API Endpoints

### OAuth2 Core Endpoints

#### Authorization Endpoint

```
GET /oauth/authorize
```

**Purpose**: Initiates the OAuth2 authorization flow
**Parameters**:

- `client_id` (required): OAuth2 client identifier
- `redirect_uri` (required): Callback URL after authorization
- `response_type` (required): Must be "code"
- `scope` (optional): Requested scopes (space-separated)
- `state` (optional): Client state parameter

**Response**: Redirects to login page or client redirect_uri with authorization code

#### Token Exchange Endpoint

```
POST /oauth/token
```

**Purpose**: Exchanges authorization code for access tokens
**Content-Type**: `application/x-www-form-urlencoded`
**Parameters**:

- `grant_type` (required): "authorization_code" or "refresh_token"
- `code` (required for auth code): Authorization code
- `client_id` (required): OAuth2 client identifier
- `client_secret` (required): OAuth2 client secret
- `redirect_uri` (required for auth code): Must match original request
- `refresh_token` (required for refresh): Refresh token

**Response**:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiI...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_token_string",
  "id_token": "eyJhbGciOiJSUzI1NiI...",
  "scope": "openid profile"
}
```

### OpenID Connect Endpoints

#### UserInfo Endpoint

```
GET /oauth/userinfo
```

**Purpose**: Returns user profile information
**Authorization**: `Bearer {access_token}`
**Response**:

```json
{
  "sub": "123",
  "phone": "13800138000"
}
```

#### JSON Web Key Set (JWKS)

```
GET /.well-known/jwks.json
```

**Purpose**: Returns public keys for token verification
**Response**:

```json
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "oauth2-rsa-key",
      "n": "...",
      "e": "AQAB"
    }
  ]
}
```

### Authentication Endpoints

#### Login Page

```
GET /login
```

**Purpose**: Displays phone authentication form
**Parameters**:

- `client_id`: OAuth2 client identifier
- `redirect_uri`: Callback URL
- `scope`: Requested scopes
- `state`: Client state

#### Send Verification Code

```
POST /api/send-code
```

**Purpose**: Sends SMS verification code to phone
**Content-Type**: `application/json`
**Body**:

```json
{
  "phone": "13800138000"
}
```

#### Verify Code and Login

```
POST /api/login
```

**Purpose**: Verifies code and completes authentication
**Content-Type**: `application/json`
**Body**:

```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

### Management Endpoints

#### Token Introspection

```
POST /oauth/introspect
```

**Purpose**: Validates and returns token metadata
**Content-Type**: `application/x-www-form-urlencoded`
**Parameters**:

- `token` (required): Access token to introspect
- `client_id` (required): Client identifier
- `client_secret` (required): Client secret

**Response**:

```json
{
  "active": true,
  "client_id": "my-app",
  "user_id": "123",
  "scope": "openid profile",
  "exp": 1640995200
}
```

## Services Documentation

### UserService

**Purpose**: Manages user authentication and verification codes

#### Methods

- **`SendVerificationCode(phone string) error`**

  - Generates and sends 6-digit verification code
  - Stores code in Redis with 5-minute expiration
  - In production, integrates with SMS service

- **`VerifyCode(phone, code string) (*models.User, error)`**

  - Validates verification code against Redis storage
  - Returns authenticated user or creates new user
  - Implements idempotent user registration

- **`GetUserByID(userID int) (*models.User, error)`**
  - Retrieves user by unique identifier
  - Used for token validation and user info endpoints

### OAuthService

**Purpose**: Implements OAuth2 authorization flows and client management

#### Methods

- **`GetClient(clientID string) (*models.OAuthClient, error)`**

  - Retrieves OAuth2 client configuration
  - Returns client metadata including redirect URIs and scopes

- **`ValidateClient(clientID, clientSecret string) (*models.OAuthClient, error)`**

  - Validates client credentials for token exchange
  - Used in confidential client authentication

- **`CreateAuthCode(clientID string, userID int, redirectURI, scope string) (*models.AuthCode, error)`**

  - Generates authorization code for OAuth2 flow
  - 10-minute expiration, single-use

- **`ExchangeAuthCode(code, clientID, redirectURI string) (*models.AuthCode, error)`**

  - Validates and exchanges authorization code
  - Ensures code hasn't expired and parameters match

- **`CreateAccessToken(clientID string, userID int, scope string) (*models.AccessToken, error)`**

  - Generates OAuth2 access token
  - 1-hour expiration for API authentication

- **`CreateRefreshToken(clientID string, userID int, scope string) (*models.RefreshToken, error)`**

  - Generates refresh token for token renewal
  - 30-day expiration

- **`ValidateAccessToken(token string) (*models.AccessToken, error)`**

  - Validates access token for API requests
  - Checks expiration and removes expired tokens

- **`RefreshAccessToken(refreshToken string) (*models.RefreshToken, error)`**
  - Validates refresh token for token renewal flow
  - Returns metadata for generating new access tokens

### JWTService

**Purpose**: Provides JWT token generation and validation with RSA encryption

#### Methods

- **`GenerateAccessToken(userID int, clientID, scope string) (string, error)`**

  - Creates signed JWT access token
  - Includes standard OAuth2 claims with 1-hour expiration

- **`GenerateIDToken(user *models.User, clientID string) (string, error)`**

  - Creates signed JWT ID token for OpenID Connect
  - Contains user identity information

- **`ValidateToken(tokenString string) (*jwt.Token, error)`**

  - Verifies JWT signature using RSA public key
  - Ensures token integrity and authenticity

- **`ParseAccessToken(tokenString string) (*models.AccessTokenClaims, error)`**
  - Parses and validates access token claims
  - Used for API request authentication

## Data Models

### User

```go
type User struct {
    ID        int       `json:"id"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### OAuthClient

```go
type OAuthClient struct {
    ID            string    `json:"id"`
    Secret        string    `json:"secret"`
    Name          string    `json:"name"`
    RedirectURIs  []string  `json:"redirect_uris"`
    GrantTypes    []string  `json:"grant_types"`
    ResponseTypes []string  `json:"response_types"`
    Scope         string    `json:"scope"`
    CreatedAt     time.Time `json:"created_at"`
}
```

### AuthCode

```go
type AuthCode struct {
    Code        string    `json:"code"`
    ClientID    string    `json:"client_id"`
    UserID      int       `json:"user_id"`
    RedirectURI string    `json:"redirect_uri"`
    Scope       string    `json:"scope"`
    ExpiresAt   time.Time `json:"expires_at"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### AccessToken

```go
type AccessToken struct {
    Token     string    `json:"token"`
    ClientID  string    `json:"client_id"`
    UserID    int       `json:"user_id"`
    Scope     string    `json:"scope"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}
```

### JWT Claims

#### AccessTokenClaims

```go
type AccessTokenClaims struct {
    UserID   int    `json:"sub"`
    ClientID string `json:"client_id"`
    Scope    string `json:"scope"`
    Exp      int64  `json:"exp"`
    Iat      int64  `json:"iat"`
    Iss      string `json:"iss"`
    Aud      string `json:"aud"`
}
```

#### IDTokenClaims

```go
type IDTokenClaims struct {
    UserID   int    `json:"sub"`
    Phone    string `json:"phone"`
    Exp      int64  `json:"exp"`
    Iat      int64  `json:"iat"`
    Iss      string `json:"iss"`
    Aud      string `json:"aud"`
    AuthTime int64  `json:"auth_time"`
}
```

## Configuration

### Environment Variables

- **`PORT`**: Server port (default: 8080)
- **`DATABASE_URL`**: PostgreSQL connection string
- **`REDIS_URL`**: Redis connection string

### Default Configuration

```go
type Config struct {
    Port        string
    DatabaseURL string
    RedisURL    string
    JWTIssuer   string
    PrivateKey  *rsa.PrivateKey
    PublicKey   *rsa.PublicKey
}
```

## Security Features

### RSA Asymmetric Encryption

- 2048-bit RSA key pairs for JWT signing
- Private key for token generation
- Public key for token verification
- Keys auto-generated on startup

### Token Security

- Cryptographically secure random token generation
- Short-lived access tokens (1 hour)
- Longer-lived refresh tokens (30 days)
- Single-use authorization codes (10 minutes)

### Phone Authentication

- 6-digit numeric verification codes
- 5-minute code expiration
- Redis-backed temporary storage
- Rate limiting ready (implementation dependent)

## Database Schema

### Tables

- **`users`**: User accounts with phone numbers
- **`oauth_clients`**: Registered OAuth2 applications
- **`auth_codes`**: Temporary authorization codes
- **`access_tokens`**: Active access tokens
- **`refresh_tokens`**: Long-lived refresh tokens

## Deployment

### Docker Support

- Dockerfile for containerized deployment
- docker-compose.yml for development environment
- PostgreSQL and Redis containers included

### Running the Server

```bash
# Development
go run main.go

# Production build
go build -o oauth2-server
./oauth2-server
```

## Standards Compliance

### OAuth2 (RFC 6749)

- Authorization Code Flow
- Client authentication
- Token introspection
- Refresh token flow

### OpenID Connect (OIDC)

- ID token generation
- UserInfo endpoint
- JWKS endpoint
- Standard claims

### JWT (RFC 7519)

- RSA-256 signatures
- Standard claims
- Proper expiration handling

## Generated with Go Doc

This documentation was generated using Go's built-in documentation system:

```bash
go doc -all                 # Complete package documentation
go doc ./services          # Services package
go doc ./models            # Models package
go doc ./handlers          # Handlers package
```

For detailed method signatures and examples, use:

```bash
go doc services.UserService.SendVerificationCode
go doc models.User
go doc handlers.Handler.Authorize
```
