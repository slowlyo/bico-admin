#!/bin/bash

# =============================================================================
# Bico Admin Docker 构建脚本
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
PROJECT_NAME="bico-admin"
IMAGE_NAME="bico-admin"
VERSION=${1:-latest}

# 函数：打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}[$(date +'%Y-%m-%d %H:%M:%S')] ${message}${NC}"
}

# 函数：检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_message $RED "错误: $1 命令未找到，请先安装 $1"
        exit 1
    fi
}

# 函数：构建 Docker 镜像
build_image() {
    print_message $BLUE "开始构建 Docker 镜像..."
    
    # 检查 Dockerfile 是否存在
    if [ ! -f "Dockerfile" ]; then
        print_message $RED "错误: Dockerfile 不存在"
        exit 1
    fi
    
    # 构建镜像
    docker build \
        --tag ${IMAGE_NAME}:${VERSION} \
        --tag ${IMAGE_NAME}:latest \
        --build-arg BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
        --build-arg VERSION=${VERSION} \
        .
    
    print_message $GREEN "Docker 镜像构建完成: ${IMAGE_NAME}:${VERSION}"
}

# 函数：显示镜像信息
show_image_info() {
    print_message $BLUE "镜像信息:"
    docker images | grep ${IMAGE_NAME} | head -5
    
    print_message $BLUE "镜像大小:"
    docker image inspect ${IMAGE_NAME}:${VERSION} --format='{{.Size}}' | numfmt --to=iec
}

# 函数：运行容器测试
test_container() {
    print_message $BLUE "测试容器运行..."
    
    # 停止并删除已存在的测试容器
    docker stop ${PROJECT_NAME}-test 2>/dev/null || true
    docker rm ${PROJECT_NAME}-test 2>/dev/null || true
    
    # 运行测试容器
    docker run -d \
        --name ${PROJECT_NAME}-test \
        -p 8899:8899 \
        -e BICO_APP_ENVIRONMENT=development \
        -e BICO_APP_DEBUG=true \
        -e LOG_LEVEL=debug \
        ${IMAGE_NAME}:${VERSION}
    
    # 等待容器启动
    print_message $YELLOW "等待容器启动..."
    sleep 10
    
    # 检查容器状态
    if docker ps | grep -q ${PROJECT_NAME}-test; then
        print_message $GREEN "容器启动成功"
        print_message $BLUE "容器日志:"
        docker logs ${PROJECT_NAME}-test --tail 20
        
        # 健康检查
        print_message $BLUE "执行健康检查..."
        if curl -f http://localhost:8899/admin-api/auth/login &>/dev/null; then
            print_message $GREEN "健康检查通过"
        else
            print_message $YELLOW "健康检查失败，但容器正在运行"
        fi
    else
        print_message $RED "容器启动失败"
        docker logs ${PROJECT_NAME}-test
        exit 1
    fi
    
    # 清理测试容器
    print_message $BLUE "清理测试容器..."
    docker stop ${PROJECT_NAME}-test
    docker rm ${PROJECT_NAME}-test
}

# 函数：推送镜像
push_image() {
    local registry=$1
    if [ -n "$registry" ]; then
        print_message $BLUE "推送镜像到 ${registry}..."
        docker tag ${IMAGE_NAME}:${VERSION} ${registry}/${IMAGE_NAME}:${VERSION}
        docker push ${registry}/${IMAGE_NAME}:${VERSION}
        print_message $GREEN "镜像推送完成"
    fi
}

# 函数：显示帮助信息
show_help() {
    echo "用法: $0 [版本号] [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -t, --test     构建后运行测试"
    echo "  -p, --push     推送镜像到仓库"
    echo "  --registry     指定镜像仓库地址"
    echo ""
    echo "示例:"
    echo "  $0                    # 构建 latest 版本"
    echo "  $0 v1.0.0             # 构建 v1.0.0 版本"
    echo "  $0 v1.0.0 --test      # 构建并测试"
    echo "  $0 v1.0.0 --push --registry registry.example.com"
}

# 主函数
main() {
    print_message $GREEN "🚀 Bico Admin Docker 构建脚本"
    print_message $BLUE "版本: ${VERSION}"
    
    # 检查必要的命令
    check_command docker
    check_command curl
    
    # 解析参数
    local run_test=false
    local push_image_flag=false
    local registry=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -t|--test)
                run_test=true
                shift
                ;;
            -p|--push)
                push_image_flag=true
                shift
                ;;
            --registry)
                registry="$2"
                shift 2
                ;;
            *)
                if [[ $1 =~ ^v?[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                    VERSION=$1
                fi
                shift
                ;;
        esac
    done
    
    # 构建镜像
    build_image
    
    # 显示镜像信息
    show_image_info
    
    # 运行测试
    if [ "$run_test" = true ]; then
        test_container
    fi
    
    # 推送镜像
    if [ "$push_image_flag" = true ]; then
        push_image $registry
    fi
    
    print_message $GREEN "✅ 构建完成!"
    print_message $BLUE "运行命令: docker run -p 8899:8899 ${IMAGE_NAME}:${VERSION}"
}

# 执行主函数
main "$@"
