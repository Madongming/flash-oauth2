# Flash OAuth2 认证服务器

一个基于 Go + Gin 框架实现的 OAuth2 + OpenID Connect 认证服务器。

## 功能特性

- ✅ 完整的 OAuth2.0 授权服务器实现
- ✅ OpenID Connect (OIDC) 支持
- ✅ JWT 访问令牌和 ID 令牌（非对称加密）
- ✅ 手机号 + 验证码认证
- ✅ 用户自动注册/登录（幂等操作）
- ✅ PostgreSQL 数据持久化
- ✅ Redis 缓存验证码等临时数据
- ✅ 现代化登录界面
- ✅ 令牌刷新机制
- ✅ 令牌内省端点
- ✅ JWKS 公钥端点

## 支持的 OAuth2 流程

- Authorization Code Flow (授权码流程)
- Refresh Token Flow (刷新令牌流程)

## 支持的 OpenID Connect 特性

- ID Token 生成
- UserInfo 端点
- JWKS 端点

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### 安装依赖

```bash
go mod tidy
```

### 环境变量配置

```bash
export PORT=8080
export DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable"
export REDIS_URL="redis://localhost:6379/2"
```

### 启动服务

```bash
go run main.go
```

服务将在 http://localhost:8080 启动

## API 端点

### OAuth2 端点

- `GET /authorize` - 授权端点
- `POST /token` - 令牌端点
- `POST /introspect` - 令牌内省端点

### OpenID Connect 端点

- `GET /userinfo` - 用户信息端点
- `GET /.well-known/jwks.json` - JSON Web Key Set

### 认证端点

- `POST /send-code` - 发送验证码
- `POST /login` - 用户登录

### 其他端点

- `GET /health` - 健康检查

## 使用示例

### 1. 获取授权码

```bash
# 浏览器访问
http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=xyz
```

### 2. 发送验证码

```bash
curl -X POST http://localhost:8080/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone": "13800138000"}'
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'phone=13800138000&code=123456&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=xyz'
```

### 4. 交换访问令牌

```bash
curl -X POST http://localhost:8080/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'grant_type=authorization_code&code=AUTH_CODE&redirect_uri=http://localhost:3000/callback&client_id=default-client&client_secret=default-secret'
```

### 5. 获取用户信息

```bash
curl -X GET http://localhost:8080/userinfo \
  -H "Authorization: Bearer ACCESS_TOKEN"
```

### 6. 刷新访问令牌

```bash
curl -X POST http://localhost:8080/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d 'grant_type=refresh_token&refresh_token=REFRESH_TOKEN&client_id=default-client&client_secret=default-secret'
```

## 数据库结构

项目包含以下数据表：

- `users` - 用户信息
- `oauth_clients` - OAuth2 客户端
- `auth_codes` - 授权码
- `access_tokens` - 访问令牌（JWT 令牌不需要存储，此表可用于审计）
- `refresh_tokens` - 刷新令牌

## 安全特性

- RSA 非对称加密签名 JWT 令牌
- 验证码有效期控制（5 分钟）
- 授权码有效期控制（10 分钟）
- 访问令牌有效期控制（1 小时）
- 刷新令牌有效期控制（30 天）
- 客户端认证
- 重定向 URI 验证

## 部署建议

1. 使用 HTTPS 部署生产环境
2. 配置适当的 CORS 策略
3. 使用环境变量管理敏感配置
4. 定期清理过期令牌
5. 监控和日志记录
6. 数据库连接池优化

## License

MIT License
