package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfigFeatures 测试新的配置功能
func TestNewConfigFeatures(t *testing.T) {
	t.Log("🔧 Testing new configuration features")

	config := GetTestConfig()

	t.Run("Callback URL Generation", func(t *testing.T) {
		// 测试默认回调URL
		defaultURL := config.GetDefaultCallbackURL()
		assert.Equal(t, "http://localhost:3000/callback", defaultURL)

		// 测试自定义回调URL
		customURL := config.GetCallbackURL(3005, "/auth/callback")
		assert.Equal(t, "http://localhost:3005/auth/callback", customURL)

		// 测试所有回调URL
		allURLs := config.GetAllCallbackURLs()
		assert.Len(t, allURLs, 3)
		assert.Contains(t, allURLs, "http://localhost:3000/callback")
		assert.Contains(t, allURLs, "http://localhost:3001/callback")
		assert.Contains(t, allURLs, "http://localhost:3002/callback")
	})

	t.Run("Debug Mode Detection", func(t *testing.T) {
		// 测试非调试模式
		assert.False(t, config.IsDebugMode())

		// 设置调试模式
		os.Setenv("TEST_LOG_LEVEL", "debug")
		debugConfig := GetTestConfig()
		assert.True(t, debugConfig.IsDebugMode())

		// 恢复原设置
		os.Unsetenv("TEST_LOG_LEVEL")
	})

	t.Run("Environment Variable Override", func(t *testing.T) {
		// 测试环境变量覆盖
		os.Setenv("TEST_CALLBACK_PORTS", "4000,4001")
		os.Setenv("TEST_DATA_FACTORY", "random")

		envConfig := GetTestConfig()
		assert.Equal(t, []int{4000, 4001}, envConfig.CallbackPorts)
		assert.Equal(t, "random", envConfig.DataFactory)

		// 清理环境变量
		os.Unsetenv("TEST_CALLBACK_PORTS")
		os.Unsetenv("TEST_DATA_FACTORY")
	})

	t.Run("Redis DB Configuration", func(t *testing.T) {
		// 测试默认Redis DB
		defaultConfig := GetTestConfig()
		assert.Contains(t, defaultConfig.RedisURL, "/15", "Default Redis URL should use DB 15")

		// 测试自定义Redis DB
		os.Setenv("TEST_REDIS_DB", "10")
		customConfig := GetTestConfig()
		assert.Contains(t, customConfig.RedisURL, "/10", "Custom Redis URL should use DB 10")

		// 清理环境变量
		os.Unsetenv("TEST_REDIS_DB")
	})

	t.Log("✅ New configuration features working correctly")
}

// TestTestDataManager 测试数据管理器功能
func TestTestDataManager(t *testing.T) {
	t.Log("📋 Testing TestDataManager features")

	config := GetTestConfig()
	dataManager := NewTestDataManager(config)

	t.Run("Test Users", func(t *testing.T) {
		users := dataManager.GetTestUsers()
		assert.NotEmpty(t, users)

		standardUser := users[DefaultUserType]
		assert.NotNil(t, standardUser)
		assert.NotEmpty(t, standardUser.Phone)
		assert.NotEmpty(t, standardUser.Name)
		assert.NotEmpty(t, standardUser.VerifyCode)
	})

	t.Run("Test Clients", func(t *testing.T) {
		clients := dataManager.GetTestClients()
		assert.NotEmpty(t, clients)

		standardClient := clients[DefaultClientType]
		assert.NotNil(t, standardClient)
		assert.NotEmpty(t, standardClient.ID)
		assert.NotEmpty(t, standardClient.Secret)
		assert.NotEmpty(t, standardClient.RedirectURIs)

		// 验证回调URL使用了配置中的值
		assert.Contains(t, standardClient.RedirectURIs[0], config.TestServerHost)
	})

	t.Run("Test Scenarios", func(t *testing.T) {
		scenarios := dataManager.GetTestScenarios()
		assert.NotEmpty(t, scenarios)

		basicFlow := scenarios["basic_oauth2"]
		assert.NotNil(t, basicFlow)
		assert.Equal(t, "Basic OAuth2 Flow", basicFlow.Name)
		assert.NotEmpty(t, basicFlow.Users)
		assert.NotEmpty(t, basicFlow.Clients)
	})

	t.Run("Test Data Factory", func(t *testing.T) {
		// 测试预定义数据
		userData := dataManager.CreateTestData("user", DefaultUserType)
		assert.NotNil(t, userData)

		clientData := dataManager.CreateTestData("client", PublicClientType)
		assert.NotNil(t, clientData)
	})

	t.Run("Test Tokens and Codes", func(t *testing.T) {
		states := dataManager.GetTestStates()
		assert.NotEmpty(t, states)
		assert.True(t, len(states) >= 5)

		authCodes := dataManager.GetTestAuthCodes()
		assert.NotEmpty(t, authCodes)
		assert.True(t, len(authCodes) >= 4)

		tokens := dataManager.GetTestTokens()
		assert.NotEmpty(t, tokens)
		assert.Contains(t, tokens, "access_token")
		assert.Contains(t, tokens, "refresh_token")
		assert.Contains(t, tokens, "id_token")

		jwks := dataManager.GetTestJWKS()
		assert.NotNil(t, jwks)
		assert.Contains(t, jwks, "keys")
	})

	t.Log("✅ TestDataManager features working correctly")
}
