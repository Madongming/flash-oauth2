#!/bin/bash

# E2E Test Runner Script for Flash OAuth2 Server
# This script sets up the test environment and runs comprehensive end-to-end tests

set -e  # Exit on any error

echo "🚀 Starting Flash OAuth2 E2E Tests"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_DB_NAME="oauth2_test"
TEST_REDIS_DB="15"
TEST_PORT="8081"

echo -e "${BLUE}📋 Test Configuration:${NC}"
echo "  Database: postgres://postgres:1q2w3e4r@localhost:5432/${TEST_DB_NAME}?sslmode=disable"
echo "  Redis: redis://localhost:6379/${TEST_REDIS_DB}"
echo "  Test Port: ${TEST_PORT}"
echo ""

# Function to check if a service is running
check_service() {
    local service=$1
    local port=$2
    
    if nc -z localhost $port 2>/dev/null; then
        echo -e "${GREEN}✅ $service is running on port $port${NC}"
        return 0
    else
        echo -e "${RED}❌ $service is not running on port $port${NC}"
        return 1
    fi
}

# Function to setup test database
setup_test_database() {
    echo -e "${YELLOW}🗄️  Setting up test database...${NC}"
    
    # Check if PostgreSQL is running
    if ! check_service "PostgreSQL" 5432; then
        echo -e "${RED}Please start PostgreSQL service${NC}"
        echo "  macOS: brew services start postgresql"
        echo "  Ubuntu: sudo systemctl start postgresql"
        echo "  Docker: docker run --name postgres -e POSTGRES_PASSWORD=1q2w3e4r -p 5432:5432 -d postgres"
        exit 1
    fi
    
    # Create test database if it doesn't exist
    echo "Creating test database if not exists..."
    PGPASSWORD=1q2w3e4r psql -h localhost -U postgres -c "CREATE DATABASE ${TEST_DB_NAME};" 2>/dev/null || echo "Database already exists"
    
    echo -e "${GREEN}✅ Test database ready${NC}"
}

# Function to setup Redis for testing
setup_test_redis() {
    echo -e "${YELLOW}🔴 Setting up test Redis...${NC}"
    
    # Check if Redis is running
    if ! check_service "Redis" 6379; then
        echo -e "${RED}Please start Redis service${NC}"
        echo "  macOS: brew services start redis"
        echo "  Ubuntu: sudo systemctl start redis"
        echo "  Docker: docker run --name redis -p 6379:6379 -d redis"
        exit 1
    fi
    
    # Clear test Redis database
    echo "Clearing test Redis database ${TEST_REDIS_DB}..."
    redis-cli -n ${TEST_REDIS_DB} FLUSHDB >/dev/null
    
    echo -e "${GREEN}✅ Test Redis ready${NC}"
}

# Function to run tests
run_tests() {
    echo -e "${YELLOW}🧪 Running E2E tests...${NC}"
    
    # Export test environment variables
    export TEST_DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/${TEST_DB_NAME}?sslmode=disable"
    export TEST_REDIS_URL="redis://localhost:6379/${TEST_REDIS_DB}"
    export TEST_PORT="${TEST_PORT}"
    export GIN_MODE="test"
    
    # Run different test suites
    echo -e "${BLUE}Running OAuth2 flow tests...${NC}"
    go test -v ./tests -run TestCompleteOAuth2Flow -timeout 30s
    
    echo -e "${BLUE}Running phone authentication tests...${NC}"
    go test -v ./tests -run TestPhoneAuthenticationFlow -timeout 15s
    
    echo -e "${BLUE}Running API endpoint tests...${NC}"
    go test -v ./tests -run TestAPIEndpoints -timeout 20s
    
    echo -e "${BLUE}Running JWKS tests...${NC}"
    go test -v ./tests -run TestJWKSEndpoint -timeout 10s
    
    echo -e "${BLUE}Running health check tests...${NC}"
    go test -v ./tests -run TestHealthEndpoint -timeout 5s
    
    echo -e "${BLUE}Running documentation tests...${NC}"
    go test -v ./tests -run TestDocumentationEndpoint -timeout 5s
    
    echo -e "${BLUE}Running error handling tests...${NC}"
    go test -v ./tests -run TestErrorHandling -timeout 15s
    
    echo -e "${BLUE}Running database integration tests...${NC}"
    go test -v ./tests -run TestDatabaseIntegration -timeout 20s
    
    echo -e "${BLUE}Running Redis integration tests...${NC}"
    go test -v ./tests -run TestRedisIntegration -timeout 15s
    
    echo -e "${BLUE}Running security tests...${NC}"
    go test -v ./tests -run TestSecurityFeatures -timeout 15s
    
    echo -e "${BLUE}Running edge case tests...${NC}"
    go test -v ./tests -run TestEdgeCases -timeout 10s
}

