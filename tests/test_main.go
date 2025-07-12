package tests

import (
	"log"
	"os"
	"testing"
)

// TestMain is the entry point for all tests in this package
// It handles global setup and teardown for the test suite
func TestMain(m *testing.M) {
	log.Println("🚀 Starting Flash OAuth2 E2E Test Suite")

	// Setup test environment
	setupTestEnvironment()

	// Run all tests
	code := m.Run()

	// Cleanup test environment
	cleanupTestEnvironment()

	log.Println("🏁 Flash OAuth2 E2E Test Suite Completed")

	// Exit with the same code as the tests
	os.Exit(code)
}

// setupTestEnvironment prepares the global test environment
func setupTestEnvironment() {
	log.Println("📋 Setting up test environment...")

	// Set default test environment variables if not already set
	if os.Getenv("TEST_DATABASE_URL") == "" {
		os.Setenv("TEST_DATABASE_URL", "postgres://postgres:1q2w3e4r@localhost:5432/oauth2_test?sslmode=disable")
	}

	if os.Getenv("TEST_REDIS_URL") == "" {
		os.Setenv("TEST_REDIS_URL", "redis://localhost:6379/15")
	}

	if os.Getenv("TEST_PORT") == "" {
		os.Setenv("TEST_PORT", "8081")
	}

	if os.Getenv("GIN_MODE") == "" {
		os.Setenv("GIN_MODE", "test")
	}

	// Set JWT secret for testing
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-jwt-secret-key-for-oauth2-server")
	}

	log.Printf("  Database: %s", os.Getenv("TEST_DATABASE_URL"))
	log.Printf("  Redis: %s", os.Getenv("TEST_REDIS_URL"))
	log.Printf("  Port: %s", os.Getenv("TEST_PORT"))
	log.Printf("  Mode: %s", os.Getenv("GIN_MODE"))

	log.Println("✅ Test environment setup completed")
}

// cleanupTestEnvironment cleans up after all tests
func cleanupTestEnvironment() {
	log.Println("🧹 Cleaning up test environment...")

	// Any global cleanup can be done here
	// Note: Individual test cleanup is handled in each test

	log.Println("✅ Test environment cleanup completed")
}

// TestConfig holds common test configuration
type TestConfig struct {
	DatabaseURL string
	RedisURL    string
	TestPort    string
	JWTSecret   string
}

// GetTestConfig returns the current test configuration
func GetTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL", "postgres://postgres:1q2w3e4r@localhost:5432/oauth2_test?sslmode=disable"),
		RedisURL:    getEnvOrDefault("TEST_REDIS_URL", "redis://localhost:6379/15"),
		TestPort:    getEnvOrDefault("TEST_PORT", "8081"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "test-jwt-secret-key-for-oauth2-server"),
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
