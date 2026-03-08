# 09 - gRPC

## 说明

演示如何使用 gRPC 构建高性能的 RPC 服务，包括服务端和客户端的实现，以及各种流式传输模式。

## 知识点

- gRPC 简介和特点
- 定义 gRPC 服务（`.proto` 文件）
- 使用 `protoc` 生成 gRPC 代码
- 实现 gRPC 服务端
- 实现 gRPC 客户端
- 四种 RPC 模式：
  - 一元 RPC（Unary RPC）
  - 服务端流（Server Streaming）
  - 客户端流（Client Streaming）
  - 双向流（Bidirectional Streaming）
- 错误处理和状态码

## 前置要求

### 1. 安装 protoc 编译器

```bash
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

# 或从源码编译
# https://github.com/protocolbuffers/protobuf/releases
```

### 2. 安装 Go 插件

```bash
# 安装 protobuf 编译器插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 gRPC 插件
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. 确保 PATH 包含 Go bin 目录

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 生成 gRPC 代码

如果修改了 `user.proto` 文件，需要重新生成 Go 代码：

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto
```

生成的文件位于 `pb/` 目录：
- `pb/user.pb.go` - Protobuf 消息定义
- `pb/user_grpc.pb.go` - gRPC 服务定义

## 运行方式

### 1. 启动服务端

```bash
# 安装依赖
go mod tidy

# 运行服务端
go run main.go server.go -mode=server -addr=:50051
```

或者：

```bash
go run . -mode=server
```

服务端将在 `:50051` 启动。

### 2. 运行客户端

在另一个终端：

```bash
# 运行客户端演示
go run main.go client.go -mode=client -addr=localhost:50051
```

或者：

```bash
go run . -mode=client
```

## gRPC 服务方法

### 1. 一元 RPC（Unary RPC）

#### GetUser - 获取单个用户

```go
req := &pb.GetUserRequest{Id: 1}
user, err := client.GetUser(ctx, req)
```

#### CreateUser - 创建用户

```go
req := &pb.CreateUserRequest{
    Username: "alice",
    Email:    "alice@example.com",
    Age:      25,
}
resp, err := client.CreateUser(ctx, req)
```

#### ListUsers - 获取用户列表

```go
req := &pb.ListUsersRequest{
    Page:     1,
    PageSize: 10,
}
resp, err := client.ListUsers(ctx, req)
```

#### UpdateUser - 更新用户

```go
req := &pb.UpdateUserRequest{
    Id:       1,
    Username: "alice-updated",
    Email:    "alice@example.com",
}
resp, err := client.UpdateUser(ctx, req)
```

#### DeleteUser - 删除用户

```go
req := &pb.DeleteUserRequest{Id: 1}
resp, err := client.DeleteUser(ctx, req)
```

### 2. 服务端流（Server Streaming）

#### StreamUsers - 流式获取用户

服务端逐个发送用户，客户端接收流：

```go
req := &pb.StreamUsersRequest{
    Limit:      10,
    IntervalMs: 500,
}
stream, err := client.StreamUsers(ctx, req)
for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    // 处理用户
}
```

**使用场景**：
- 实时数据推送
- 大文件分块传输
- 服务器端推送通知

### 3. 客户端流（Client Streaming）

#### BatchCreateUsers - 批量创建用户

客户端逐个发送请求，服务端接收流后返回单个响应：

```go
stream, err := client.BatchCreateUsers(ctx)
for _, user := range users {
    stream.Send(user)
}
resp, err := stream.CloseAndRecv()
```

**使用场景**：
- 批量上传
- 数据采集
- 日志收集

### 4. 双向流（Bidirectional Streaming）

#### ChatUsers - 聊天式交互

客户端和服务端都可以随时发送消息：

```go
stream, err := client.ChatUsers(ctx)
// 发送消息
stream.Send(&pb.ChatMessage{...})
// 接收消息
msg, err := stream.Recv()
```

**使用场景**：
- 实时聊天
- 游戏服务器
- 协作编辑

## 使用 grpcurl 测试

`grpcurl` 是一个命令行工具，类似于 `curl`，用于测试 gRPC 服务。

### 安装 grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### 列出服务

```bash
grpcurl -plaintext localhost:50051 list
```

### 列出方法

```bash
grpcurl -plaintext localhost:50051 list user.UserService
```

### 调用方法

```bash
# 获取用户
grpcurl -plaintext -d '{"id": 1}' \
  localhost:50051 user.UserService/GetUser

# 创建用户
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "age": 25
}' localhost:50051 user.UserService/CreateUser

# 获取用户列表
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 user.UserService/ListUsers
```

