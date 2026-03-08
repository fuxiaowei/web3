# 示例测试指南

## 快速验证所有示例

使用以下命令快速测试所有示例是否可以正常编译：

```bash
cd /Users/zhouxing/Documents/rcc/new-golang/lesson-03/examples

# 测试所有示例编译
for dir in 01-hello-world 02-routing 03-request-handling 04-middleware 05-jwt-auth 06-file-upload 07-config; do
    echo "Testing $dir..."
    cd $dir
    go build -o test_binary
    if [ $? -eq 0 ]; then
        echo "✅ $dir 编译成功"
        rm test_binary
    else
        echo "❌ $dir 编译失败"
    fi
    cd ..
done
```

## 单个示例测试

### 01 - Hello World

```bash
cd 01-hello-world
go run main.go &
PID=$!
sleep 2

# 测试
curl http://localhost:8080/hello
curl http://localhost:8080/ping

# 停止服务
kill $PID
```

### 02 - Routing

```bash
cd 02-routing
go run main.go &
PID=$!
sleep 2

# 测试
curl http://localhost:8080/
curl http://localhost:8080/users/123
curl "http://localhost:8080/search?keyword=golang&page=2"
curl http://localhost:8080/api/v1/users

# 停止服务
kill $PID
```

### 03 - Request Handling

```bash
cd 03-request-handling
go run main.go &
PID=$!
sleep 2

# 测试 JSON 绑定
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'

# 测试查询参数
curl "http://localhost:8080/api/products?page=1&size=20"

# 停止服务
kill $PID
```

### 04 - Middleware

```bash
cd 04-middleware
go run main.go &
PID=$!
sleep 2

# 测试公开路由
curl http://localhost:8080/public

# 测试需要认证的路由（无 token）
curl http://localhost:8080/api/users

# 测试需要认证的路由（有效 token）
curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/users

# 停止服务
kill $PID
```

### 05 - JWT Auth

```bash
cd 05-jwt-auth
go run main.go &
PID=$!
sleep 2

# 登录获取 token
RESPONSE=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $RESPONSE | jq -r '.token')

# 使用 token 访问受保护的接口
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/protected

# 停止服务
kill $PID
```

### 06 - File Upload

```bash
cd 06-file-upload
go run main.go &
PID=$!
sleep 2

# 创建测试文件
echo "Test file content" > test.txt

# 上传文件
curl -X POST http://localhost:8080/api/upload \
  -F "file=@test.txt"

# 下载文件
curl http://localhost:8080/api/download/test.txt

# 清理
rm test.txt

# 停止服务
kill $PID
```

### 07 - Config

```bash
cd 07-config
go run main.go &
PID=$!
sleep 2

# 测试
curl http://localhost:8080/config
curl http://localhost:8080/health

# 停止服务
kill $PID
```

## 一键测试脚本

创建一个测试脚本 `test_all.sh`：

```bash
#!/bin/bash

cd "$(dirname "$0")"

echo "开始测试所有示例..."

# 01 - Hello World
echo -e "\n=== 测试 01-hello-world ==="
cd 01-hello-world
go run main.go &
PID=$!
sleep 2
curl -s http://localhost:8080/hello
kill $PID
cd ..

# 02 - Routing
echo -e "\n=== 测试 02-routing ==="
cd 02-routing
go run main.go &
PID=$!
sleep 2
curl -s http://localhost:8080/api/v1/users
kill $PID
cd ..

# 03 - Request Handling
echo -e "\n=== 测试 03-request-handling ==="
cd 03-request-handling
go run main.go &
PID=$!
sleep 2
curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","age":25}'
kill $PID
cd ..

# 04 - Middleware
echo -e "\n=== 测试 04-middleware ==="
cd 04-middleware
go run main.go &
PID=$!
sleep 2
curl -s http://localhost:8080/public
kill $PID
cd ..

# 05 - JWT Auth
echo -e "\n=== 测试 05-jwt-auth ==="
cd 05-jwt-auth
go run main.go &
PID=$!
sleep 2
curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
kill $PID
cd ..

# 06 - File Upload
echo -e "\n=== 测试 06-file-upload ==="
cd 06-file-upload
go run main.go &
PID=$!
sleep 2
echo "test" > test.txt
curl -s -X POST http://localhost:8080/api/upload -F "file=@test.txt"
rm test.txt
kill $PID
cd ..

# 07 - Config
echo -e "\n=== 测试 07-config ==="
cd 07-config
go run main.go &
PID=$!
sleep 2
curl -s http://localhost:8080/health
kill $PID
cd ..

echo -e "\n所有测试完成！"
```

使用方法：

```bash
chmod +x test_all.sh
./test_all.sh
```

## 常见问题

### Q: 端口被占用怎么办？

```bash
# 查找占用 8080 端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>
```

### Q: 依赖下载失败？

```bash
# 设置 Go 代理
export GOPROXY=https://goproxy.cn,direct

# 重新下载依赖
go mod download
```

### Q: 编译失败？

```bash
# 清理并重新下载依赖
go clean -modcache
go mod tidy
```

