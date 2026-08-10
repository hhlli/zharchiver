#!/bin/bash
# ZHArchiver 启动脚本
# 用法：双击此脚本，或在终端执行 bash start.sh

cd "$(dirname "$0")/zharchiver" || exit 1
pkill -x zharchiver 2>/dev/null
echo "启动 ZHArchiver..."
./zharchiver &
sleep 1
echo "服务已启动：http://localhost:8080"
open "http://localhost:8080"
