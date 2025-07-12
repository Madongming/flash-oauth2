# 配置修改总结

## 已修改的配置

### 数据库配置

- **旧配置**: `postgres://user:password@localhost/oauth2?sslmode=disable`
- **新配置**: `postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable`

### Redis 配置

- **旧配置**: `redis://localhost:6379`
- **新配置**: `redis://localhost:6379/2`

## 修改的文件

1. **config/config.go** - 更新默认配置值
2. **.env.example** - 更新环境变量示例
3. **docker-compose.yml** - 更新 Docker 配置
4. **scripts/start-dev.sh** - 更新开发脚本
5. **scripts/init-db.sh** - 更新数据库初始化脚本
6. **scripts/start-local.sh** - 新增本地启动脚本
7. **README.md** - 更新文档中的示例
8. **QUICKSTART.md** - 更新快速开始指南

## 使用方法

### 1. 本地开发（直接运行）

```bash
# 设置环境变量
export DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable"
export REDIS_URL="redis://localhost:6379/2"

# 初始化数据库
./scripts/init-db.sh

# 启动应用
./scripts/start-local.sh
```

### 2. Docker 开发

```bash
# 启动完整环境
make start

# 或使用开发脚本
./scripts/start-dev.sh
```

### 3. 环境变量文件

创建 `.env` 文件：

```bash
PORT=8080
DATABASE_URL=postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable
REDIS_URL=redis://localhost:6379/2
```

## 注意事项

1. **PostgreSQL**:

   - 用户名: `postgres`
   - 密码: `1q2w3e4r`
   - 数据库: `oauth2`
   - 端口: `5432`

2. **Redis**:

   - 主机: `localhost`
   - 端口: `6379`
   - 数据库: `2`

3. **安全提醒**:
   - 生产环境请更改密码
   - 建议使用更强的密码
   - 考虑启用 SSL 连接

## 验证配置

运行以下命令验证配置是否正确：

```bash
# 检查编译
make build

# 检查数据库连接
psql "postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable" -c '\q'

# 检查Redis连接
redis-cli -h localhost -p 6379 -n 2 ping
```

## 故障排除

如果遇到连接问题：

1. **PostgreSQL**:

   ```bash
   # 检查服务状态
   pg_isready -h localhost -p 5432 -U postgres

   # 创建数据库
   createdb -h localhost -U postgres oauth2
   ```

2. **Redis**:

   ```bash
   # 检查服务状态
   redis-cli ping

   # 测试特定数据库
   redis-cli -n 2 ping
   ```