## 使用 grpcui 测试（Web UI）

`grpcui` 提供了一个 Web 界面来测试 gRPC 服务。

### 安装 grpcui

```bash
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

### 启动 grpcui

```bash
grpcui -plaintext localhost:50051
```

然后在浏览器中打开 `http://localhost:8080`，可以看到一个 Web 界面，可以直接调用 gRPC 方法。

## gRPC vs REST API

| 特性 | gRPC | REST API |
|------|------|----------|
| 协议 | HTTP/2 | HTTP/1.1 |
| 数据格式 | Protobuf（二进制） | JSON（文本） |
| 性能 | 高 | 中 |
| 流式传输 | 支持 | 有限支持（SSE/WebSocket） |
| 浏览器支持 | 需要 gRPC-Web | 原生支持 |
| 代码生成 | 自动生成 | 手动编写 |
| 类型安全 | 强类型 | 弱类型 |

## 错误处理

gRPC 使用状态码来表示错误：

```go
import "google.golang.org/grpc/status"
import "google.golang.org/grpc/codes"

// 返回错误
return nil, status.Errorf(codes.NotFound, "User not found")

// 检查错误
st, ok := status.FromError(err)
if ok {
    switch st.Code() {
    case codes.NotFound:
        // 处理未找到错误
    case codes.InvalidArgument:
        // 处理参数错误
    }
}
```

### 常用状态码

- `OK` - 成功
- `InvalidArgument` - 参数无效
- `NotFound` - 资源未找到
- `AlreadyExists` - 资源已存在
- `PermissionDenied` - 权限不足
- `Internal` - 服务器内部错误
- `Unavailable` - 服务不可用

## 文件结构

```
09-grpc/
├── user.proto           # gRPC 服务定义
├── pb/
│   ├── user.pb.go       # 生成的 Protobuf 代码
│   └── user_grpc.pb.go  # 生成的 gRPC 代码
├── server.go            # 服务端实现
├── client.go            # 客户端实现
├── main.go              # 主程序入口
├── go.mod               # 依赖管理
└── README.md            # 说明文档
```

## 安全

### TLS/SSL 加密

生产环境应该使用 TLS：

```go
// 服务端
creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
grpcServer := grpc.NewServer(grpc.Creds(creds))

// 客户端
creds, err := credentials.NewClientTLSFromFile("ca.crt", "")
conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
```

### 认证

可以使用 Token 认证：

```go
// 客户端
token := "your-auth-token"
conn, err := grpc.NewClient(
    addr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithPerRPCCredentials(&authCreds{token: token}),
)
```

## 性能优化

1. **连接复用**：客户端应该复用连接
2. **流式传输**：大数据使用流式传输
3. **压缩**：启用压缩减少网络传输
4. **超时设置**：合理设置超时时间
5. **负载均衡**：使用 gRPC 负载均衡

## 最佳实践

1. **版本控制**：将 `.proto` 文件纳入版本控制
2. **向后兼容**：遵循 Protobuf 版本演进规则
3. **错误处理**：使用标准状态码
4. **超时设置**：为所有 RPC 调用设置超时
5. **日志记录**：记录请求和响应（注意敏感信息）
6. **监控**：监控 RPC 调用指标
7. **文档**：为 `.proto` 文件添加注释

## 常见问题

**Q: 如何调试 gRPC 调用？**

A: 可以使用：
- `grpcurl` 命令行工具
- `grpcui` Web 界面
- 日志记录中间件
- 拦截器（Interceptor）

**Q: gRPC 支持哪些语言？**

A: gRPC 支持多种语言，包括 Go、Java、Python、C++、C#、Node.js、Ruby 等。

**Q: 如何在浏览器中使用 gRPC？**

A: 使用 gRPC-Web，需要代理服务器将 gRPC-Web 转换为 gRPC。

**Q: gRPC 和 REST 可以共存吗？**

A: 可以。可以在同一个服务中同时提供 gRPC 和 REST API。

**Q: 如何处理流式传输的错误？**

A: 在流式传输中，错误可以通过 `Recv()` 或 `Send()` 返回，需要检查 `io.EOF` 来判断流是否结束。

## 下一步

- 学习 gRPC 拦截器（Interceptor）
- 学习 gRPC 中间件
- 学习 gRPC 网关（将 gRPC 转换为 REST）
- 学习服务发现和负载均衡
- 学习监控和追踪

