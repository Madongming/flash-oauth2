# Flash OAuth2 认证服务器

一个基于 Go + Gin 框架实现的完整 OAuth2 + OpenID Connect 认证服务器，支持手机号验证、应用管理和管理员认证。

## 📋 项目概述

Flash OAuth2 是一个功能完整的认证授权服务器，实现了现代认证标准和最佳实践：

### 🚀 核心功能

- **OAuth2.0 授权服务器** - 完整实现 RFC 6749 标准
- **OpenID Connect (OIDC)** - 身份认证层支持
- **JWT 令牌系统** - RSA 非对称加密签名
- **手机号认证** - 验证码登录机制
- **用户自动注册** - 幂等操作，无需预注册
- **应用管理平台** - 完整的 OAuth2 客户端管理
- **管理员系统** - 基于角色的访问控制
- **数据持久化** - PostgreSQL + Redis 双存储
- **容器化部署** - Docker 和 Docker Compose 支持

### 🔐 支持的认证流程

- **OAuth2 Authorization Code Flow** - 标准授权码流程
- **Refresh Token Flow** - 令牌刷新机制
- **Phone + Verification Code** - 手机号验证登录
- **Admin Authentication** - 管理员专用认证

### 🛡️ 安全特性

- RSA 2048 位密钥对 JWT 签名
- 短期访问令牌（1 小时）
- 长期刷新令牌（30 天）
- 验证码限时（5 分钟）
- 客户端认证和重定向 URI 验证
- CORS 安全策略

## 📁 项目结构

```
flash-oauth2/
├── main.go                     # 应用程序入口
├── config/
│   └── config.go              # 配置管理和 RSA 密钥生成
├── database/
│   └── database.go            # PostgreSQL 数据库连接和迁移
├── redis_client/
│   └── redis.go               # Redis 连接客户端
├── models/
│   └── models.go              # 数据模型定义
├── services/
│   ├── user_service.go        # 用户管理服务
│   ├── oauth_service.go       # OAuth2 核心服务
│   ├── jwt_service.go         # JWT 令牌服务
│   ├── sms_service.go         # 短信服务
│   └── app_management_service.go # 应用管理服务
├── handlers/
│   ├── handler.go             # 通用处理器
│   ├── oauth.go               # OAuth2/OIDC 端点处理器
│   └── app_management.go      # 应用管理处理器
├── middleware/
│   ├── admin_auth.go          # 管理员认证中间件
│   └── cors.go                # CORS 中间件
├── routes/
│   └── routes.go              # 路由配置
├── templates/
│   ├── login.gohtml           # 用户登录页面
│   ├── admin_login.gohtml     # 管理员登录页面
│   ├── dashboard.gohtml       # 管理仪表板
│   ├── app_details.gohtml     # 应用详情页面
│   └── register_developer.gohtml # 开发者注册页面
├── examples/
│   └── oauth2-client.html     # OAuth2 客户端演示
├── tests/
│   ├── e2e_oauth2_test.go     # OAuth2 流程测试
│   ├── e2e_api_test.go        # API 端点测试
│   ├── e2e_app_management_test.go # 应用管理测试
│   ├── e2e_test_helper.go     # 测试工具函数
│   ├── basic_test.go          # 基础测试
│   ├── config_test.go         # 配置测试
│   ├── test_main.go           # 测试主入口
│   └── test_data.go           # 测试数据
├── scripts/
│   ├── init-db.sh             # 数据库初始化脚本
│   ├── start-dev.sh           # 开发环境启动脚本
│   └── test-api.sh            # API 测试脚本
├── docker-compose.yml         # Docker Compose 配置
├── docker-compose.test.yml    # 测试环境 Docker 配置
├── Dockerfile                 # Docker 镜像构建文件
├── Makefile                   # 构建和部署命令
└── README.md                  # 项目文档
```

### 🗃️ 核心文件说明

| 文件/目录     | 功能说明                             |
| ------------- | ------------------------------------ |
| `main.go`     | 程序入口，初始化服务器和依赖         |
| `config/`     | 配置管理，RSA 密钥生成               |
| `database/`   | 数据库连接，迁移和初始化             |
| `models/`     | 数据模型定义（用户、客户端、令牌等） |
| `services/`   | 业务逻辑服务层                       |
| `handlers/`   | HTTP 请求处理器                      |
| `middleware/` | 中间件（认证、CORS 等）              |
| `routes/`     | 路由配置和映射                       |
| `templates/`  | HTML 模板文件                        |
| `tests/`      | 测试套件和工具                       |

