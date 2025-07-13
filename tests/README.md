# Flash OAuth2 测试指南

本目录包含 Flash OAuth2 服务器的完整测试套件。

## 📋 测试文件

| 文件                         | 说明            | 测试内容                          |
| ---------------------------- | --------------- | --------------------------------- |
| `basic_test.go`              | 基础测试        | 基础验证（无外部依赖）            |
| `e2e_oauth2_test.go`         | OAuth2 流程测试 | 完整授权流程、手机验证、JWKS 端点 |
| `e2e_api_test.go`            | API 端点测试    | 所有 REST API、错误处理、安全验证 |
| `e2e_app_management_test.go` | 应用管理测试    | 开发者注册、应用管理、密钥管理    |
| `config_test.go`             | 配置测试        | 配置管理和测试数据                |
| `environment_test.go`        | 环境测试        | 环境配置验证                      |
| `e2e_test_helper.go`         | 测试工具        | 测试辅助函数和工具                |
| `test_main.go`               | 测试入口        | 测试主入口和配置                  |
| `test_data.go`               | 测试数据        | 测试数据工厂                      |
| `test_utils.go`              | 测试工具        | 通用测试工具函数                  |

## 🚀 快速测试

### 基础测试（推荐开始）

```bash
# 无需数据库的快速测试
make test-e2e-quick
```

### 完整测试

```bash
# 运行所有测试（需要数据库）
make test-e2e-all
```

### 特定测试

```bash
# 运行特定测试
make test-e2e-specific TEST=TestCompleteOAuth2Flow
```

## 📊 测试命令

| 命令                     | 说明                                  |
| ------------------------ | ------------------------------------- |
| `make test-e2e-all`      | 完整测试套件（设置+测试+基准+覆盖率） |
| `make test-e2e-run`      | 运行 E2E 测试                         |
| `make test-e2e-quick`    | 快速测试（无环境设置）                |
| `make test-e2e-bench`    | 性能基准测试                          |
| `make test-e2e-coverage` | 生成覆盖率报告                        |

## 🔧 测试环境

测试环境自动配置：

- **测试数据库**: `oauth2_test`
- **Redis 数据库**: `15`
- **测试端口**: `8081`

## 📈 性能测试

```bash
# 运行性能测试
make test-e2e-bench

# 查看覆盖率
make test-e2e-coverage
make test-e2e-view-coverage
```

更多详细信息请参考主项目的 [README.md](../README.md)。
