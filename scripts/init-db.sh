#!/bin/bash

# 数据库初始化脚本
# 创建数据库和必要的配置

set -e

echo "🗄️ 初始化OAuth2数据库..."

# 数据库配置
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="1q2w3e4r"
DB_NAME="oauth2"

# 检查PostgreSQL是否运行
echo "🔍 检查PostgreSQL服务..."
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER 2>/dev/null; then
    echo "❌ PostgreSQL服务未运行，请先启动PostgreSQL"
    exit 1
fi

# 设置密码环境变量
export PGPASSWORD=$DB_PASSWORD

# 检查数据库是否存在
echo "🔍 检查数据库是否存在..."
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt 2>/dev/null | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo "✅ 数据库 '$DB_NAME' 已存在"
else
    echo "📝 创建数据库 '$DB_NAME'..."
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    echo "✅ 数据库 '$DB_NAME' 创建成功"
fi

# 测试连接
echo "🔗 测试数据库连接..."
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo "✅ 数据库连接成功"
else
    echo "❌ 数据库连接失败"
    exit 1
fi

# 检查Redis
echo "🔍 检查Redis服务..."
if redis-cli -h localhost -p 6379 ping 2>/dev/null | grep -q PONG; then
    echo "✅ Redis服务正常"
    
    # 测试Redis数据库2
    if redis-cli -h localhost -p 6379 -n 2 ping 2>/dev/null | grep -q PONG; then
        echo "✅ Redis数据库2可用"
    else
        echo "❌ Redis数据库2不可用"
        exit 1
    fi
else
    echo "❌ Redis服务未运行，请先启动Redis"
    exit 1
fi

echo ""
echo "🎉 数据库初始化完成！"
echo ""
echo "📝 连接信息："
echo "  PostgreSQL: postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
echo "  Redis: redis://localhost:6379/2"
echo ""
echo "🚀 现在可以启动应用了："
echo "  ./scripts/start-local.sh"
