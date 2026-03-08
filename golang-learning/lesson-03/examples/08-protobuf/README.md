# 08 - Protocol Buffers (Protobuf)

## 说明

演示如何在 Gin 框架中使用 Protocol Buffers（Protobuf）进行数据序列化和反序列化。

## 知识点

- Protocol Buffers 简介和特点
- 定义 `.proto` 文件
- 使用 `protoc` 编译生成 Go 代码
- 在 Gin 中使用 Protobuf 处理请求和响应
- Protobuf 与 JSON 的对比

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
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### 3. 安装依赖

```bash
go get google.golang.org/protobuf/proto
```

## 生成 Protobuf 代码

如果修改了 `user.proto` 文件，需要重新生成 Go 代码：

```bash
protoc --go_out=. --go_opt=paths=source_relative user.proto
```

生成的文件位于 `pb/user.pb.go`。

## 运行方式

```bash
# 安装依赖
go mod tidy

# 运行服务器
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 端点

### 1. 获取单个用户（Protobuf 格式）

```bash
curl -X GET http://localhost:8080/api/proto/user/1 \
  -H "Accept: application/x-protobuf" \
  --output user.bin

# 查看二进制数据大小
ls -lh user.bin
```

### 2. 获取用户列表（Protobuf 格式）

```bash
curl -X GET http://localhost:8080/api/proto/users \
  -H "Accept: application/x-protobuf" \
  --output users.bin
```

### 3. 创建用户（Protobuf 格式）

首先需要准备一个 Protobuf 格式的请求文件。项目已包含 Python 脚本 `create_user.py` 用于生成请求文件。

**步骤 1：生成 Python Protobuf 代码**

```bash
# 安装 Python protobuf 编译器（如果还没有）
pip install protobuf

# 生成 Python 的 protobuf 代码
protoc --python_out=. user.proto
```

这会生成 `user_pb2.py` 文件。

**步骤 2：安装 Python 依赖**

```bash
pip install -r requirements.txt
```

**步骤 3：运行脚本生成请求文件**

```bash
# 使用默认值生成
python create_user.py

# 或指定自定义参数
python create_user.py "myuser" "myuser@example.com" 30
```

脚本会生成 `create_user.bin` 文件。

**步骤 4：发送请求**

```bash
curl -X POST http://localhost:8080/api/proto/user \
  -H "Content-Type: application/x-protobuf" \
  -H "Accept: application/x-protobuf" \
  --data-binary @create_user.bin \
  --output response.bin
```

#### 使用 Postman 发送 Protobuf 请求

在 Postman 中发送 Protobuf 请求需要以下步骤：

**步骤 1：准备 Protobuf 二进制文件**

首先使用 Python 脚本生成二进制文件（参考上面的 `create_user.py`），或者使用其他工具生成。

**步骤 2：在 Postman 中配置请求**

1. **设置请求方法和 URL**
   - Method: `POST`
   - URL: `http://localhost:8080/api/proto/user`

2. **设置 Headers**
   - `Content-Type`: `application/x-protobuf`
   - `Accept`: `application/x-protobuf`

3. **设置 Body**
   - 选择 `Body` 标签
   - 选择 `binary` 选项（不是 `raw`）
   - 点击 `Select File` 按钮
   - 选择之前生成的 `create_user.bin` 文件

4. **发送请求**
   - 点击 `Send` 按钮
   - 响应会以二进制格式返回，可以保存为文件查看

**注意事项：**
- Postman 默认会将二进制响应显示为十六进制，可以点击 `Save Response` → `Save to a file` 保存为 `.bin` 文件
- 如果需要查看响应内容，可以使用 `protoc` 解码（见下方说明）
- 确保服务器已启动（`go run main.go`）

**查看 Protobuf 响应：**

保存响应文件后，可以使用以下命令解码：

```bash
# 解码 CreateUserResponse
protoc --decode=user.CreateUserResponse user.proto < response.bin
```

### 4. 获取用户（JSON 格式，用于对比）

```bash
curl -X GET http://localhost:8080/api/proto/user/1/json
```

## Protobuf vs JSON 对比

### 性能对比

