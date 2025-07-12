#!/bin/bash

# 测试环境检查脚本
# 检查所有必需的服务和依赖是否正确配置

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Flash OAuth2 E2E 测试环境检查${NC}"
echo "================================="

# 检查Go环境
check_go() {
    echo -e "\n${YELLOW}📦 检查Go环境...${NC}"
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go未安装${NC}"
        return 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go已安装: ${GO_VERSION}${NC}"
    
    # 检查Go模块
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod文件不存在${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ Go模块配置正确${NC}"
    return 0
}

# 检查PostgreSQL连接
check_postgresql() {
    echo -e "\n${YELLOW}🐘 检查PostgreSQL连接...${NC}"
    
    # 检查PostgreSQL是否运行
    if ! nc -z localhost 5432 2>/dev/null; then
        echo -e "${RED}❌ PostgreSQL未在端口5432运行${NC}"
        echo -e "${YELLOW}💡 启动建议:${NC}"
        echo "  macOS: brew services start postgresql"
        echo "  Ubuntu: sudo systemctl start postgresql"
        echo "  Docker: docker run --name postgres -e POSTGRES_PASSWORD=1q2w3e4r -p 5432:5432 -d postgres"
        return 1
    fi
    
    echo -e "${GREEN}✅ PostgreSQL正在运行${NC}"
    
    # 检查数据库连接
    if command -v psql &> /dev/null; then
        if PGPASSWORD=1q2w3e4r psql -h localhost -U postgres -c '\q' 2>/dev/null; then
            echo -e "${GREEN}✅ PostgreSQL连接成功${NC}"
        else
            echo -e "${RED}❌ PostgreSQL连接失败${NC}"
            echo -e "${YELLOW}💡 请检查用户名和密码配置${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  psql未安装，跳过连接测试${NC}"
    fi
    
    return 0
}

# 检查Redis连接
check_redis() {
    echo -e "\n${YELLOW}🔴 检查Redis连接...${NC}"
    
    # 检查Redis是否运行
    if ! nc -z localhost 6379 2>/dev/null; then
        echo -e "${RED}❌ Redis未在端口6379运行${NC}"
        echo -e "${YELLOW}💡 启动建议:${NC}"
        echo "  macOS: brew services start redis"
        echo "  Ubuntu: sudo systemctl start redis"
        echo "  Docker: docker run --name redis -p 6379:6379 -d redis"
        return 1
    fi
    
    echo -e "${GREEN}✅ Redis正在运行${NC}"
    
    # 检查Redis连接
    if command -v redis-cli &> /dev/null; then
        if redis-cli ping >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Redis连接成功${NC}"
        else
            echo -e "${RED}❌ Redis连接失败${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  redis-cli未安装，跳过连接测试${NC}"
    fi
    
    return 0
}

# 检查网络工具
check_network_tools() {
    echo -e "\n${YELLOW}🔧 检查网络工具...${NC}"
    
    if ! command -v nc &> /dev/null; then
        echo -e "${RED}❌ netcat (nc)未安装${NC}"
        echo -e "${YELLOW}💡 安装建议:${NC}"
        echo "  macOS: brew install netcat"
        echo "  Ubuntu: sudo apt-get install netcat"
        return 1
    fi
    
    echo -e "${GREEN}✅ 网络工具可用${NC}"
    return 0
}

# 检查测试依赖
check_test_dependencies() {
    echo -e "\n${YELLOW}🧪 检查测试依赖...${NC}"
    
    # 检查testify依赖
    if ! go list -m github.com/stretchr/testify >/dev/null 2>&1; then
        echo -e "${RED}❌ testify依赖未安装${NC}"
        echo -e "${YELLOW}💡 安装命令: go get github.com/stretchr/testify@latest${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ 测试依赖已安装${NC}"
    
    # 检查测试文件
    if [ ! -d "tests" ]; then
        echo -e "${RED}❌ tests目录不存在${NC}"
        return 1
    fi
    
    if [ ! -f "tests/e2e_test_helper.go" ]; then
        echo -e "${RED}❌ E2E测试助手文件不存在${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ 测试文件结构正确${NC}"
    return 0
}

# 检查端口可用性
check_ports() {
    echo -e "\n${YELLOW}🚪 检查端口可用性...${NC}"
    
    TEST_PORT=8081
    
    if nc -z localhost $TEST_PORT 2>/dev/null; then
        echo -e "${RED}❌ 测试端口${TEST_PORT}已被占用${NC}"
        echo -e "${YELLOW}💡 请停止占用端口的服务或更改TEST_PORT环境变量${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ 测试端口${TEST_PORT}可用${NC}"
    return 0
}

# 检查文件权限
check_permissions() {
    echo -e "\n${YELLOW}🔐 检查文件权限...${NC}"
    
    if [ ! -x "tests/run_e2e_tests.sh" ]; then
        echo -e "${RED}❌ 测试脚本没有执行权限${NC}"
        echo -e "${YELLOW}💡 修复命令: chmod +x tests/run_e2e_tests.sh${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ 文件权限正确${NC}"
    return 0
}

# 运行快速连接测试
run_quick_tests() {
    echo -e "\n${YELLOW}⚡ 运行快速连接测试...${NC}"
    
    # 测试数据库连接
    if command -v psql &> /dev/null && PGPASSWORD=1q2w3e4r psql -h localhost -U postgres -c 'SELECT 1;' >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL连接测试通过${NC}"
    else
        echo -e "${RED}❌ PostgreSQL连接测试失败${NC}"
        return 1
    fi
    
    # 测试Redis连接
    if command -v redis-cli &> /dev/null && redis-cli ping >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Redis连接测试通过${NC}"
    else
        echo -e "${RED}❌ Redis连接测试失败${NC}"
        return 1
    fi
    
    return 0
}

# 显示环境信息
show_environment_info() {
    echo -e "\n${BLUE}📋 环境信息:${NC}"
    echo "  操作系统: $(uname -s)"
    echo "  架构: $(uname -m)"
    echo "  Go版本: $(go version 2>/dev/null | awk '{print $3}' || echo '未安装')"
    echo "  PostgreSQL: $(nc -z localhost 5432 && echo '运行中' || echo '未运行')"
    echo "  Redis: $(nc -z localhost 6379 && echo '运行中' || echo '未运行')"
    echo "  测试端口: 8081 $(nc -z localhost 8081 && echo '(占用)' || echo '(可用)')"
}

# 主函数
main() {
    local errors=0
    
    # 运行所有检查
    check_go || ((errors++))
    check_network_tools || ((errors++))
    check_postgresql || ((errors++))
    check_redis || ((errors++))
    check_test_dependencies || ((errors++))
    check_ports || ((errors++))
    check_permissions || ((errors++))
    
    # 显示环境信息
    show_environment_info
    
    # 运行快速测试
    if [ $errors -eq 0 ]; then
        run_quick_tests || ((errors++))
    fi
    
    # 显示结果
    echo -e "\n================================="
    if [ $errors -eq 0 ]; then
        echo -e "${GREEN}🎉 环境检查通过！可以运行E2E测试${NC}"
        echo -e "\n${BLUE}下一步:${NC}"
        echo "  make test-e2e-all    # 运行完整测试套件"
        echo "  make test-e2e-quick  # 快速测试"
        echo "  ./tests/run_e2e_tests.sh test  # 使用脚本运行"
        exit 0
    else
        echo -e "${RED}❌ 发现 $errors 个问题，请修复后再运行测试${NC}"
        echo -e "\n${BLUE}修复建议:${NC}"
        echo "  1. 启动必需的服务 (PostgreSQL, Redis)"
        echo "  2. 安装缺失的依赖包"
        echo "  3. 检查网络连接和端口配置"
        echo "  4. 确保文件权限正确"
        exit 1
    fi
}

# 运行主函数
main "$@"
