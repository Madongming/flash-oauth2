package tests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// Test data type constants for better separation
const (
	// Default test data types
	DefaultDeveloperType = "standard"
	DefaultAppType       = "web"
	DefaultClientType    = "standard"
	DefaultUserType      = "standard"

	// Alternative test data types
	EnterpriseDeveloperType = "enterprise"
	StartupDeveloperType    = "startup"
	MobileAppType           = "mobile"
	APIAppType              = "api"

	// Additional client types for different scenarios
	PublicClientType       = "public"
	ConfidentialClientType = "confidential"

	// Additional user types for different scenarios
	PremiumUserType  = "premium"
	InactiveUserType = "inactive"
	NewUserType      = "new"
)

// TestDataManager 管理测试数据的生成和提供
type TestDataManager struct {
	config *TestConfig
}

// NewTestDataManager 创建测试数据管理器
func NewTestDataManager(config *TestConfig) *TestDataManager {
	return &TestDataManager{
		config: config,
	}
}

// TestClientData 测试客户端数据结构
type TestClientData struct {
	ID           string
	Secret       string
	Name         string
	RedirectURIs []string
	GrantTypes   []string
	Scopes       []string
}

// TestUserData 测试用户数据结构
type TestUserData struct {
	ID         int
	Phone      string
	Email      string
	Name       string
	VerifyCode string
	CreatedAt  time.Time
	IsActive   bool
}

// TestDeveloperData 测试开发者数据结构
type TestDeveloperData struct {
	ID        string
	Name      string
	Email     string
	Company   string
	Status    string
	CreatedAt time.Time
}

// TestAppData 测试应用数据结构
type TestAppData struct {
	ID          string
	Name        string
	Description string
	Website     string
	DeveloperID string
	Status      string
	CreatedAt   time.Time
}

// TestKeyPairData 测试密钥对数据结构
type TestKeyPairData struct {
	ID         string
	AppID      string
	PrivateKey string
	PublicKey  string
	Status     string
	CreatedAt  time.Time
}

// Predefined test clients for different scenarios
func (tdm *TestDataManager) GetTestClients() map[string]*TestClientData {
	// 使用配置中的回调URL而不是硬编码
	callbackURLs := tdm.config.GetAllCallbackURLs()

	return map[string]*TestClientData{
		DefaultClientType: {
			ID:           "test-client-standard",
			Secret:       "test-secret-standard-123",
			Name:         "Standard Test Client",
			RedirectURIs: callbackURLs,
			GrantTypes:   []string{"authorization_code", "refresh_token"},
			Scopes:       []string{"openid", "profile", "email"},
		},
		PublicClientType: {
			ID:           "test-client-public",
			Secret:       "", // Public client has no secret
			Name:         "Public Test Client",
			RedirectURIs: []string{tdm.config.GetDefaultCallbackURL()},
			GrantTypes:   []string{"authorization_code"},
			Scopes:       []string{"openid", "profile"},
		},
		ConfidentialClientType: {
			ID:           "test-client-confidential",
			Secret:       "test-secret-confidential-456",
			Name:         "Confidential Test Client",
			RedirectURIs: []string{"https://app.example.com/callback"},
			GrantTypes:   []string{"authorization_code", "refresh_token", "client_credentials"},
			Scopes:       []string{"openid", "profile", "email", "admin"},
		},
		MobileAppType: {
			ID:           "test-client-mobile",
			Secret:       "test-secret-mobile-789",
			Name:         "Mobile Test Client",
			RedirectURIs: []string{"com.example.app://oauth/callback"},
			GrantTypes:   []string{"authorization_code", "refresh_token"},
			Scopes:       []string{"openid", "profile"},
		},
	}
}

// Predefined test users for different scenarios
func (tdm *TestDataManager) GetTestUsers() map[string]*TestUserData {
	return map[string]*TestUserData{
		DefaultUserType: {
			Phone:      "13800138000",
			Email:      "user1@example.com",
			Name:       "Standard Test User",
			VerifyCode: "123456",
			IsActive:   true,
		},
		PremiumUserType: {
			Phone:      "13800138001",
			Email:      "premium@example.com",
			Name:       "Premium Test User",
			VerifyCode: "654321",
			IsActive:   true,
		},
		InactiveUserType: {
			Phone:      "13800138002",
			Email:      "inactive@example.com",
			Name:       "Inactive Test User",
			VerifyCode: "111111",
			IsActive:   false,
		},
		NewUserType: {
			Phone:      "13800138003",
			Email:      "newuser@example.com",
			Name:       "New Test User",
			VerifyCode: "999999",
			IsActive:   true,
		},
	}
}

