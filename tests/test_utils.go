package tests

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// TestConfig 测试配置结构
// 这是一个公共的配置结构，供所有测试文件使用
type TestConfig struct {
	DatabaseURL string
	RedisURL    string
	TestPort    string
	JWTSecret   string

	// Database settings
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis settings
	RedisHost     string
	RedisPort     int
	RedisDB       int
	RedisPassword string

	// Test settings
	TestTimeout time.Duration
	CleanupData bool

	// Test endpoints and URLs
	TestServerHost string
	CallbackPorts  []int

	// Test data settings
	DataFactory     string // "predefined", "random", "mixed"
	UseRealServices bool   // 是否使用真实服务而不是mock
	LogLevel        string // "debug", "info", "warn", "error"
}

// GetTestConfig 获取测试配置
// 优先使用环境变量，如果没有则使用默认值
func GetTestConfig() *TestConfig {
	config := &TestConfig{
		// Default values
		DBHost:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		DBPort:     getEnvAsIntOrDefault("TEST_DB_PORT", 5432),
		DBUser:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("TEST_DB_PASSWORD", "1q2w3e4r"),
		DBName:     getEnvOrDefault("TEST_DB_NAME", "oauth2_test"),

		RedisHost:     getEnvOrDefault("TEST_REDIS_HOST", "localhost"),
		RedisPort:     getEnvAsIntOrDefault("TEST_REDIS_PORT", 6379),
		RedisDB:       getEnvAsIntOrDefault("TEST_REDIS_DB", 15),
		RedisPassword: getEnvOrDefault("TEST_REDIS_PASSWORD", ""),

		TestPort:    getEnvOrDefault("TEST_PORT", "8081"),
		JWTSecret:   getEnvOrDefault("TEST_JWT_SECRET", "test-jwt-secret-key-for-oauth2-server"),
		TestTimeout: getEnvAsDurationOrDefault("TEST_TIMEOUT", 30*time.Second),
		CleanupData: getEnvAsBoolOrDefault("TEST_CLEANUP_DATA", true),

		// New configurable settings
		TestServerHost:  getEnvOrDefault("TEST_SERVER_HOST", "localhost"),
		CallbackPorts:   getEnvAsIntSliceOrDefault("TEST_CALLBACK_PORTS", []int{3000, 3001, 3002}),
		DataFactory:     getEnvOrDefault("TEST_DATA_FACTORY", "predefined"),
		UseRealServices: getEnvAsBoolOrDefault("TEST_USE_REAL_SERVICES", false),
		LogLevel:        getEnvOrDefault("TEST_LOG_LEVEL", "info"),
	}

	// Build connection strings
	config.DatabaseURL = getEnvOrDefault("TEST_DATABASE_URL",
		buildPostgresURL(config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName))

	config.RedisURL = getEnvOrDefault("TEST_REDIS_URL",
		buildRedisURL(config.RedisHost, config.RedisPort, config.RedisPassword, config.RedisDB))

	return config
}

// Helper functions for environment variable handling
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsIntSliceOrDefault(key string, defaultValue []int) []int {
	if value := os.Getenv(key); value != "" {
		// Parse comma-separated integers: "3000,3001,3002"
		var result []int
		parts := strings.Split(value, ",")
		for _, part := range parts {
			if intValue, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				result = append(result, intValue)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func buildPostgresURL(host string, port int, user, password, dbname string) string {
	return "postgres://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname + "?sslmode=disable"
}

func buildRedisURL(host string, port int, password string, db int) string {
	if password != "" {
		return "redis://:" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + strconv.Itoa(db)
	}
	return "redis://" + host + ":" + strconv.Itoa(port) + "/" + strconv.Itoa(db)
}

// 便利方法
func (c *TestConfig) GetCallbackURL(port int, path string) string {
	if path == "" {
		path = "/callback"
	}
	return "http://" + c.TestServerHost + ":" + strconv.Itoa(port) + path
}

func (c *TestConfig) GetDefaultCallbackURL() string {
	return c.GetCallbackURL(c.CallbackPorts[0], "/callback")
}

func (c *TestConfig) GetAllCallbackURLs() []string {
	var urls []string
	for _, port := range c.CallbackPorts {
		urls = append(urls, c.GetCallbackURL(port, "/callback"))
	}
	return urls
}

func (c *TestConfig) IsDebugMode() bool {
	return c.LogLevel == "debug"
}
