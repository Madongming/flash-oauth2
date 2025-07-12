.PHONY: help build run test clean docker-build docker-run dev start stop logs
VERSION ?= $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2' )
TAG ?= $(VERSION)
FLASHX_IMG ?= flash-oauth2:${TAG}

# 默认目标
help:
	@echo "Flash OAuth2 认证服务器"
	@echo ""
	@echo "可用命令:"
	@echo ""
	@ec# E2E测试：运行OAuth2流程测试（需要完整环境）
test-e2e-oauth2:
	@echo "🔐 运行OAuth2流程测试..."
	@go test -v ./tests -run "TestCompleteOAuth2Flow|TestPhoneAuthenticationFlow|TestJWKSEndpoint" -timeout 60s

# E2E测试：运行完整API测试（需要数据库）
test-e2e-api-full:
	@echo "🌐 运行完整API测试..."
	@go test -v ./tests -run "TestAPIEndpoints|TestErrorHandling|TestSecurityFeatures" -timeout 60s🏗️  构建和运行:"
	@echo "  build        编译应用"
	@echo "  run          运行应用"
	@echo "  clean        清理构建文件"
	@echo ""
	@echo "🧪 测试命令:"
	@echo "  test         运行单元测试"
	@echo "  test-e2e     运行所有E2E测试（自动跳过不可用服务）"
	@echo "  test-e2e-quick 快速E2E测试（无外部依赖）"
	@echo "  test-e2e-connectivity 检查环境连通性"
	@echo "  test-e2e-unit 基础单元测试"
	@echo "  test-e2e-integration 集成测试"
	@echo "  test-e2e-oauth2 OAuth2流程测试"
	@echo "  test-e2e-bench 性能基准测试"
	@echo "  test-e2e-coverage 测试覆盖率报告"
	@echo ""
	@echo "🐳 Docker命令:"
	@echo "  dev          启动开发环境（Docker）"
	@echo "  start        启动所有服务"
	@echo "  stop         停止所有服务"
	@echo "  logs         查看日志"
	@echo "  docker-build 构建Docker镜像"
	@echo "  docker-run   运行Docker容器"
	@echo ""
	@echo "🔧 开发工具:"
	@echo "  deps         安装依赖"
	@echo "  lint         代码检查"
	@echo "  fmt          格式化代码"
	@echo "  generate-keys 生成RSA密钥对"
	@echo "  health       健康检查"

# 编译应用
build:
	@echo "🏗️ 编译应用..."
	go build -o bin/flash-oauth2 .

# 运行应用
run:
	@echo "🚀 运行应用..."
	go run .

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	docker-compose down --volumes --remove-orphans

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	go mod download

# 代码检查
lint:
	@echo "🔍 代码检查..."
	go vet ./...
	go fmt ./...

# 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	go fmt ./...

# 构建Docker镜像
docker-build:
	@echo "🐳 构建Docker镜像..."
	#docker build -t flash-oauth2 .
	docker buildx build --platform linux/amd64 . -t ${FLASHX_IMG} --load

# 运行Docker容器
docker-run: docker-build
	@echo "🐳 运行Docker容器..."
	docker run -p 8080:8080 flash-oauth2

# 启动开发环境
dev:
	@echo "🚀 启动开发环境..."
	./scripts/start-dev.sh

# 启动所有服务
start:
	@echo "🚀 启动所有服务..."
	docker-compose up -d

# 停止所有服务
stop:
	@echo "🛑 停止所有服务..."
	docker-compose down

# 查看日志
logs:
	@echo "📋 查看日志..."
	docker-compose logs -f

# 重新启动服务
restart: stop start

# 初始化数据库
init-db:
	@echo "🗄️ 初始化数据库..."
	./scripts/init-db.sh

# 测试API
test-api:
	@echo "🧪 测试API..."
	./scripts/test-api.sh

# 生成新的密钥对
generate-keys:
	@echo "🔑 生成新的RSA密钥对..."
	mkdir -p keys
	openssl genrsa -out keys/private.pem 2048
	openssl rsa -in keys/private.pem -pubout -out keys/public.pem
	@echo "密钥对已生成在 keys/ 目录"

# 健康检查
health:
	@echo "🏥 健康检查..."
	curl -f http://localhost:8080/health || exit 1

# 完整的开发环境设置
setup: deps start
	@echo "⏳ 等待服务启动..."
	sleep 15
	@echo "✅ 开发环境设置完成！"
	@echo "🔗 授权页面: http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=test"
	@echo "🔗 健康检查: http://localhost:8080/health"

# E2E 测试相关命令 - 纯Go实现，无外部依赖
.PHONY: test-e2e test-e2e-quick test-e2e-connectivity test-e2e-unit test-e2e-integration

# E2E测试：检查环境连通性（不修改任何数据）
test-e2e-connectivity:
	@echo "� 检查测试环境连通性..."
	@go test -v ./tests -run TestEnvironmentConnectivity -timeout 10s

# E2E测试：运行基础单元测试（无外部依赖）
test-e2e-unit:
	@echo "🧪 运行基础单元测试..."
	@go test -v ./tests -run TestBasic -timeout 10s

# E2E测试：运行集成测试（需要数据库/Redis，但会优雅跳过）
test-e2e-integration:
	@echo "� 运行集成测试..."
	@go test -v ./tests -run "TestWith" -timeout 30s

# E2E测试：运行OAuth2流程测试（需要完整环境）
test-e2e-oauth2:
	@echo "� 运行OAuth2流程测试..."
	@go test -v ./tests -run TestOAuth2 -timeout 60s

# E2E测试：快速测试（只运行无依赖的测试）
test-e2e-quick:
	@echo "⚡ 运行快速E2E测试..."
	@go test -v ./tests -run "TestBasic|TestEnvironmentConnectivity" -timeout 15s

# E2E测试：运行所有测试（优雅降级）
test-e2e:
	@echo "🚀 运行所有E2E测试（自动跳过不可用的服务）..."
	@go test -v ./tests -run "TestBasic|TestEnvironment|TestConnectivity|TestWith|TestOffline" -timeout 120s

# E2E测试：性能基准测试
test-e2e-bench:
	@echo "⚡ 运行性能基准测试..."
	@go test -v ./tests -bench=. -run=^$$ -benchtime=3s

# E2E测试：生成覆盖率报告
test-e2e-coverage:
	@echo "📊 生成测试覆盖率报告..."
	@go test -v ./tests -coverprofile=coverage.out -covermode=count
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📋 覆盖率报告已生成: coverage.html"

# E2E测试：竞态检测
test-e2e-race:
	@echo "🏃 运行E2E竞态检测测试..."
	@go test -v ./tests -race -timeout 60s

# E2E测试：验证测试文件语法
test-e2e-validate:
	@echo "✅ 验证E2E测试文件..."
	@go vet ./tests
	@go fmt ./tests

# E2E测试：查看覆盖率报告
test-e2e-view-coverage:
	@echo "🌐 在浏览器中打开覆盖率报告..."
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "请在浏览器中打开 coverage.html"
