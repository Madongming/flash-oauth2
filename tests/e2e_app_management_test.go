package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAppManagementE2E tests the complete application management platform
func TestAppManagementE2E(t *testing.T) {
	// Check if we can set up a test server (requires database)
	ts := TrySetupTestServer(t)
	if ts == nil {
		t.Skip("Cannot setup test server (likely database not available)")
		return
	}
	defer ts.TeardownTestServer(t)

	t.Run("Complete App Management Flow", func(t *testing.T) {
		// 1. Register a developer
		developer := ts.RegisterTestDeveloper(t)

		// 2. Register an external application
		app := ts.RegisterTestExternalApp(t, developer.ID)

		// 3. Generate key pairs for the application
		keyPair := ts.GenerateTestKeyPair(t, app.ID)

		// 4. Get application details with keys
		ts.GetAppWithKeys(t, app.ID)

		// 5. Revoke a key
		ts.RevokeTestKey(t, keyPair.ID)

		// 6. Test dashboard functionality
		ts.TestManagementDashboard(t)

		// 7. Test app details page
		ts.TestAppDetailsPage(t, app.ID)

		// 8. Test API endpoints
		ts.TestAppManagementAPI(t, developer.ID, app.ID)
	})
}

// TestDeveloperRegistration tests developer registration functionality
func TestDeveloperRegistration(t *testing.T) {
	ts := TrySetupTestServer(t)
	if ts == nil {
		t.Skip("Cannot setup test server (likely database not available)")
		return
	}
	defer ts.TeardownTestServer(t)

	t.Run("Developer Registration", func(t *testing.T) {
		// Test valid developer registration
		t.Run("Valid Registration", func(t *testing.T) {
			developer := ts.RegisterTestDeveloper(t)
			assert.NotEmpty(t, developer.ID)
			assert.Equal(t, "Test Company", developer.Name)
			assert.Equal(t, "test@company.com", developer.Email)
			assert.Equal(t, "active", developer.Status)
		})
	})
}

// TestKeyManagement tests key generation and management
func TestKeyManagement(t *testing.T) {
	ts := TrySetupTestServer(t)
	if ts == nil {
		t.Skip("Cannot setup test server (likely database not available)")
		return
	}
	defer ts.TeardownTestServer(t)

	t.Run("Key Management", func(t *testing.T) {
		// Setup test data
		developer := ts.RegisterTestDeveloper(t)
		app := ts.RegisterTestExternalApp(t, developer.ID)

		// Test key generation
		t.Run("Generate Key Pair", func(t *testing.T) {
			keyPair := ts.GenerateTestKeyPair(t, app.ID)
			assert.NotEmpty(t, keyPair.ID)
			assert.NotEmpty(t, keyPair.PrivateKey)
			assert.NotEmpty(t, keyPair.PublicKey)
			assert.Equal(t, app.ID, keyPair.AppID)
			assert.Equal(t, "active", keyPair.Status)
			assert.True(t, strings.Contains(keyPair.PrivateKey, "BEGIN") && strings.Contains(keyPair.PrivateKey, "PRIVATE KEY"))
			assert.True(t, strings.Contains(keyPair.PublicKey, "BEGIN") && strings.Contains(keyPair.PublicKey, "PUBLIC KEY"))
		})
	})
}
