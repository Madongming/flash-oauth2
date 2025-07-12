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

	// Get test configuration (will use environment variables if available)
	config := GetTestConfig()

	// Set environment variables for other components that might need them
	os.Setenv("TEST_DATABASE_URL", config.DatabaseURL)
	os.Setenv("TEST_REDIS_URL", config.RedisURL)
	os.Setenv("TEST_PORT", config.TestPort)
	os.Setenv("GIN_MODE", "test")
	os.Setenv("JWT_SECRET", config.JWTSecret)

	log.Printf("  Database: %s", config.DatabaseURL)
	log.Printf("  Redis: %s", config.RedisURL)
	log.Printf("  Port: %s", config.TestPort)
	log.Printf("  JWT Secret: %s", maskSecret(config.JWTSecret))
	log.Printf("  Test Timeout: %v", config.TestTimeout)
	log.Printf("  Cleanup Data: %v", config.CleanupData)

	log.Println("✅ Test environment setup completed")
}

// cleanupTestEnvironment cleans up after all tests
func cleanupTestEnvironment() {
	log.Println("🧹 Cleaning up test environment...")

	// Get configuration to check if cleanup is enabled
	config := GetTestConfig()
	if config.CleanupData {
		log.Println("  Cleanup enabled - test data will be cleaned")
		// Individual test cleanup is handled in each test
	} else {
		log.Println("  Cleanup disabled - test data will be preserved")
	}

	log.Println("✅ Test environment cleanup completed")
}

// Helper function to mask secrets in logs
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}
