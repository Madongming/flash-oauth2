package tests

import (
	"context"
	"database/sql"
	"net"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// TestEnvironmentSetup verifies the test environment connectivity
// It does NOT create database structure - that's the application's job
func TestEnvironmentSetup(t *testing.T) {
	t.Log("� Checking test environment connectivity (read-only)")

	// Test 1: Check PostgreSQL connectivity
	t.Run("PostgreSQL Connection", func(t *testing.T) {
		config := GetTestConfig()

		// Check if PostgreSQL port is accessible
		conn, err := net.DialTimeout("tcp", "localhost:5432", 5*time.Second)
		if err != nil {
			t.Skip("PostgreSQL not running on localhost:5432, skipping database tests")
			return
		}
		conn.Close()

		// Test database connection
		db, err := sql.Open("postgres", config.DatabaseURL)
		if err != nil {
			t.Skipf("Cannot connect to test database: %v", err)
			return
		}
		defer db.Close()

		// Test database ping
		err = db.Ping()
		if err != nil {
			t.Skipf("Cannot ping test database: %v (database should be created by application, not tests)", err)
			return
		}

		// Simple connectivity test
		var result int
		err = db.QueryRow("SELECT 1").Scan(&result)
		require.NoError(t, err, "Should execute basic query")
		require.Equal(t, 1, result, "Query should return 1")

		t.Log("✅ PostgreSQL connection verified")
	})
	// Test 2: Check Redis connectivity
	t.Run("Redis Connection", func(t *testing.T) {
		// Check if Redis port is accessible
		conn, err := net.DialTimeout("tcp", "localhost:6379", 5*time.Second)
		if err != nil {
			t.Skip("Redis not running on localhost:6379, skipping Redis tests")
			return
		}
		conn.Close()

		// Test Redis connection
		rdb := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15, // Use database 15 for testing
		})
		defer rdb.Close()

		// Test Redis ping
		ctx := context.Background()
		pong, err := rdb.Ping(ctx).Result()
		require.NoError(t, err, "Should ping Redis successfully")
		require.Equal(t, "PONG", pong, "Redis should respond with PONG")

		// Clear test database
		err = rdb.FlushDB(ctx).Err()
		require.NoError(t, err, "Should clear test Redis database")

		// Test Redis operations
		err = rdb.Set(ctx, "test-key", "test-value", time.Minute).Err()
		require.NoError(t, err, "Should set test key")

		value, err := rdb.Get(ctx, "test-key").Result()
		require.NoError(t, err, "Should get test key")
		require.Equal(t, "test-value", value, "Should retrieve correct value")

		t.Log("✅ Redis connection verified")
	})
	// Test 3: Check port availability
	t.Run("Port Availability", func(t *testing.T) {
		// Check if test port is available
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			t.Skipf("Test port 8081 is not available: %v", err)
			return
		}
		listener.Close()

		t.Log("✅ Test port 8081 is available")
	})

	t.Log("🎉 Test environment connectivity check completed")
}

// TestDatabaseConnectivity tests database connectivity without modifying anything
func TestDatabaseConnectivity(t *testing.T) {
	config := GetTestConfig()

	// Check if PostgreSQL is available
	conn, err := net.DialTimeout("tcp", "localhost:5432", 2*time.Second)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping database tests")
		return
	}
	conn.Close()

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		t.Skipf("Cannot connect to database: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Skipf("Cannot ping database: %v (database should exist before running tests)", err)
		return
	}

	t.Run("Basic Database Operations", func(t *testing.T) {
		// Test basic query
		var result int
		err := db.QueryRow("SELECT 1").Scan(&result)
		require.NoError(t, err, "Should execute basic query")
		require.Equal(t, 1, result, "Query should return 1")

		// Check if expected tables exist (created by application)
		tables := []string{"users", "oauth2_clients"}
		for _, table := range tables {
			var exists bool
			err := db.QueryRow(`
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' 
					AND table_name = $1
				)`, table).Scan(&exists)

			if err == nil && exists {
				t.Logf("✅ Table %s exists", table)
			} else {
				t.Logf("⚠️  Table %s does not exist - should be created by application", table)
			}
		}
	})
}

// TestRedisCRUD tests basic Redis operations
func TestRedisCRUD(t *testing.T) {
	// Check if Redis is available
	conn, err := net.DialTimeout("tcp", "localhost:6379", 2*time.Second)
	if err != nil {
		t.Skip("Redis not available, skipping Redis CRUD tests")
		return
	}
	conn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use test database
	})
	defer rdb.Close()

	ctx := context.Background()

	t.Run("String Operations", func(t *testing.T) {
		key := "test:string"
		value := "test-value"

		// Set
		err := rdb.Set(ctx, key, value, time.Minute).Err()
		require.NoError(t, err, "Should set string value")

		// Get
		result, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err, "Should get string value")
		require.Equal(t, value, result, "Value should match")

		// Delete
		err = rdb.Del(ctx, key).Err()
		require.NoError(t, err, "Should delete key")

		// Verify deletion
		_, err = rdb.Get(ctx, key).Result()
		require.Error(t, err, "Key should be deleted")
		require.Equal(t, redis.Nil, err, "Should return nil error")
	})

	t.Run("Hash Operations", func(t *testing.T) {
		key := "test:hash"
		field := "field1"
		value := "value1"

		// Set hash field
		err := rdb.HSet(ctx, key, field, value).Err()
		require.NoError(t, err, "Should set hash field")

		// Get hash field
		result, err := rdb.HGet(ctx, key, field).Result()
		require.NoError(t, err, "Should get hash field")
		require.Equal(t, value, result, "Hash field value should match")

		// Set expiration
		err = rdb.Expire(ctx, key, time.Minute).Err()
		require.NoError(t, err, "Should set expiration")

		// Check TTL
		ttl, err := rdb.TTL(ctx, key).Result()
		require.NoError(t, err, "Should get TTL")
		require.Greater(t, ttl, time.Duration(0), "TTL should be positive")

		// Clean up
		err = rdb.Del(ctx, key).Err()
		require.NoError(t, err, "Should delete hash")
	})
}