## 🚀 部署说明

### 环境要求

- **Go 1.21+**
- **PostgreSQL 12+**
- **Redis 6+**
- **Docker & Docker Compose** (可选)

### 方式一：Docker Compose 部署（推荐）

1. **克隆项目**

   ```bash
   git clone <your-repo>
   cd flash-oauth2
   ```

2. **一键启动**

   ```bash
   make setup
   ```

   这将自动：

   - 启动 PostgreSQL 数据库
   - 启动 Redis 缓存服务
   - 构建并启动 OAuth2 服务器
   - 运行数据库迁移

3. **验证服务**
   ```bash
   make health
   # 或直接访问
   curl http://localhost:8080/health
   ```

### 方式二：本地开发部署

1. **安装依赖**

   ```bash
   go mod tidy
   ```

2. **启动依赖服务**

   ```bash
   # 启动 PostgreSQL 和 Redis
   docker-compose up -d postgres redis
   ```

3. **配置环境变量**

   ```bash
   export PORT=8080
   export DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable"
   export REDIS_URL="redis://localhost:6379/2"
   ```

4. **启动服务器**
   ```bash
   go run main.go
   ```

### 方式三：生产环境部署

1. **构建应用**

   ```bash
   make build
   ```

2. **配置生产环境变量**

   ```bash
   export PORT=8080
   export DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"
   export REDIS_URL="redis://:password@host:6379"
   export SMS_ENABLED=true
   export SMS_ACCESS_KEY_ID="your-key"
   export SMS_ACCESS_KEY_SECRET="your-secret"
   ```

3. **运行应用**
   ```bash
   ./bin/flash-oauth2
   ```

### 服务依赖

| 服务          | 端口 | 用途       | 配置                     |
| ------------- | ---- | ---------- | ------------------------ |
| PostgreSQL    | 5432 | 主数据库   | 存储用户、客户端、令牌等 |
| Redis         | 6379 | 缓存数据库 | 验证码、会话数据         |
| OAuth2 Server | 8080 | 认证服务器 | 主服务端口               |

### Docker 配置说明

- **开发环境**: `docker-compose.yml`
- **测试环境**: `docker-compose.test.yml` (不同端口避免冲突)
- **生产镜像**: `Dockerfile` (多阶段构建)

## 🧪 测试说明

### 测试套件概览

项目包含完整的端到端（E2E）测试套件，覆盖所有核心功能：

| 测试类型        | 文件                         | 测试内容                          |
| --------------- | ---------------------------- | --------------------------------- |
| OAuth2 流程测试 | `e2e_oauth2_test.go`         | 完整授权流程、手机验证、JWKS 端点 |
| API 端点测试    | `e2e_api_test.go`            | 所有 REST API、错误处理、安全验证 |
| 应用管理测试    | `e2e_app_management_test.go` | 开发者注册、应用管理、密钥管理    |
| 基础功能测试    | `basic_test.go`              | 基础验证（无外部依赖）            |
| 配置测试        | `config_test.go`             | 配置管理和测试数据                |

### 快速测试

```bash
# 基础测试（无需数据库）
make test-e2e-quick

# 完整测试套件
make test-e2e-all

# 特定测试
make test-e2e-specific TEST=TestCompleteOAuth2Flow
```

### 测试命令详解

| 命令                     | 说明          | 用途                         |
| ------------------------ | ------------- | ---------------------------- |
| `make test-e2e-all`      | 完整测试套件  | 包含设置、测试、基准、覆盖率 |
| `make test-e2e-run`      | 运行 E2E 测试 | 主要功能测试                 |
| `make test-e2e-quick`    | 快速测试      | 无环境设置的快速验证         |
| `make test-e2e-bench`    | 性能测试      | 基准测试和性能分析           |
| `make test-e2e-coverage` | 覆盖率报告    | 测试覆盖率分析               |

### 测试环境配置

测试环境自动配置：

- **测试数据库**: `oauth2_test`
- **Redis 数据库**: `15` (避免与主应用冲突)
- **测试端口**: `8081`

### 性能基准测试

```bash
# 运行性能测试
make test-e2e-bench

# 示例输出：
# BenchmarkTokenGeneration-8    1000    1.2ms/op    512 B/op    8 allocs/op
# BenchmarkUserInfoRetrieval-8  2000    0.8ms/op    256 B/op    4 allocs/op
```

### 测试覆盖率

```bash
# 生成覆盖率报告
make test-e2e-coverage

# 在浏览器中查看
make test-e2e-view-coverage
```

