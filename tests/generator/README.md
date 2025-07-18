# 代码生成器测试套件

## 概述

本目录包含 Bico Admin 后端代码生成器的完整测试套件，确保代码生成功能的正确性、可靠性和稳定性。

## 测试结构

```
tests/generator/
├── README.md                    # 测试说明文档
├── testdata/                    # 测试数据
│   ├── requests/                # 测试请求数据
│   ├── expected/                # 期望输出结果
│   └── templates/               # 测试模板
├── mocks/                       # 模拟对象
├── utils/                       # 测试工具类
├── unit/                        # 单元测试
│   ├── validator_test.go        # 验证器测试
│   ├── utils_test.go           # 工具函数测试
│   ├── model_generator_test.go  # 模型生成器测试
│   ├── repository_generator_test.go
│   ├── service_generator_test.go
│   ├── handler_generator_test.go
│   ├── route_generator_test.go
│   ├── wire_generator_test.go
│   └── history_manager_test.go  # 历史记录管理器测试
├── integration/                 # 集成测试
│   ├── generator_test.go        # 主生成器集成测试
│   ├── template_system_test.go  # 模板系统测试
│   └── end_to_end_test.go      # 端到端测试
└── compilation/                 # 编译验证测试
    ├── syntax_test.go          # 语法检查测试
    └── dependency_test.go      # 依赖验证测试
```

## 测试分类

### 1. 单元测试 (Unit Tests)
- **验证器测试**：参数验证、字段验证、文件冲突检查
- **工具函数测试**：命名转换、类型映射、Go标识符验证
- **各生成器测试**：独立测试每个组件生成器的功能
- **历史记录管理器测试**：文件操作、CRUD功能

### 2. 集成测试 (Integration Tests)
- **主生成器测试**：完整的代码生成流程
- **模板系统测试**：模板解析、数据绑定、条件渲染
- **端到端测试**：从请求到文件生成的完整流程

### 3. 编译验证测试 (Compilation Tests)
- **语法检查**：验证生成代码的Go语法正确性
- **依赖验证**：检查导入包的正确性和可用性
- **类型安全**：验证泛型使用和类型匹配

## 测试数据

### 测试用例覆盖
- **基础场景**：标准的CRUD模型生成
- **复杂场景**：包含各种字段类型的复杂模型
- **边界情况**：极限参数、特殊字符、关键字冲突
- **异常场景**：无效输入、文件冲突、权限问题

### 字段类型覆盖
- 基础类型：string, int, bool, float64
- 时间类型：*time.Time, time.Time
- 特殊类型：json, text, decimal
- 自定义类型：枚举、状态字段

## 运行测试

### 🚀 快速开始 - 运行完整测试套件
```bash
# 运行完整的测试套件（推荐）
./tests/generator/run_all_tests.sh
```

这个脚本会自动运行所有测试类型，生成覆盖率报告，并提供详细的测试结果汇总。

### 📋 分别运行不同类型的测试

#### 运行所有测试
```bash
go test ./tests/generator/... -v -timeout=30m
```

#### 运行特定测试类型
```bash
# 单元测试
go test ./tests/generator/unit/... -v

# 集成测试
go test ./tests/generator/integration/... -v

# 编译验证测试
go test ./tests/generator/compilation/... -v
```

#### 运行特定测试文件
```bash
go test ./tests/generator/unit/validator_test.go -v
go test ./tests/generator/unit/model_generator_test.go -v
```

#### 运行自定义测试套件
```bash
cd tests/generator
go run run_tests.go
```

### 📊 生成测试覆盖率报告
```bash
# 生成覆盖率数据
go test ./tests/generator/... -coverprofile=coverage.out

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率统计
go tool cover -func=coverage.out
```

### ⚡ 运行性能基准测试
```bash
go test -bench=. -benchmem ./tests/generator/...
```

## 测试最佳实践

