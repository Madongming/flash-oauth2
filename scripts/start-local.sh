#!/bin/bash

# 本地开发环境启动脚本
# 使用新的数据库和Redis配置

set -e

echo "🚀 启动Flash OAuth2本地开发环境..."

# 设置环境变量
export PORT=8080
export DATABASE_URL="postgres://postgres:1q2w3e4r@localhost:5432/oauth2?sslmode=disable"
export REDIS_URL="redis://localhost:6379/2"

echo "📝 当前配置："
echo "  - 端口: $PORT"
echo "  - 数据库: $DATABASE_URL"
echo "  - Redis: $REDIS_URL"
echo ""

# 检查PostgreSQL是否运行
echo "🔍 检查PostgreSQL连接..."
if ! psql "$DATABASE_URL" -c '\q' 2>/dev/null; then
    echo "❌ 无法连接到PostgreSQL。请确保："
    echo "  1. PostgreSQL服务正在运行"
    echo "  2. 数据库 'oauth2' 已创建"
    echo "  3. 用户 'postgres' 密码为 '1q2w3e4r'"
    echo ""
    echo "创建数据库的命令："
    echo "  createdb -h localhost -U postgres oauth2"
    echo ""
    exit 1
fi

# 检查Redis是否运行
echo "🔍 检查Redis连接..."
if ! redis-cli -h localhost -p 6379 -n 2 ping 2>/dev/null | grep -q PONG; then
    echo "❌ 无法连接到Redis。请确保："
    echo "  1. Redis服务正在运行"
    echo "  2. Redis可以访问数据库2"
    echo ""
    echo "启动Redis的命令："
    echo "  redis-server"
    echo ""
    exit 1
fi

echo "✅ 所有依赖服务正常"

# 构建应用
echo "🏗️ 构建应用..."
make build

# 启动应用
echo "🚀 启动OAuth2服务器..."
echo "访问地址: http://localhost:8080"
echo "健康检查: http://localhost:8080/health"
echo ""

./bin/flash-oauth2
