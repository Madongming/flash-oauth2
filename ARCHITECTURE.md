# Flash OAuth2 项目架构说明

## 项目结构

```
flash-oauth2/
├── main.go                     # 应用入口点
├── config/
│   └── config.go              # 配置管理和RSA密钥生成
├── database/
│   └── database.go            # PostgreSQL数据库连接和迁移
├── redis_client/
│   └── redis.go               # Redis连接客户端
├── models/
│   └── models.go              # 数据模型定义
├── services/
│   ├── user_service.go        # 用户管理服务
│   ├── oauth_service.go       # OAuth2核心服务
│   └── jwt_service.go         # JWT令牌服务
├── handlers/
│   ├── handler.go             # 处理器基础结构
│   └── oauth.go               # OAuth2/OIDC端点处理器
├── middleware/
│   └── cors.go                # CORS中间件
├── routes/
│   └── routes.go              # 路由配置
├── templates/
│   └── login.html             # 登录页面模板
├── examples/
│   └── oauth2-client.html     # 客户端演示页面
├── scripts/
│   ├── init-db.sh             # 数据库初始化脚本
│   ├── start-dev.sh           # 开发环境启动脚本
│   └── test-api.sh            # API测试脚本
├── docker-compose.yml         # Docker Compose配置
├── Dockerfile                 # Docker镜像构建文件
├── Makefile                   # 构建和部署命令
└── README.md                  # 项目文档
```

## 核心特性实现

### 1. OAuth2.0 标准实现

- ✅ Authorization Code Flow
- ✅ Token Endpoint
- ✅ Introspection Endpoint
- ✅ 客户端认证
- ✅ 重定向 URI 验证
- ✅ Scope 管理

### 2. OpenID Connect (OIDC)

- ✅ ID Token 生成
- ✅ UserInfo Endpoint
- ✅ JWKS Endpoint (.well-known/jwks.json)

### 3. JWT 令牌系统

- ✅ RSA 非对称加密签名
- ✅ Access Token (访问令牌)
- ✅ ID Token (身份令牌)
- ✅ 令牌验证和解析
- ✅ 公钥分发

### 4. 用户认证系统

- ✅ 手机号+验证码认证
- ✅ 用户自动注册/登录（幂等操作）
- ✅ 验证码生成和验证
- ✅ Redis 缓存验证码

### 5. 数据持久化

- ✅ PostgreSQL 存储用户、客户端、令牌等数据
- ✅ Redis 存储临时数据（验证码、会话）
- ✅ 数据库自动迁移

## 安全设计

### 1. 加密和签名

- RSA 2048 位密钥对
- JWT 使用 RS256 算法签名
- 客户端可使用公钥验证令牌

### 2. 令牌生命周期

- 验证码：5 分钟
- 授权码：10 分钟
- 访问令牌：1 小时
- 刷新令牌：30 天
- ID 令牌：1 小时

### 3. 访问控制

- 客户端密钥验证
- 重定向 URI 白名单
- Scope 权限控制
- CORS 安全策略

## 部署方案

### 1. 本地开发

```bash
# 使用Docker Compose
make setup

# 或手动启动
make start
```

### 2. 生产部署

```bash
# 构建Docker镜像
make docker-build

# 使用环境变量配置
export DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"
export REDIS_URL="redis://host:6379"
export PORT="8080"

# 启动容器
docker run -d -p 8080:8080 \
  -e DATABASE_URL="$DATABASE_URL" \
  -e REDIS_URL="$REDIS_URL" \
  flash-oauth2
```

## API 端点总览

| 端点                     | 方法 | 描述                             |
| ------------------------ | ---- | -------------------------------- |
| `/authorize`             | GET  | OAuth2 授权端点                  |
| `/token`                 | POST | 令牌端点（授权码交换、刷新令牌） |
| `/userinfo`              | GET  | OIDC 用户信息端点                |
| `/introspect`            | POST | 令牌内省端点                     |
| `/.well-known/jwks.json` | GET  | JSON Web Key Set                 |
| `/login`                 | POST | 用户登录                         |
| `/send-code`             | POST | 发送验证码                       |
| `/health`                | GET  | 健康检查                         |

## 扩展能力

1. **多种认证方式**：可扩展支持密码、生物识别等
2. **多租户支持**：可为不同客户端配置不同策略
3. **审计日志**：可记录所有认证和授权事件
4. **速率限制**：可添加防刷机制
5. **联邦身份**：可对接第三方身份提供商

## 性能优化

1. **连接池**：数据库和 Redis 连接池
2. **缓存策略**：客户端信息缓存
3. **异步处理**：验证码发送异步化
4. **负载均衡**：支持水平扩展
