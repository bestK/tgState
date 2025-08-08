#!/bin/bash

# TGState 启动脚本
# 用于启动 Telegram 文件存储服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 默认配置
DEFAULT_PORT="8088"
DEFAULT_MODE="p"
BINARY_NAME="tgstate"
PID_FILE="tgstate.pid"

# 显示帮助信息
show_help() {
    echo "TGState 启动脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -p, --port PORT        设置Web端口 (默认: $DEFAULT_PORT)"
    echo "  -m, --mode MODE        运行模式 (p=完整模式, m=最小模式) (默认: $DEFAULT_MODE)"
    echo "  -d, --daemon           后台运行"
    echo "  -s, --stop             停止服务"
    echo "  -r, --restart          重启服务"
    echo "  -t, --status           查看服务状态"
    echo "  -l, --logs             查看日志"
    echo "  -h, --help             显示此帮助信息"
    echo ""
    echo "环境变量 (可在 .env 文件中设置):"
    echo "  token      - Telegram Bot Token"
    echo "  target     - 目标频道名称或ID"
    echo "  pass       - 访问密码"
    echo "  apiPass    - API访问密码"
    echo "  url        - 基础URL"
    echo "  proxyUrl   - 代理URL"
    echo "  exts       - 允许的文件扩展名"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    
    # 检查网络工具
    if ! command -v netstat &> /dev/null && ! command -v ss &> /dev/null; then
        log_warn "未找到 netstat 或 ss 命令，端口检查功能可能受限"
    fi
    
    # 检查 .env 文件
    if [ ! -f ".env" ]; then
        log_warn ".env 文件不存在，将使用命令行参数或默认值"
    fi
    
    log_info "依赖检查完成"
}

# 获取端口占用进程信息
get_port_process() {
    local port=$1
    
    if command -v lsof &> /dev/null; then
        lsof -ti:$port 2>/dev/null
    elif command -v ss &> /dev/null; then
        ss -tulpn | grep ":$port " | awk -F',' '{print $2}' | awk -F'=' '{print $2}' | head -1
    elif command -v netstat &> /dev/null; then
        netstat -tulpn 2>/dev/null | grep ":$port " | awk '{print $7}' | cut -d'/' -f1 | head -1
    else
        echo ""
    fi
}

# 检查端口是否被占用
is_port_occupied() {
    local port=$1
    
    if command -v ss &> /dev/null; then
        ss -tuln | grep -q ":$port "
    elif command -v netstat &> /dev/null; then
        netstat -tuln 2>/dev/null | grep -q ":$port "
    else
        # 尝试连接端口来检查
        (echo >/dev/tcp/localhost/$port) &>/dev/null
    fi
}

# 构建项目
build_project() {
    log_info "构建项目..."
    
    if [ ! -f "go.mod" ]; then
        log_error "go.mod 文件不存在，请确保在项目根目录运行此脚本"
        exit 1
    fi
    
    # 下载依赖
    go mod tidy
    
    # 构建二进制文件
    go build -o $BINARY_NAME .
    
    if [ $? -eq 0 ]; then
        log_info "构建成功"
    else
        log_error "构建失败"
        exit 1
    fi
}

# 检查服务状态
check_status() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat $PID_FILE)
        if ps -p $PID > /dev/null 2>&1; then
            return 0  # 服务正在运行
        else
            rm -f $PID_FILE
            return 1  # PID文件存在但进程不存在
        fi
    else
        return 1  # PID文件不存在
    fi
}

# 停止服务
stop_service() {
    log_info "停止服务..."
    
    if check_status; then
        PID=$(cat $PID_FILE)
        kill $PID
        
        # 等待进程结束
        for i in {1..10}; do
            if ! ps -p $PID > /dev/null 2>&1; then
                break
            fi
            sleep 1
        done
        
        # 如果进程仍在运行，强制杀死
        if ps -p $PID > /dev/null 2>&1; then
            log_warn "强制停止服务..."
            kill -9 $PID
        fi
        
        rm -f $PID_FILE
        log_info "服务已停止"
    else
        log_warn "服务未运行"
    fi
}