| 特性 | Protobuf | JSON |
|------|----------|------|
| 大小 | 小（二进制） | 大（文本） |
| 速度 | 快 | 慢 |
| 可读性 | 差（二进制） | 好（文本） |
| 类型安全 | 是 | 否 |
| 浏览器支持 | 需要解码 | 原生支持 |

### 使用场景

- **Protobuf**：微服务通信、高性能 API、数据存储
- **JSON**：Web API、前端交互、配置文件

## 文件结构

```
08-protobuf/
├── user.proto          # Protobuf 定义文件
├── pb/
│   └── user.pb.go     # 生成的 Go 代码
├── main.go            # 主程序
├── go.mod             # Go 依赖管理
├── create_user.py     # Python 脚本：生成 Protobuf 请求文件
├── requirements.txt   # Python 依赖
├── user_pb2.py        # 生成的 Python protobuf 代码（运行脚本后生成）
└── README.md          # 说明文档
```

## .proto 文件说明

### user.proto

定义了以下消息类型：

1. **User**：用户信息
   - 基本字段：id, username, email, age, active
   - 数组字段：tags（字符串数组）
   - 映射字段：metadata（键值对）

2. **UserList**：用户列表
   - users：用户数组
   - total：总数

3. **CreateUserRequest**：创建用户请求
   - username, email, age

4. **CreateUserResponse**：创建用户响应
   - user：创建的用户
   - success：是否成功
   - message：消息

### 字段编号规则

- **1-15**：常用字段（占用 1 字节）
- **16-2047**：不常用字段（占用 2 字节）
- **重要**：不要使用已删除的字段编号

## 版本演进

Protobuf 支持向后兼容的版本演进：

### ✅ 允许的操作

1. 添加新字段（使用新编号）
2. 删除可选字段
3. 重命名字段（编号不变）

### ❌ 不允许的操作

1. 修改字段类型
2. 重用已删除的字段编号

### 示例

```protobuf
// v1
message User {
  int64 id = 1;
  string name = 2;
}

// v2（兼容 v1）
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;  // 新增字段 ✅
  // int32 age = 4;   // 已删除，但编号不能重用 ❌
}
```

## 测试工具

### 使用 grpcurl（推荐）

```bash
# 安装 grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### 使用 Python 脚本

```python
import requests
import user_pb2

# 创建请求
req = user_pb2.CreateUserRequest()
req.username = "testuser"
req.email = "test@example.com"
req.age = 25

# 发送请求
response = requests.post(
    "http://localhost:8080/api/proto/user",
    data=req.SerializeToString(),
    headers={"Content-Type": "application/x-protobuf"}
)

# 解析响应
resp = user_pb2.CreateUserResponse()
resp.ParseFromString(response.content)
print(resp)
```

## 最佳实践

1. **字段编号**：常用字段使用 1-15
2. **命名规范**：使用驼峰命名（CamelCase）
3. **包管理**：使用 `option go_package` 指定包路径
4. **版本控制**：将 `.proto` 文件纳入版本控制
5. **文档注释**：使用注释说明字段用途

## 注意事项

1. 生成的 `.pb.go` 文件不要手动修改
2. 修改 `.proto` 文件后需要重新生成代码
3. Protobuf 是二进制格式，调试不如 JSON 方便
4. 浏览器不支持原生 Protobuf，需要 JavaScript 库解码
5. 确保 `protoc` 和 `protoc-gen-go` 版本兼容

## 常见问题

**Q: 如何查看 Protobuf 二进制数据？**

A: 可以使用 `protoc` 解码：

```bash
# 需要先定义 .proto 文件
protoc --decode=User user.proto < user.bin
```

**Q: Protobuf 和 JSON 可以混用吗？**

A: 可以。同一个 API 可以同时支持 Protobuf 和 JSON，根据 `Accept` 头返回不同格式。

**Q: 如何实现 Protobuf 和 JSON 的自动转换？**

A: Protobuf 生成的 Go 代码支持 JSON 序列化，可以使用 `jsonpb` 包进行转换。

**Q: Protobuf 文件很大怎么办？**

A: 可以：
1. 使用流式传输
2. 分页加载
3. 只传输必要字段（使用 FieldMask）