// Predefined test developers
func (tdm *TestDataManager) GetTestDevelopers() map[string]*TestDeveloperData {
	return map[string]*TestDeveloperData{
		DefaultDeveloperType: {
			ID:      "dev-001",
			Name:    "John Developer",
			Email:   "john@testcompany.com",
			Company: "Test Company Ltd",
			Status:  "active",
		},
		EnterpriseDeveloperType: {
			ID:      "dev-002",
			Name:    "Enterprise Developer",
			Email:   "enterprise@bigcorp.com",
			Company: "Big Corporation",
			Status:  "active",
		},
		StartupDeveloperType: {
			ID:      "dev-003",
			Name:    "Startup Developer",
			Email:   "dev@startup.io",
			Company: "Cool Startup Inc",
			Status:  "pending",
		},
	}
}

// Predefined test applications
func (tdm *TestDataManager) GetTestApps() map[string]*TestAppData {
	return map[string]*TestAppData{
		DefaultAppType: {
			ID:          "app-web-001",
			Name:        "Web Application",
			Description: "A test web application",
			Website:     "https://webapp.example.com",
			DeveloperID: "dev-001",
			Status:      "active",
		},
		MobileAppType: {
			ID:          "app-mobile-001",
			Name:        "Mobile Application",
			Description: "A test mobile application",
			Website:     "https://mobileapp.example.com",
			DeveloperID: "dev-001",
			Status:      "active",
		},
		APIAppType: {
			ID:          "app-api-001",
			Name:        "API Service",
			Description: "A test API service",
			Website:     "https://api.example.com",
			DeveloperID: "dev-002",
			Status:      "active",
		},
	}
}

// Generate random test data for edge cases and stress testing
func (tdm *TestDataManager) GenerateRandomClient() *TestClientData {
	return &TestClientData{
		ID:           generateRandomID("client"),
		Secret:       generateRandomSecret(),
		Name:         fmt.Sprintf("Random Test Client %s", generateRandomString(5)),
		RedirectURIs: []string{fmt.Sprintf("http://localhost:%d/callback", generateRandomPort())},
		GrantTypes:   []string{"authorization_code"},
		Scopes:       []string{"openid", "profile"},
	}
}

func (tdm *TestDataManager) GenerateRandomUser() *TestUserData {
	return &TestUserData{
		Phone:      generateRandomPhone(),
		Email:      fmt.Sprintf("random%s@example.com", generateRandomString(8)),
		Name:       fmt.Sprintf("Random User %s", generateRandomString(6)),
		VerifyCode: generateRandomVerifyCode(),
		IsActive:   true,
	}
}

func (tdm *TestDataManager) GenerateRandomDeveloper() *TestDeveloperData {
	return &TestDeveloperData{
		ID:      generateRandomID("dev"),
		Name:    fmt.Sprintf("Random Developer %s", generateRandomString(8)),
		Email:   fmt.Sprintf("dev%s@company.com", generateRandomString(6)),
		Company: fmt.Sprintf("Company %s", generateRandomString(10)),
		Status:  "active",
	}
}

// Test data validation scenarios
func (tdm *TestDataManager) GetInvalidPhoneNumbers() []string {
	return []string{
		"",                            // Empty
		"123",                         // Too short
		"1380013800013800138000",      // Too long
		"abc123456789",                // Contains letters
		"138-0013-8000",               // Contains dashes
		"+86 138 0013 8000",           // Contains spaces and plus
		"138001380001380013800013800", // Extremely long
	}
}

func (tdm *TestDataManager) GetInvalidVerifyCodes() []string {
	return []string{
		"",        // Empty
		"12345",   // Too short
		"1234567", // Too long
		"abcdef",  // Contains letters
		"12 34",   // Contains spaces
		"123-45",  // Contains dash
	}
}

func (tdm *TestDataManager) GetInvalidClientData() map[string]*TestClientData {
	defaultCallback := tdm.config.GetDefaultCallbackURL()
	return map[string]*TestClientData{
		"empty_id": {
			ID:           "",
			Secret:       "valid-secret",
			Name:         "Test Client",
			RedirectURIs: []string{defaultCallback},
		},
		"invalid_redirect": {
			ID:           "test-client",
			Secret:       "valid-secret",
			Name:         "Test Client",
			RedirectURIs: []string{"not-a-valid-url"},
		},
		"empty_name": {
			ID:           "test-client",
			Secret:       "valid-secret",
			Name:         "",
			RedirectURIs: []string{defaultCallback},
		},
	}
}

// Helper functions for generating random data
func generateRandomID(prefix string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, generateRandomString(8), generateRandomString(4))
}