## 📖 使用说明

### API 端点总览

| 类型               | 端点                     | 方法     | 说明             |
| ------------------ | ------------------------ | -------- | ---------------- |
| **OAuth2 核心**    | `/authorize`             | GET      | 授权端点         |
|                    | `/token`                 | POST     | 令牌交换端点     |
|                    | `/introspect`            | POST     | 令牌内省端点     |
| **OpenID Connect** | `/userinfo`              | GET      | 用户信息端点     |
|                    | `/.well-known/jwks.json` | GET      | JSON Web Key Set |
| **用户认证**       | `/login`                 | POST     | 用户登录         |
|                    | `/send-code`             | POST     | 发送验证码       |
| **管理员**         | `/admin/login`           | GET/POST | 管理员登录       |
|                    | `/admin/dashboard`       | GET      | 管理仪表板       |
| **应用管理**       | `/api/admin/apps`        | GET/POST | 应用管理         |
|                    | `/api/admin/developers`  | POST     | 开发者注册       |
| **其他**           | `/health`                | GET      | 健康检查         |

### 完整 OAuth2 流程示例

#### 1. 客户端注册

首先需要管理员注册 OAuth2 客户端：

```bash
# 访问管理员登录页面
http://localhost:8080/admin/login

# 使用默认管理员账户登录
# 手机号: admin
# 从服务器日志获取验证码
```

#### 2. 发起授权请求

```bash
# 浏览器访问授权端点
http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=xyz
```

#### 3. 用户认证

发送验证码：

```bash
curl -X POST http://localhost:8080/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone": "13800138000"}'
```

在网页登录表单中输入手机号和验证码完成登录。

#### 4. 交换访问令牌

```bash
curl -X POST http://localhost:8080/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'grant_type=authorization_code&code=YOUR_AUTH_CODE&redirect_uri=http://localhost:3000/callback&client_id=default-client&client_secret=default-secret'
```

响应示例：

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiI...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "refresh_token_string",
  "id_token": "eyJhbGciOiJSUzI1NiI...",
  "scope": "openid profile"
}
```

#### 5. 获取用户信息

```bash
curl -X GET http://localhost:8080/userinfo \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

响应示例：

```json
{
  "sub": "123",
  "phone": "13800138000"
}
```

#### 6. 刷新访问令牌

```bash
curl -X POST http://localhost:8080/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'grant_type=refresh_token&refresh_token=REFRESH_TOKEN&client_id=default-client&client_secret=default-secret'
```

### 应用管理流程

#### 1. 注册开发者

```bash
curl -X POST http://localhost:8080/api/admin/developers \
  -H "Content-Type: application/json" \
  -H "Cookie: admin_session=SESSION_TOKEN" \
  -d '{
    "name": "开发者名称",
    "email": "dev@example.com",
    "phone": "13800138000"
  }'
```

#### 2. 注册应用

```bash
curl -X POST http://localhost:8080/api/admin/apps \
  -H "Content-Type: application/json" \
  -H "Cookie: admin_session=SESSION_TOKEN" \
  -d '{
    "name": "我的应用",
    "description": "应用描述",
    "callback_urls": ["https://myapp.com/callback"],
    "developer_id": "DEVELOPER_ID"
  }'
```

#### 3. 生成密钥对

```bash
curl -X POST http://localhost:8080/api/admin/apps/APP_ID/keys \
  -H "Cookie: admin_session=SESSION_TOKEN"
```

### 默认配置

| 配置项     | 默认值                                             | 说明               |
| ---------- | -------------------------------------------------- | ------------------ |
| 服务端口   | 8080                                               | HTTP 服务端口      |
| 数据库     | postgres://postgres:1q2w3e4r@localhost:5432/oauth2 | PostgreSQL 连接    |
| Redis      | redis://localhost:6379/2                           | Redis 连接         |
| 客户端 ID  | default-client                                     | 默认 OAuth2 客户端 |
| 客户端密钥 | default-secret                                     | 默认客户端密钥     |
| 重定向 URI | http://localhost:3000/callback                     | 默认回调地址       |
| 管理员账户 | admin                                              | 默认管理员手机号   |

### 令牌生命周期

| 令牌类型 | 生命周期 | 用途        |
| -------- | -------- | ----------- |
| 验证码   | 5 分钟   | 手机号验证  |
| 授权码   | 10 分钟  | OAuth2 授权 |
| 访问令牌 | 1 小时   | API 访问    |
| 刷新令牌 | 30 天    | 令牌刷新    |
| ID 令牌  | 1 小时   | 身份验证    |

