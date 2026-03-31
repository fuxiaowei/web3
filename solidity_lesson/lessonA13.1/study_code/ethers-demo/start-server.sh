#!/bin/bash
# 启动本地服务器来运行 demo.html
# 这样可以避免 file:// 协议的问题

echo "正在启动本地服务器..."
echo "访问地址: http://localhost:8080/demo.html"
echo "按 Ctrl+C 停止服务器"
echo ""

# 检查是否安装了 Python
if command -v python3 &> /dev/null; then
    python3 -m http.server 8080
elif command -v python &> /dev/null; then
    python -m http.server 8080
elif command -v npx &> /dev/null; then
    npx http-server -p 8080
else
    echo "错误: 未找到 Python 或 npx"
    echo "请安装 Python 或 Node.js"
    exit 1
fi

