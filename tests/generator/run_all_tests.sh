#!/bin/bash

# 代码生成器测试套件运行脚本
# 作者: Bico Admin Team
# 描述: 运行完整的代码生成器测试套件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TEST_DIR="$PROJECT_ROOT/tests/generator"

echo -e "${BLUE}🚀 Bico Admin 代码生成器测试套件${NC}"
echo -e "${BLUE}===============================================${NC}"
echo -e "项目根目录: ${PROJECT_ROOT}"
echo -e "测试目录: ${TEST_DIR}"
echo ""

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到Go环境，请先安装Go${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go环境检查通过${NC}"
go version

# 进入项目根目录
cd "$PROJECT_ROOT"

# 检查依赖
echo -e "\n${YELLOW}📦 检查项目依赖...${NC}"
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ 错误: 未找到go.mod文件${NC}"
    exit 1
fi

# 下载依赖
echo -e "${YELLOW}📥 下载依赖包...${NC}"
go mod download
go mod tidy

# 创建测试输出目录
TEST_OUTPUT_DIR="$TEST_DIR/output"
mkdir -p "$TEST_OUTPUT_DIR"

echo -e "\n${PURPLE}🧪 开始运行测试套件...${NC}"
echo -e "${PURPLE}===============================================${NC}"

# 运行单元测试
echo -e "\n${CYAN}1️⃣  运行单元测试${NC}"
echo -e "${CYAN}-------------------${NC}"
cd "$PROJECT_ROOT"
go test -v ./tests/generator/unit/... -timeout=30m | tee "$TEST_OUTPUT_DIR/unit_tests.log"
UNIT_TEST_EXIT_CODE=${PIPESTATUS[0]}

if [ $UNIT_TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ 单元测试通过${NC}"
else
    echo -e "${RED}❌ 单元测试失败 (退出码: $UNIT_TEST_EXIT_CODE)${NC}"
fi

# 运行集成测试
echo -e "\n${CYAN}2️⃣  运行集成测试${NC}"
echo -e "${CYAN}-------------------${NC}"
go test -v ./tests/generator/integration/... -timeout=30m | tee "$TEST_OUTPUT_DIR/integration_tests.log"
INTEGRATION_TEST_EXIT_CODE=${PIPESTATUS[0]}

if [ $INTEGRATION_TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ 集成测试通过${NC}"
else
    echo -e "${RED}❌ 集成测试失败 (退出码: $INTEGRATION_TEST_EXIT_CODE)${NC}"
fi

# 运行编译验证测试
echo -e "\n${CYAN}3️⃣  运行编译验证测试${NC}"
echo -e "${CYAN}-------------------------${NC}"
go test -v ./tests/generator/compilation/... -timeout=30m | tee "$TEST_OUTPUT_DIR/compilation_tests.log"
COMPILATION_TEST_EXIT_CODE=${PIPESTATUS[0]}

if [ $COMPILATION_TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ 编译验证测试通过${NC}"
else
    echo -e "${RED}❌ 编译验证测试失败 (退出码: $COMPILATION_TEST_EXIT_CODE)${NC}"
fi

# 运行自定义测试套件
echo -e "\n${CYAN}4️⃣  运行自定义测试套件${NC}"
echo -e "${CYAN}-------------------------${NC}"
cd "$TEST_DIR"
if [ -f "run_tests.go" ]; then
    go run run_tests.go | tee "$TEST_OUTPUT_DIR/custom_tests.log"
    CUSTOM_TEST_EXIT_CODE=${PIPESTATUS[0]}
    
    if [ $CUSTOM_TEST_EXIT_CODE -eq 0 ]; then
        echo -e "${GREEN}✅ 自定义测试套件通过${NC}"
    else
        echo -e "${RED}❌ 自定义测试套件失败 (退出码: $CUSTOM_TEST_EXIT_CODE)${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  未找到自定义测试套件${NC}"
    CUSTOM_TEST_EXIT_CODE=0
fi

# 生成测试覆盖率报告
echo -e "\n${CYAN}5️⃣  生成测试覆盖率报告${NC}"
echo -e "${CYAN}---------------------------${NC}"
cd "$PROJECT_ROOT"
COVERAGE_FILE="$TEST_OUTPUT_DIR/coverage.out"
COVERAGE_HTML="$TEST_OUTPUT_DIR/coverage.html"

echo -e "${YELLOW}📊 生成覆盖率数据...${NC}"
go test -coverprofile="$COVERAGE_FILE" ./tests/generator/... -timeout=30m

if [ -f "$COVERAGE_FILE" ]; then
    echo -e "${YELLOW}📈 生成HTML覆盖率报告...${NC}"
    go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
    
    echo -e "${YELLOW}📋 覆盖率统计:${NC}"
    go tool cover -func="$COVERAGE_FILE" | tail -1
    
    echo -e "${GREEN}✅ 覆盖率报告已生成: $COVERAGE_HTML${NC}"
else
    echo -e "${YELLOW}⚠️  未生成覆盖率数据${NC}"
fi

# 运行性能基准测试
echo -e "\n${CYAN}6️⃣  运行性能基准测试${NC}"
echo -e "${CYAN}------------------------${NC}"
BENCHMARK_FILE="$TEST_OUTPUT_DIR/benchmark.txt"
echo -e "${YELLOW}⏱️  运行基准测试...${NC}"
go test -bench=. -benchmem ./tests/generator/... -timeout=30m | tee "$BENCHMARK_FILE"
BENCHMARK_EXIT_CODE=${PIPESTATUS[0]}

if [ $BENCHMARK_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✅ 基准测试完成${NC}"
else
    echo -e "${YELLOW}⚠️  基准测试可能有问题 (退出码: $BENCHMARK_EXIT_CODE)${NC}"
fi

# 检查生成的文件
echo -e "\n${CYAN}7️⃣  检查生成的测试文件${NC}"
echo -e "${CYAN}---------------------------${NC}"
GENERATED_FILES_COUNT=0
if [ -d "internal/shared/models" ]; then
    GENERATED_FILES_COUNT=$((GENERATED_FILES_COUNT + $(find internal/shared/models -name "*.go" -type f | wc -l)))
fi
if [ -d "internal/admin" ]; then
    GENERATED_FILES_COUNT=$((GENERATED_FILES_COUNT + $(find internal/admin -name "*_gen.go" -o -name "*test*.go" -type f | wc -l)))
fi

echo -e "${YELLOW}📁 发现 $GENERATED_FILES_COUNT 个生成的测试文件${NC}"

# 清理测试生成的文件（可选）
read -p "是否清理测试生成的文件? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🧹 清理测试生成的文件...${NC}"
    
    # 清理模型文件
    find internal/shared/models -name "*test*.go" -type f -delete 2>/dev/null || true
    
    # 清理其他测试生成的文件
    find internal/admin -name "*test*_gen.go" -type f -delete 2>/dev/null || true
    find internal/admin -name "*Test*.go" -type f -delete 2>/dev/null || true
    
    echo -e "${GREEN}✅ 清理完成${NC}"
else
    echo -e "${BLUE}ℹ️  保留测试生成的文件${NC}"
fi

# 汇总测试结果
echo -e "\n${PURPLE}📊 测试结果汇总${NC}"
echo -e "${PURPLE}===============================================${NC}"

TOTAL_FAILURES=0

echo -e "单元测试: $([ $UNIT_TEST_EXIT_CODE -eq 0 ] && echo -e "${GREEN}✅ 通过${NC}" || echo -e "${RED}❌ 失败${NC}")"
[ $UNIT_TEST_EXIT_CODE -ne 0 ] && TOTAL_FAILURES=$((TOTAL_FAILURES + 1))

echo -e "集成测试: $([ $INTEGRATION_TEST_EXIT_CODE -eq 0 ] && echo -e "${GREEN}✅ 通过${NC}" || echo -e "${RED}❌ 失败${NC}")"
[ $INTEGRATION_TEST_EXIT_CODE -ne 0 ] && TOTAL_FAILURES=$((TOTAL_FAILURES + 1))

echo -e "编译验证测试: $([ $COMPILATION_TEST_EXIT_CODE -eq 0 ] && echo -e "${GREEN}✅ 通过${NC}" || echo -e "${RED}❌ 失败${NC}")"
[ $COMPILATION_TEST_EXIT_CODE -ne 0 ] && TOTAL_FAILURES=$((TOTAL_FAILURES + 1))

echo -e "自定义测试套件: $([ $CUSTOM_TEST_EXIT_CODE -eq 0 ] && echo -e "${GREEN}✅ 通过${NC}" || echo -e "${RED}❌ 失败${NC}")"
[ $CUSTOM_TEST_EXIT_CODE -ne 0 ] && TOTAL_FAILURES=$((TOTAL_FAILURES + 1))

echo ""
echo -e "测试输出目录: ${BLUE}$TEST_OUTPUT_DIR${NC}"
echo -e "测试日志文件:"
echo -e "  - 单元测试: $TEST_OUTPUT_DIR/unit_tests.log"
echo -e "  - 集成测试: $TEST_OUTPUT_DIR/integration_tests.log"
echo -e "  - 编译验证: $TEST_OUTPUT_DIR/compilation_tests.log"
echo -e "  - 自定义测试: $TEST_OUTPUT_DIR/custom_tests.log"
echo -e "  - 基准测试: $TEST_OUTPUT_DIR/benchmark.txt"
[ -f "$COVERAGE_HTML" ] && echo -e "  - 覆盖率报告: $COVERAGE_HTML"

echo ""
if [ $TOTAL_FAILURES -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试都通过了！代码生成器工作正常。${NC}"
    exit 0
else
    echo -e "${RED}💥 有 $TOTAL_FAILURES 个测试套件失败，请检查日志文件。${NC}"
    exit 1
fi
