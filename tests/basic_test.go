package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBasicSetup verifies that the test environment is set up correctly
func TestBasicSetup(t *testing.T) {
	t.Log("🚀 Testing basic E2E test setup")

	// Test basic assertions
	assert.Equal(t, 1, 1, "Basic assertion should work")
	assert.True(t, true, "Boolean assertion should work")

	// Test time operations
	now := time.Now()
	assert.True(t, now.Before(now.Add(time.Second)), "Time comparison should work")

	t.Log("✅ Basic test setup is working correctly")
}

// TestEnvironmentVariables checks if environment variables are accessible
func TestEnvironmentVariables(t *testing.T) {
	t.Log("🔧 Testing environment variable access")

	config := GetTestConfig()

	assert.NotEmpty(t, config.DatabaseURL, "Database URL should not be empty")
	assert.NotEmpty(t, config.RedisURL, "Redis URL should not be empty")
	assert.NotEmpty(t, config.TestPort, "Test port should not be empty")

	t.Logf("Database URL: %s", config.DatabaseURL)
	t.Logf("Redis URL: %s", config.RedisURL)
	t.Logf("Test Port: %s", config.TestPort)

	t.Log("✅ Environment variables are accessible")
}