func generateRandomSecret() string {
	return fmt.Sprintf("secret-%s-%s", generateRandomString(16), generateRandomString(8))
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func generateRandomPhone() string {
	// Generate valid Chinese mobile phone number format
	prefixes := []string{"138", "139", "150", "151", "152", "157", "158", "159", "182", "183", "184", "187", "188"}
	prefix := prefixes[mustGetRandomInt(len(prefixes))]
	suffix := fmt.Sprintf("%08d", mustGetRandomInt(100000000))
	return prefix + suffix
}

func generateRandomVerifyCode() string {
	return fmt.Sprintf("%06d", mustGetRandomInt(1000000))
}

func generateRandomPort() int {
	// Generate port between 30000-39999 for testing
	port := 30000 + mustGetRandomInt(10000)
	return port
}

func mustGetRandomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

// 测试场景数据生成器
type TestScenario struct {
	Name        string
	Description string
	Users       []*TestUserData
	Clients     []*TestClientData
	Setup       func(*TestDataManager) error
	Teardown    func(*TestDataManager) error
}

// GetTestScenarios 获取预定义的测试场景
func (tdm *TestDataManager) GetTestScenarios() map[string]*TestScenario {
	return map[string]*TestScenario{
		"basic_oauth2": {
			Name:        "Basic OAuth2 Flow",
			Description: "标准的OAuth2授权码流程测试",
			Users:       []*TestUserData{tdm.GetTestUsers()[DefaultUserType]},
			Clients:     []*TestClientData{tdm.GetTestClients()[DefaultClientType]},
		},
		"public_client": {
			Name:        "Public Client Flow",
			Description: "公开客户端（无secret）OAuth2流程测试",
			Users:       []*TestUserData{tdm.GetTestUsers()[DefaultUserType]},
			Clients:     []*TestClientData{tdm.GetTestClients()[PublicClientType]},
		},
		"mobile_app": {
			Name:        "Mobile App Flow",
			Description: "移动应用OAuth2流程测试",
			Users:       []*TestUserData{tdm.GetTestUsers()[PremiumUserType]},
			Clients:     []*TestClientData{tdm.GetTestClients()[MobileAppType]},
		},
		"multi_user": {
			Name:        "Multi User Scenario",
			Description: "多用户测试场景",
			Users: []*TestUserData{
				tdm.GetTestUsers()[DefaultUserType],
				tdm.GetTestUsers()[PremiumUserType],
				tdm.GetTestUsers()[InactiveUserType],
			},
			Clients: []*TestClientData{tdm.GetTestClients()[DefaultClientType]},
		},
		"edge_cases": {
			Name:        "Edge Cases",
			Description: "边界情况和错误处理测试",
			Users:       []*TestUserData{tdm.GetTestUsers()[NewUserType]},
			Clients:     []*TestClientData{tdm.GetTestClients()[ConfidentialClientType]},
		},
	}
}

// 常用的测试数据组合
func (tdm *TestDataManager) GetTestStates() []string {
	return []string{
		"test-state-basic-001",
		"test-state-mobile-002",
		"test-state-web-003",
		"test-state-edge-004",
		"test-state-security-005",
	}
}

func (tdm *TestDataManager) GetTestAuthCodes() []string {
	return []string{
		"test-auth-code-standard-001",
		"test-auth-code-public-002",
		"test-auth-code-mobile-003",
		"test-auth-code-refresh-004",
	}
}

func (tdm *TestDataManager) GetTestTokens() map[string]string {
	return map[string]string{
		"access_token":  "test-access-token-" + generateRandomString(8),
		"refresh_token": "test-refresh-token-" + generateRandomString(8),
		"id_token":      "test-id-token-" + generateRandomString(8),
	}
}

// 测试数据工厂方法 - 根据配置选择数据生成策略
func (tdm *TestDataManager) CreateTestData(dataType, scenario string) interface{} {
	switch tdm.config.DataFactory {
	case "random":
		return tdm.createRandomData(dataType)
	case "mixed":
		// 50% 概率使用随机数据
		if mustGetRandomInt(2) == 0 {
			return tdm.createRandomData(dataType)
		}
		return tdm.createPredefinedData(dataType, scenario)
	default: // "predefined"
		return tdm.createPredefinedData(dataType, scenario)
	}
}

func (tdm *TestDataManager) createRandomData(dataType string) interface{} {
	switch dataType {
	case "user":
		return tdm.GenerateRandomUser()
	case "client":
		return tdm.GenerateRandomClient()
	case "developer":
		return tdm.GenerateRandomDeveloper()
	default:
		return nil
	}
}

func (tdm *TestDataManager) createPredefinedData(dataType, scenario string) interface{} {
	switch dataType {
	case "user":
		users := tdm.GetTestUsers()
		if user, ok := users[scenario]; ok {
			return user
		}
		return users[DefaultUserType] // fallback
	case "client":
		clients := tdm.GetTestClients()
		if client, ok := clients[scenario]; ok {
			return client
		}
		return clients[DefaultClientType] // fallback
	case "developer":
		developers := tdm.GetTestDevelopers()
		if dev, ok := developers[scenario]; ok {
			return dev
		}
		return developers[DefaultDeveloperType] // fallback
	default:
		return nil
	}
}

// JWKS测试数据
func (tdm *TestDataManager) GetTestJWKS() map[string]any {
	keyID := "test-key-" + generateRandomString(4)
	return map[string]any{
		"keys": []any{
			map[string]any{
				"kty": "RSA",
				"use": "sig",
				"kid": keyID,
				"n":   "test-modulus-" + generateRandomString(8),
				"e":   "AQAB",
			},
		},
	}
}