# 启动服务
start_service() {
    local port=$1
    local mode=$2
    local daemon=$3
    
    log_info "启动服务..."
    
    # 检查服务是否已在运行
    if check_status; then
        log_error "服务已在运行 (PID: $(cat $PID_FILE))"
        exit 1
    fi
    
    # 检查端口是否被占用
    if is_port_occupied $port; then
        log_warn "端口 $port 已被占用"
        
        # 查找占用端口的进程
        PROCESS_PID=$(get_port_process $port)
        
        if [ ! -z "$PROCESS_PID" ] && [ "$PROCESS_PID" != "0" ]; then
            # 获取进程详细信息
            if command -v ps &> /dev/null; then
                PROCESS_DETAILS=$(ps -p $PROCESS_PID -o pid,ppid,user,cmd --no-headers 2>/dev/null)
                if [ ! -z "$PROCESS_DETAILS" ]; then
                    log_info "占用端口 $port 的进程信息:"
                    echo "  PID   PPID  USER     COMMAND"
                    echo "  $PROCESS_DETAILS"
                else
                    log_info "占用端口 $port 的进程 PID: $PROCESS_PID"
                fi
            else
                log_info "占用端口 $port 的进程 PID: $PROCESS_PID"
            fi
            
            echo ""
            echo -n "是否强制结束占用端口 $port 的进程并继续启动? [y/N]: "
            read -r response
            case "$response" in
                [yY][eE][sS]|[yY])
                    log_info "正在结束占用端口的进程 (PID: $PROCESS_PID)..."
                    
                    # 首先尝试优雅地结束进程
                    if kill $PROCESS_PID 2>/dev/null; then
                        log_debug "发送 TERM 信号到进程 $PROCESS_PID"
                        
                        # 等待进程结束
                        for i in {1..8}; do
                            if ! is_port_occupied $port; then
                                log_info "端口 $port 已释放"
                                break
                            fi
                            sleep 1
                            log_debug "等待进程结束... ($i/8)"
                        done
                        
                        # 如果端口仍被占用，强制杀死进程
                        if is_port_occupied $port; then
                            log_warn "正在强制结束进程 (PID: $PROCESS_PID)..."
                            if kill -9 $PROCESS_PID 2>/dev/null; then
                                sleep 2
                                
                                if is_port_occupied $port; then
                                    log_error "无法释放端口 $port，进程可能无法被终止"
                                    exit 1
                                else
                                    log_info "端口 $port 已强制释放"
                                fi
                            else
                                log_error "无法强制结束进程 $PROCESS_PID"
                                exit 1
                            fi
                        fi
                    else
                        log_error "无法发送信号到进程 $PROCESS_PID，可能权限不足"
                        exit 1
                    fi
                    ;;
                *)
                    log_info "用户取消操作"
                    exit 1
                    ;;
            esac
        else
            # 无法获取进程信息，但端口确实被占用
            log_warn "无法获取占用端口 $port 的进程信息"
            
            # 显示端口占用情况
            if command -v ss &> /dev/null; then
                PORT_INFO=$(ss -tulpn | grep ":$port ")
                if [ ! -z "$PORT_INFO" ]; then
                    log_info "端口占用详情:"
                    echo "$PORT_INFO"
                fi
            elif command -v netstat &> /dev/null; then
                PORT_INFO=$(netstat -tulpn 2>/dev/null | grep ":$port ")
                if [ ! -z "$PORT_INFO" ]; then
                    log_info "端口占用详情:"
                    echo "$PORT_INFO"
                fi
            fi
            
            echo ""
            echo -n "端口 $port 已被占用但无法获取进程信息，是否继续启动? (可能会失败) [y/N]: "
            read -r response
            case "$response" in
                [yY][eE][sS]|[yY])
                    log_warn "用户选择继续启动，启动可能会失败"
                    ;;
                *)
                    log_info "用户取消操作"
                    exit 1
                    ;;
            esac
        fi
    fi
    
    # 构建项目
    build_project
    
    # 准备启动参数
    ARGS="-port $port -mode $mode"
    
    if [ "$daemon" = true ]; then
        # 后台运行
        log_info "以守护进程模式启动服务 (端口: $port, 模式: $mode)"
        nohup ./$BINARY_NAME $ARGS > tgstate.log 2>&1 &
        echo $! > $PID_FILE
        log_info "服务已启动 (PID: $!)"
        log_info "日志文件: tgstate.log"
    else
        # 前台运行
        log_info "启动服务 (端口: $port, 模式: $mode)"
        ./$BINARY_NAME $ARGS
    fi
}

# 查看服务状态
show_status() {
    if check_status; then
        PID=$(cat $PID_FILE)
        log_info "服务正在运行 (PID: $PID)"
        
        # 显示端口信息
        if command -v netstat &> /dev/null; then
            PORTS=$(netstat -tuln 2>/dev/null | grep $PID 2>/dev/null | awk '{print $4}' | cut -d: -f2 | sort -u)
            if [ ! -z "$PORTS" ]; then
                log_info "监听端口: $PORTS"
            fi
        fi
        
        # 显示内存使用情况
        if command -v ps &> /dev/null; then
            MEMORY=$(ps -p $PID -o rss= 2>/dev/null | awk '{print int($1/1024)"MB"}')
            if [ ! -z "$MEMORY" ]; then
                log_info "内存使用: $MEMORY"
            fi
        fi
    else
        log_warn "服务未运行"
    fi
}

# 查看日志
show_logs() {
    if [ -f "tgstate.log" ]; then
        tail -f tgstate.log
    else
        log_error "日志文件不存在"
        exit 1
    fi
}

# 重启服务
restart_service() {
    local port=$1
    local mode=$2
    
    log_info "重启服务..."
    stop_service
    sleep 2
    start_service $port $mode true
}

# 主函数
main() {
    local port=$DEFAULT_PORT
    local mode=$DEFAULT_MODE
    local daemon=false
    local action="start"
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--port)
                port="$2"
                shift 2
                ;;
            -m|--mode)
                mode="$2"
                shift 2
                ;;
            -d|--daemon)
                daemon=true
                shift
                ;;
            -s|--stop)
                action="stop"
                shift
                ;;
            -r|--restart)
                action="restart"
                shift
                ;;
            -t|--status)
                action="status"
                shift
                ;;
            -l|--logs)
                action="logs"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 检查依赖
    check_dependencies
    
    # 执行相应操作
    case $action in
        start)
            start_service $port $mode $daemon
            ;;
        stop)
            stop_service
            ;;
        restart)
            restart_service $port $mode
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs
            ;;
        *)
            log_error "未知操作: $action"
            exit 1
            ;;
    esac
}

# 脚本入口
main "$@"