# Function to run performance benchmarks
run_benchmarks() {
    echo -e "${YELLOW}⚡ Running performance benchmarks...${NC}"
    
    export TEST_DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/${TEST_DB_NAME}?sslmode=disable"
    export TEST_REDIS_URL="redis://localhost:6379/${TEST_REDIS_DB}"
    export GIN_MODE="test"
    
    echo -e "${BLUE}Benchmarking token generation...${NC}"
    go test -v ./tests -bench=BenchmarkTokenGeneration -run=^$ -benchtime=5s
    
    echo -e "${BLUE}Benchmarking user info retrieval...${NC}"
    go test -v ./tests -bench=BenchmarkUserInfoRetrieval -run=^$ -benchtime=5s
}

# Function to generate test coverage report
generate_coverage() {
    echo -e "${YELLOW}📊 Generating test coverage report...${NC}"
    
    export TEST_DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/${TEST_DB_NAME}?sslmode=disable"
    export TEST_REDIS_URL="redis://localhost:6379/${TEST_REDIS_DB}"
    export GIN_MODE="test"
    
    # Run tests with coverage
    go test -v ./tests -coverprofile=coverage.out -covermode=count
    
    # Generate HTML coverage report
    go tool cover -html=coverage.out -o coverage.html
    
    # Display coverage summary
    go tool cover -func=coverage.out | tail -1
    
    echo -e "${GREEN}✅ Coverage report generated: coverage.html${NC}"
}

# Function to cleanup test resources
cleanup() {
    echo -e "${YELLOW}🧹 Cleaning up test resources...${NC}"
    
    # Clean test database
    PGPASSWORD=1q2w3e4r psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS ${TEST_DB_NAME};" 2>/dev/null || true
    
    # Clear test Redis database
    redis-cli -n ${TEST_REDIS_DB} FLUSHDB >/dev/null 2>&1 || true
    
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Main execution
main() {
    local command=${1:-"all"}
    
    case $command in
        "setup")
            setup_test_database
            setup_test_redis
            ;;
        "test")
            setup_test_database
            setup_test_redis
            run_tests
            ;;
        "bench")
            setup_test_database
            setup_test_redis
            run_benchmarks
            ;;
        "coverage")
            setup_test_database
            setup_test_redis
            generate_coverage
            ;;
        "cleanup")
            cleanup
            ;;
        "all")
            setup_test_database
            setup_test_redis
            run_tests
            run_benchmarks
            generate_coverage
            echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
            ;;
        *)
            echo "Usage: $0 {setup|test|bench|coverage|cleanup|all}"
            echo ""
            echo "Commands:"
            echo "  setup    - Setup test environment"
            echo "  test     - Run E2E tests"
            echo "  bench    - Run performance benchmarks"
            echo "  coverage - Generate test coverage report"
            echo "  cleanup  - Clean up test resources"
            echo "  all      - Run everything (default)"
            exit 1
            ;;
    esac
}

# Trap cleanup on script exit
trap cleanup EXIT

# Run main function with all arguments
main "$@"
