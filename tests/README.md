# Flash OAuth2 E2E 测试指南

本指南详细说明如何运行 Flash OAuth2 服务器的端到端（E2E）测试套件。

## 📋 概述

E2E 测试套件提供了对 OAuth2 服务器的全面测试，包括：

- **OAuth2 授权码流程测试** - 完整的授权流程验证
- **手机号验证流程测试** - 短信验证码登录流程
- **API 端点测试** - 所有 REST API 端点的功能测试
- **JWKS 端点测试** - JWT 密钥集验证
- **安全功能测试** - 安全漏洞和边界情况测试
- **数据库集成测试** - PostgreSQL 操作验证
- **Redis 集成测试** - 缓存和会话存储测试
- **性能基准测试** - 性能指标和负载测试

## 🔧 环境要求

### 必需服务

在运行 E2E 测试之前，请确保以下服务正在运行：

1. **PostgreSQL** (端口 5432)

   ```bash
   # macOS (使用Homebrew)
   brew services start postgresql

   # Ubuntu/Debian
   sudo systemctl start postgresql

   # Docker
   docker run --name postgres -e POSTGRES_PASSWORD=1q2w3e4r -p 5432:5432 -d postgres
   ```

2. **Redis** (端口 6379)

   ```bash
   # macOS (使用Homebrew)
   brew services start redis

   # Ubuntu/Debian
   sudo systemctl start redis

   # Docker
   docker run --name redis -p 6379:6379 -d redis
   ```

### 依赖包

确保安装了测试依赖：

```bash
make test-e2e-deps
# 或者
go get github.com/stretchr/testify@latest
go mod tidy
```

## 🚀 快速开始

### 运行完整测试套件

```bash
# 运行所有测试（包括设置、测试、基准测试和覆盖率报告）
make test-e2e-all
```

### 运行基本测试

```bash
# 仅运行测试（假设环境已设置）
make test-e2e-run

# 快速测试（无环境设置/清理）
make test-e2e-quick
```

## 📊 测试命令详解

### 基本测试命令

| 命令                    | 描述          | 用途                       |
| ----------------------- | ------------- | -------------------------- |
| `make test-e2e-setup`   | 设置测试环境  | 创建测试数据库，清理 Redis |
| `make test-e2e-run`     | 运行 E2E 测试 | 执行所有功能测试           |
| `make test-e2e-quick`   | 快速测试      | 无环境设置的快速测试       |
| `make test-e2e-cleanup` | 清理测试资源  | 删除测试数据库和缓存       |

### 高级测试命令

| 命令                                   | 描述           | 示例                                                 |
| -------------------------------------- | -------------- | ---------------------------------------------------- |
| `make test-e2e-specific TEST=<测试名>` | 运行特定测试   | `make test-e2e-specific TEST=TestCompleteOAuth2Flow` |
| `make test-e2e-bench`                  | 性能基准测试   | 测试 token 生成和用户信息检索性能                    |
| `make test-e2e-coverage`               | 生成覆盖率报告 | 生成 HTML 覆盖率报告                                 |
| `make test-e2e-race`                   | 竞态条件检测   | 检测并发安全问题                                     |

### 开发工具命令

| 命令                          | 描述           | 用途                         |
| ----------------------------- | -------------- | ---------------------------- |
| `make test-e2e-validate`      | 验证测试文件   | 检查语法和格式               |
| `make test-e2e-view-coverage` | 查看覆盖率报告 | 在浏览器中打开 coverage.html |

## 📝 测试配置

### 环境变量

测试使用以下环境变量：

```bash
# 测试数据库连接
TEST_DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/oauth2_test?sslmode=disable"

# 测试Redis连接
TEST_REDIS_URL="redis://localhost:6379/15"

# 测试服务端口
TEST_PORT="8081"

# Gin运行模式
GIN_MODE="test"

# JWT密钥（测试用）
JWT_SECRET="test-jwt-secret-key-for-oauth2-server"
```

### 数据库配置

测试会自动创建和管理测试数据库：

- 数据库名：`oauth2_test`
- Redis 数据库：`15`（避免与主应用冲突）

## 🧪 测试套件详情

### 1. OAuth2 授权流程测试

**文件**: `tests/e2e_oauth2_test.go`

- `TestCompleteOAuth2Flow` - 完整授权码流程
- `TestPhoneAuthenticationFlow` - 手机号验证流程
- `TestJWKSEndpoint` - JWKS 端点测试
- `TestHealthEndpoint` - 健康检查端点
- `TestDocumentationEndpoint` - 文档端点测试

**运行示例**:

```bash
make test-e2e-specific TEST=TestCompleteOAuth2Flow
```

### 2. API 端点测试

**文件**: `tests/e2e_api_test.go`

- `TestAPIEndpoints` - 所有 API 端点功能测试
- `TestErrorHandling` - 错误处理测试
- `TestSecurityFeatures` - 安全功能测试
- `TestDatabaseIntegration` - 数据库集成测试
- `TestRedisIntegration` - Redis 集成测试
- `TestEdgeCases` - 边界情况测试

### 3. 性能基准测试

**测试项目**:

- Token 生成性能
- 用户信息检索性能
- 并发请求处理能力

**运行示例**:

```bash
make test-e2e-bench
```

## 📊 测试报告

### 覆盖率报告

生成测试覆盖率报告：

```bash
make test-e2e-coverage
```

这会生成：

- `coverage.out` - 覆盖率数据文件
- `coverage.html` - HTML 格式的覆盖率报告

查看报告：

```bash
make test-e2e-view-coverage
```

### 性能报告

基准测试会输出性能指标：

- 操作执行时间
- 内存分配统计
- 并发性能数据

## 🔍 故障排除

### 常见问题

1. **PostgreSQL 连接失败**

   ```
   确保PostgreSQL服务正在运行
   检查连接字符串中的密码是否正确
   验证数据库服务器地址和端口
   ```

2. **Redis 连接失败**

   ```
   确保Redis服务正在运行
   检查Redis是否监听在6379端口
   验证Redis配置允许本地连接
   ```

3. **端口冲突**

   ```
   确保测试端口8081未被占用
   检查是否有其他OAuth2服务实例在运行
   ```

4. **权限错误**
   ```
   确保测试脚本有执行权限
   chmod +x tests/run_e2e_tests.sh
   ```

### 调试模式

启用详细日志输出：

```bash
# 设置详细日志级别
export LOG_LEVEL=debug

# 运行测试
make test-e2e-run
```

### 清理环境

如果测试环境出现问题，完全清理：

```bash
# 停止所有服务
make stop

# 清理测试资源
make test-e2e-cleanup

# 重新设置环境
make test-e2e-setup
```

## 🚀 持续集成

### GitHub Actions

可以在 CI/CD 流水线中集成 E2E 测试：

```yaml
- name: Run E2E Tests
  run: |
    make test-e2e-deps
    make test-e2e-all
```

### Docker 测试环境

对于 CI 环境，可以使用 Docker Compose 启动测试依赖：

```bash
# 启动测试依赖服务
docker-compose -f docker-compose.test.yml up -d

# 运行测试
make test-e2e-run

# 清理
docker-compose -f docker-compose.test.yml down
```

## 📚 更多信息

- **项目文档**: [README.md](../README.md)
- **API 文档**: 运行服务器后访问 `http://localhost:8080/docs`
- **测试代码**: [tests/](./tests/) 目录
- **主要代码**: [主项目目录](../)

---

**提示**: 如果遇到问题，请确保所有依赖服务正在运行，并且已安装所有必需的 Go 包。使用`make test-e2e-validate`检查测试文件的语法错误。