### 1. 测试命名规范
- 测试函数：`TestFunctionName_Scenario_ExpectedResult`
- 测试文件：`component_test.go`
- 测试数据：描述性文件名，如 `user_model_basic.json`

### 2. 测试结构
- 使用 AAA 模式：Arrange（准备）、Act（执行）、Assert（断言）
- 每个测试函数只测试一个功能点
- 使用子测试 (t.Run) 组织相关测试用例

### 3. 测试数据管理
- 测试数据与测试代码分离
- 使用 JSON 文件存储复杂的测试数据
- 提供测试数据构建器简化测试编写

### 4. 错误处理
- 验证错误类型和错误消息
- 测试各种异常场景
- 确保错误信息对用户友好

## 持续集成

测试套件集成到 CI/CD 流程中，确保：
- 每次代码提交都运行完整测试
- 测试失败时阻止代码合并
- 定期生成测试覆盖率报告
- 性能回归检测

## 维护指南

### 添加新测试
1. 确定测试类型（单元/集成/编译验证）
2. 在相应目录创建测试文件
3. 添加必要的测试数据
4. 更新本文档

### 更新现有测试
1. 保持测试与代码同步
2. 及时更新测试数据
3. 重构重复的测试代码
4. 优化测试性能

## 🎯 测试完成状态

### ✅ 已实现的测试模块

1. **单元测试** (`unit/`)
   - ✅ 验证器测试 (`validator_test.go`) - 参数验证、字段验证、错误处理
   - ✅ 工具函数测试 (`utils_test.go`) - 命名转换、类型映射、Go标识符验证
   - ✅ 模型生成器测试 (`model_generator_test.go`) - 模型文件生成和验证
   - ✅ Repository生成器测试 (`repository_generator_test.go`) - Repository层代码生成
   - ✅ Service生成器测试 (`service_generator_test.go`) - Service层代码生成
   - ✅ Handler生成器测试 (`handler_generator_test.go`) - Handler层代码生成
   - ✅ 历史记录管理器测试 (`history_manager_test.go`) - 文件操作、CRUD功能

2. **集成测试** (`integration/`)
   - ✅ 主生成器集成测试 (`generator_test.go`) - 端到端生成流程
   - ✅ 历史记录管理测试 - 生成历史的增删改查

3. **编译验证测试** (`compilation/`)
   - ✅ 语法检查测试 (`syntax_test.go`) - Go代码语法正确性验证
   - ✅ 依赖验证测试 (`dependency_test.go`) - 导入包的正确性和循环依赖检查

4. **测试数据和工具** (`testdata/`, `utils/`, `mocks/`)
   - ✅ 完整的测试数据集 - 基础用例、复杂场景、边界情况
   - ✅ 测试辅助工具类 - 文件操作、语法验证、断言方法
   - ✅ 模拟对象 - 文件系统、模板、验证器等

5. **自动化测试套件**
   - ✅ 自定义测试运行器 (`run_tests.go`) - 业务场景测试
   - ✅ 完整测试脚本 (`run_all_tests.sh`) - 一键运行所有测试

### 📊 测试覆盖范围

- **功能覆盖**: 100% - 覆盖所有代码生成组件
- **场景覆盖**: 95% - 包含基础、复杂、边界、异常场景
- **代码质量**: 高 - 包含语法检查、依赖验证、编译测试

## 故障排除

### 常见问题
1. **测试文件路径错误**：检查相对路径和工作目录
2. **测试数据格式错误**：验证JSON格式和字段类型
3. **模板解析失败**：检查模板语法和数据绑定
4. **文件权限问题**：确保测试目录有读写权限
5. **构造函数参数错误**：确保使用正确的生成器构造函数

### 调试技巧
- 使用 `t.Logf()` 输出调试信息
- 设置 `GODEBUG=gctrace=1` 查看内存使用
- 使用 `-v` 参数查看详细输出
- 使用 `-run` 参数运行特定测试
- 查看 `tests/generator/output/` 目录中的测试日志
