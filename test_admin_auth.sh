#!/bin/bash

# Test script for admin authentication functionality
# This script tests the admin login and authentication system

echo "=== Flash OAuth2 管理员认证测试 ==="
echo

# 1. Test admin login page access
echo "1. 测试管理员登录页面访问..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/admin/login
echo

# 2. Test dashboard access without authentication (should redirect to login)
echo "2. 测试未认证访问仪表板（应重定向到登录页面）..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" -L http://localhost:8080/admin/dashboard
echo

# 3. Test admin API access without authentication (should return 403 or redirect)
echo "3. 测试未认证访问管理API（应返回403或重定向）..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/api/admin/apps
echo

# 4. Test sending verification code to admin phone
echo "4. 测试发送验证码到管理员手机..."
curl -X POST http://localhost:8080/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone": "admin"}' \
  -s -o /dev/null -w "HTTP Status: %{http_code}\n"
echo

echo "=== 测试说明 ==="
echo "- 登录页面应返回 200"
echo "- 未认证访问仪表板应重定向（302）"
echo "- 未认证访问API应返回403或重定向"
echo "- 发送验证码应返回200"
echo
echo "启动服务器后运行此脚本："
echo "  ./test_admin_auth.sh"
echo
echo "管理员登录信息："
echo "  手机号: admin"
echo "  验证码: 从服务器控制台获取"
echo "  登录后访问: http://localhost:8080/admin/dashboard"