## 📚 常用命令

### 开发命令

```bash
# 查看帮助
make help

# 编译应用
make build

# 运行应用
make run

# 运行测试
make test

# 代码检查
make lint

# 格式化代码
make fmt
```

### 服务管理

```bash
# 启动开发环境
make dev

# 启动所有服务
make start

# 停止所有服务
make stop

# 重启服务
make restart

# 查看日志
make logs

# 健康检查
make health
```

### Docker 命令

```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 清理构建文件
make clean
```

### 数据库管理

```bash
# 初始化数据库
make init-db

# 测试 API
make test-api

# 生成 RSA 密钥对
make generate-keys
```

## 🔧 环境变量配置

### 基础配置

```bash
# 服务配置
PORT=8080                                    # 服务端口
GIN_MODE=release                            # Gin 运行模式

# 数据库配置
DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"
REDIS_URL="redis://:password@host:6379/2"

# JWT 配置 (自动生成，也可手动指定)
JWT_ISSUER="flash-oauth2"                   # JWT 发行者
```

### 短信服务配置（可选）

```bash
# 阿里云短信服务
SMS_ENABLED=true                            # 启用短信发送
SMS_ACCESS_KEY_ID="your-access-key"         # 阿里云 Access Key ID
SMS_ACCESS_KEY_SECRET="your-secret"         # 阿里云 Access Key Secret
SMS_SIGN_NAME="your-signature"              # 短信签名
SMS_TEMPLATE_CODE="SMS_123456789"           # 短信模板代码
```

### 生产环境配置示例

创建 `.env` 文件：

```bash
# 生产环境配置
PORT=8080
GIN_MODE=release

# 数据库配置（使用 SSL）
DATABASE_URL="postgres://oauth2:secure_password@db.example.com:5432/oauth2_prod?sslmode=require"
REDIS_URL="redis://:redis_password@cache.example.com:6379/0"

# 启用短信服务
SMS_ENABLED=true
SMS_ACCESS_KEY_ID="LTAI..."
SMS_ACCESS_KEY_SECRET="xxx..."
SMS_SIGN_NAME="MyApp"
SMS_TEMPLATE_CODE="SMS_123456789"
```

## 🛠️ 故障排除

### 常见问题

#### 1. 端口占用

```bash
# 检查端口占用
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 停止所有服务
make stop
```

#### 2. 数据库连接失败

```bash
# 检查 PostgreSQL 状态
docker-compose logs postgres

# 重启数据库
docker-compose restart postgres

# 手动连接测试
psql postgres://postgres:1q2w3e4r@localhost:5432/oauth2
```

#### 3. Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose logs redis

# 重启 Redis
docker-compose restart redis

# 手动连接测试
redis-cli -h localhost -p 6379
```

#### 4. 验证码收不到

```bash
# 检查服务器日志中的验证码
make logs

# 确认 SMS 配置
echo $SMS_ENABLED
echo $SMS_ACCESS_KEY_ID
```

### 调试命令

```bash
# 查看详细日志
export LOG_LEVEL=debug && go run main.go

# 查看特定服务日志
docker-compose logs -f oauth2-server

# 检查数据库连接
make init-db

# 测试 API 端点
make test-api

# 运行健康检查
make health
```

### 数据库问题

```bash
# 重新创建数据库
docker-compose down -v
docker-compose up -d postgres
make init-db

# 查看数据库表
psql postgres://postgres:1q2w3e4r@localhost:5432/oauth2 -c "\dt"

# 查看用户数据
psql postgres://postgres:1q2w3e4r@localhost:5432/oauth2 -c "SELECT * FROM users;"
```

### 性能优化建议

1. **数据库优化**

   - 使用连接池
   - 创建适当索引
   - 定期清理过期令牌

2. **Redis 优化**

   - 配置内存限制
   - 启用持久化
   - 监控连接数

3. **应用优化**
   - 启用 Gin 生产模式
   - 配置日志级别
   - 使用反向代理

## 📄 许可证

MIT License

---

## 🚀 快速开始

```bash
# 1. 克隆项目
git clone <your-repo>
cd flash-oauth2

# 2. 一键启动
make setup

# 3. 验证服务
curl http://localhost:8080/health

# 4. 访问管理界面
open http://localhost:8080/admin/login
```

**就是这么简单！** 🎉

如需更多帮助，请查看项目内的具体文档或提交 Issue。
