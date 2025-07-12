# Config 模块测试

本目录包含 config 模块的完整单元测试和集成测试。

## 测试文件说明

### 测试文件
- `config_test.go` - 基础配置加载和验证测试
- `loader_test.go` - 配置加载器功能测试
- `switch_test.go` - 配置切换和环境变量覆盖测试
- `integration_test.go` - 集成测试，测试真实配置文件

### 测试数据
- `testdata/test_config.yml` - 测试用的完整配置文件
- `testdata/test_config_override.yml` - 测试用的覆盖配置文件
- `testdata/invalid_config.yml` - 测试用的无效配置文件

## 测试覆盖的功能

### 1. 配置加载 (config_test.go)
- ✅ 加载有效配置文件
- ✅ 处理无效配置文件
- ✅ 处理不存在的配置文件
- ✅ 环境变量覆盖配置
- ✅ 配置项类型验证（字符串、整数、布尔值、时间）

### 2. 配置加载器 (loader_test.go)
- ✅ 默认配置加载
- ✅ 指定配置文件加载
- ✅ 配置验证逻辑
- ✅ 数据库驱动验证
- ✅ 时间配置解析
- ✅ 环境变量格式验证

### 3. 配置切换 (switch_test.go)
- ✅ 不同配置文件之间切换
- ✅ 环境变量在不同配置中的覆盖
- ✅ 部分环境变量覆盖
- ✅ 复杂环境变量覆盖
- ✅ 配置切换时的验证

### 4. 集成测试 (integration_test.go)
- ✅ 真实配置文件测试
- ✅ 配置与环境变量集成
- ✅ 配置验证集成
- ✅ 数据库驱动切换

## 运行测试

### 运行所有配置测试
```bash
# 在项目根目录运行
make test-config

# 或者直接运行
go test ./tests/config/... -v
```

### 运行特定测试
```bash
# 运行基础配置测试
go test ./tests/config -run TestLoadValidConfig -v

# 运行配置切换测试
go test ./tests/config -run TestSwitchBetweenConfigs -v

# 运行集成测试
go test ./tests/config -run TestRealConfigFiles -v
```

### 查看测试覆盖率
```bash
go test ./tests/config/... -cover -v
```

## 测试环境要求

- Go 1.19+
- testify 测试框架
- 项目根目录下的配置文件：
  - `config/app.yml`
  - `config/app.dev.yml`

## 测试数据说明

测试使用独立的测试数据文件，不会影响项目的实际配置文件。所有测试都会在测试后清理环境变量，确保测试之间不会相互影响。

## 添加新测试

当添加新的配置功能时，请：

1. 在相应的测试文件中添加测试用例
2. 如需要新的测试数据，在 `testdata/` 目录下添加
3. 更新本 README 文档
4. 确保所有测试通过

## 常见问题

### 测试失败：找不到配置文件
确保在项目根目录运行测试，或者测试会自动切换到正确的目录。

### 环境变量影响测试
测试会自动清理环境变量，但如果手动设置了 `BICO_*` 环境变量，可能会影响测试结果。
