#!/bin/bash

# 启动开发环境脚本

set -e

echo "🚀 启动Flash OAuth2开发环境..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查docker-compose是否可用
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose未安装"
    exit 1
fi

# 停止现有的容器
echo "🛑 停止现有容器..."
docker-compose down

# 启动数据库和Redis
echo "📦 启动PostgreSQL和Redis..."
docker-compose up -d postgres redis

# 等待数据库和Redis就绪
echo "⏳ 等待数据库和Redis启动..."
sleep 10

# 检查数据库连接
echo "🔍 检查数据库连接..."
while ! docker-compose exec postgres pg_isready -U postgres -d oauth2 > /dev/null 2>&1; do
    echo "等待数据库启动..."
    sleep 2
done

# 检查Redis连接
echo "🔍 检查Redis连接..."
while ! docker-compose exec redis redis-cli ping > /dev/null 2>&1; do
    echo "等待Redis启动..."
    sleep 2
done

echo "✅ 数据库和Redis已启动"

# 构建并启动OAuth2服务器
echo "🏗️ 构建OAuth2服务器..."
docker-compose build oauth2-server

echo "🚀 启动OAuth2服务器..."
docker-compose up oauth2-server

echo "✅ OAuth2服务器已启动在 http://localhost:8080"
echo ""
echo "🔗 有用的链接："
echo "  - 授权页面: http://localhost:8080/authorize?response_type=code&client_id=default-client&redirect_uri=http://localhost:3000/callback&scope=openid%20profile&state=xyz"
echo "  - 健康检查: http://localhost:8080/health"
echo "  - JWKS端点: http://localhost:8080/.well-known/jwks.json"
echo ""
echo "📝 默认客户端信息："
echo "  - Client ID: default-client"
echo "  - Client Secret: default-secret"
echo "  - Redirect URI: http://localhost:3000/callback"
