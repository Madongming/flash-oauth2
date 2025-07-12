#!/bin/bash

# 测试脚本 - 演示完整的OAuth2流程

set -e

BASE_URL="http://localhost:8080"
CLIENT_ID="default-client"
CLIENT_SECRET="default-secret"
REDIRECT_URI="http://localhost:3000/callback"
PHONE="13800138000"

echo "🧪 Flash OAuth2 API测试"
echo "========================"

# 1. 健康检查
echo "1️⃣ 健康检查..."
curl -s $BASE_URL/health | jq .
echo ""

# 2. 发送验证码
echo "2️⃣ 发送验证码到 $PHONE..."
SEND_CODE_RESPONSE=$(curl -s -X POST $BASE_URL/send-code \
  -H "Content-Type: application/json" \
  -d "{\"phone\": \"$PHONE\"}")
echo $SEND_CODE_RESPONSE | jq .

# 从日志中获取验证码（实际应用中应该通过短信获取）
echo ""
echo "🔍 请从服务器日志中查看验证码，然后手动测试以下步骤："
echo ""

# 3. 模拟用户登录获取授权码（需要手动操作）
echo "3️⃣ 模拟获取授权码："
echo "浏览器访问: $BASE_URL/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$REDIRECT_URI&scope=openid%20profile&state=test-state"
echo "然后在登录页面输入手机号和验证码"
echo ""

echo "4️⃣ 获取访问令牌（假设授权码为 test-auth-code）："
echo "curl -X POST $BASE_URL/token \\"
echo "  -H \"Content-Type: application/x-www-form-urlencoded\" \\"
echo "  -d \"grant_type=authorization_code&code=test-auth-code&redirect_uri=$REDIRECT_URI&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET\""
echo ""

echo "5️⃣ 获取用户信息（假设访问令牌为 ACCESS_TOKEN）："
echo "curl -X GET $BASE_URL/userinfo \\"
echo "  -H \"Authorization: Bearer ACCESS_TOKEN\""
echo ""

echo "6️⃣ 获取JWKS公钥："
curl -s $BASE_URL/.well-known/jwks.json | jq .
echo ""

echo "🎉 测试完成！请手动完成OAuth2流程测试。"
