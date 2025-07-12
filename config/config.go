// Package config provides configuration management for the OAuth2 server.
// It handles loading configuration from environment variables and generates
// RSA key pairs for JWT token signing.
package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

// SMSConfig holds configuration for SMS service (Alibaba Cloud)
type SMSConfig struct {
	AccessKeyId     string // Alibaba Cloud Access Key ID
	AccessKeySecret string // Alibaba Cloud Access Key Secret
	SignName        string // SMS signature name
	TemplateCode    string // SMS template code
	Enabled         bool   // Whether SMS sending is enabled (false for testing)
}

// Config holds all configuration values for the OAuth2 server.
// It includes server settings, database connections, and cryptographic keys.
type Config struct {
	Port          string          // HTTP server port
	DatabaseURL   string          // PostgreSQL connection string
	RedisURL      string          // Redis connection string
	JWTPrivateKey *rsa.PrivateKey // RSA private key for JWT signing
	JWTPublicKey  *rsa.PublicKey  // RSA public key for JWT verification
	SMS           *SMSConfig      // SMS configuration
}

// Load creates and returns a new Config instance with values loaded from
// environment variables. If environment variables are not set, it uses
// sensible defaults. It also generates a new RSA key pair for JWT signing.
//
// Environment variables:
//   - PORT: Server port (default: "8080")
//   - DATABASE_URL: PostgreSQL connection string
//   - REDIS_URL: Redis connection string
//
// The function will terminate the program if RSA key generation fails.
func Load() *Config {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate RSA key:", err)
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/2"),
		JWTPrivateKey: privateKey,
		JWTPublicKey:  &privateKey.PublicKey,
		SMS: &SMSConfig{
			AccessKeyId:     getEnv("SMS_ACCESS_KEY_ID", ""),
			AccessKeySecret: getEnv("SMS_ACCESS_KEY_SECRET", ""),
			SignName:        getEnv("SMS_SIGN_NAME", ""),
			TemplateCode:    getEnv("SMS_TEMPLATE_CODE", ""),
			Enabled:         getEnv("SMS_ENABLED", "false") == "true",
		},
	}
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set or empty, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetJWTPublicKeyPEM returns the public key in PEM format as a string.
// This is used for JWT token verification by clients and for the JWKS endpoint.
// The function will terminate the program if key marshaling fails.
func (c *Config) GetJWTPublicKeyPEM() string {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(c.JWTPublicKey)
	if err != nil {
		log.Fatal("Failed to marshal public key:", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM)
}
