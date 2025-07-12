# 🚀 Flash OAuth2 快速开始指南

## 5 分钟启动指南

### 前提条件

- Docker 和 Docker Compose
- Go 1.21+ (可选，如果要本地开发)

### 1. 克隆项目

```bash
git clone <your-repo>
cd flash-oauth2
```

### 2. 启动开发环境

```bash
# 一键启动所有服务
make setup
```

等待服务启动后，你会看到：

```
✅ 开发环境设置完成！
🔗 授权页面: http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=test
🔗 健康检查: http://localhost:8080/health
```

### 3. 测试 OAuth2 流程

#### 方法一：使用演示客户端

1. 在浏览器中打开 `examples/oauth2-client.html`
2. 点击"开始授权"按钮
3. 在登录页面输入手机号（如：13800138000）
4. 点击"发送验证码"，查看服务器日志获取验证码
5. 输入验证码并登录
6. 完成 OAuth2 流程

#### 方法二：手动测试

1. **访问授权页面**

   ```
   http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=test-state
   ```

2. **发送验证码**

   ```bash
   curl -X POST http://localhost:8080/send-code \
     -H "Content-Type: application/json" \
     -d '{"phone": "13800138000"}'
   ```

3. **查看验证码**

   ```bash
   # 查看服务器日志
   make logs
   ```

4. **在网页上登录**
   输入手机号和验证码，完成登录

5. **交换访问令牌**

   ```bash
   curl -X POST http://localhost:8080/token \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=authorization_code&code=YOUR_AUTH_CODE&redirect_uri=http://localhost:3000/callback&client_id=default-client&client_secret=default-secret"
   ```

6. **获取用户信息**
   ```bash
   curl -X GET http://localhost:8080/userinfo \
     -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
   ```

## 常用命令

```bash
# 查看帮助
make help

# 启动服务
make start

# 停止服务
make stop

# 查看日志
make logs

# 重新启动
make restart

# 健康检查
make health

# 测试API
make test-api
```

## 默认配置

| 项目       | 值                             |
| ---------- | ------------------------------ |
| 服务端口   | 8080                           |
| PostgreSQL | localhost:5432                 |
| Redis      | localhost:6379                 |
| 客户端 ID  | default-client                 |
| 客户端密钥 | default-secret                 |
| 重定向 URI | http://localhost:3000/callback |

## 环境变量配置

创建 `.env` 文件：

```bash
PORT=8080
DATABASE_URL=postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable
REDIS_URL=redis://localhost:6379/2
```

## 故障排除

### 1. 端口占用

```bash
# 检查端口占用
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 停止所有服务
make stop
```

### 2. 数据库连接失败

```bash
# 重新启动数据库
docker-compose restart postgres

# 检查数据库状态
docker-compose logs postgres
```

### 3. Redis 连接失败

```bash
# 重新启动Redis
docker-compose restart redis

# 检查Redis状态
docker-compose logs redis
```

### 4. 查看详细日志

```bash
# 查看所有服务日志
make logs

# 查看特定服务日志
docker-compose logs oauth2-server
```

## 生产部署建议

1. **环境变量配置**

   ```bash
   export DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"
   export REDIS_URL="redis://:password@host:6379"
   ```

2. **HTTPS 配置**

   - 使用反向代理（Nginx/Traefik）
   - 配置 SSL 证书

3. **安全加固**

   - 更改默认客户端密钥
   - 配置防火墙规则
   - 启用日志审计

4. **监控和备份**
   - 数据库定期备份
   - 应用监控和告警
   - 日志收集分析

## 下一步

- 查看 [完整 API 文档](README.md#api端点)
- 了解 [项目架构](ARCHITECTURE.md)
- 自定义 [OAuth2 客户端配置](README.md#数据库结构)
- 集成到你的应用中

## 支持和反馈

如果遇到问题，请：

1. 查看日志：`make logs`
2. 检查健康状态：`make health`
3. 查看文档：README.md 和 ARCHITECTURE.md
