# 07 - 配置管理

## 说明

演示如何使用 Viper 管理应用配置，支持配置文件、环境变量和默认值。

## 知识点

- Viper 配置管理库
- YAML 配置文件
- 环境变量读取
- 配置默认值
- 配置结构体映射
- 配置文件搜索路径

## 运行方式

```bash
go run main.go
```

## 配置优先级

Viper 配置的优先级从高到低：

1. 显式调用 `Set` 设置的值
2. 环境变量
3. 配置文件
4. 默认值

## 测试

### 1. 使用配置文件

```bash
# 默认读取当前目录的 config.yaml
go run main.go

# 访问配置信息
curl http://localhost:8080/config
curl http://localhost:8080/health
```

### 2. 使用环境变量

```bash
# 环境变量会覆盖配置文件
export APP_SERVER_PORT=9000
export APP_SERVER_MODE=release
go run main.go

# 服务器会在 9000 端口启动
curl http://localhost:9000/config
```

### 3. 修改配置文件

编辑 `config.yaml`：

```yaml
server:
  port: "3000"
  host: "127.0.0.1"
  mode: "release"
```

重新运行程序，服务器将在 3000 端口启动。

## 配置文件说明

### config.yaml

```yaml
server:
  port: "8080"        # 服务器端口
  host: "0.0.0.0"     # 服务器地址
  mode: "debug"       # Gin 模式：debug, release, test

database:
  host: "localhost"   # 数据库地址
  port: 3306          # 数据库端口
  username: "root"    # 数据库用户名
  password: "password" # 数据库密码
  dbname: "myapp"     # 数据库名称

jwt:
  secret: "your-secret-key"  # JWT 密钥
  expire: "24h"             # Token 过期时间
```

## 环境变量

使用环境变量时，Viper 会自动将 `.` 替换为 `_`，并添加前缀 `APP_`：

```bash
# 配置文件中的 server.port 对应环境变量 APP_SERVER_PORT
export APP_SERVER_PORT=9000

# 配置文件中的 database.host 对应环境变量 APP_DATABASE_HOST
export APP_DATABASE_HOST=192.168.1.100
```

## 生产环境最佳实践

### 1. 敏感信息使用环境变量

```bash
export APP_DATABASE_PASSWORD=secure_password
export APP_JWT_SECRET=very_secure_secret_key
```

### 2. 不同环境使用不同配置文件

```bash
# 开发环境
config.dev.yaml

# 测试环境
config.test.yaml

# 生产环境
config.prod.yaml
```

在代码中根据环境变量选择配置文件：

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "dev"
}
viper.SetConfigName("config." + env)
```

### 3. 配置验证

添加配置验证逻辑：

```go
func (c *Config) Validate() error {
    if c.Server.Port == "" {
        return errors.New("server port is required")
    }
    if c.JWT.Secret == "" {
        return errors.New("jwt secret is required")
    }
    return nil
}
```

## Viper 常用方法

```go
// 读取配置
viper.Get("server.port")
viper.GetString("server.host")
viper.GetInt("database.port")
viper.GetBool("server.debug")

// 设置配置
viper.Set("server.port", "9000")

// 检查配置是否存在
viper.IsSet("server.port")

// 设置默认值
viper.SetDefault("server.port", "8080")

// 监听配置文件变化
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
})
```

## 注意事项

1. 敏感信息（密码、密钥）不要提交到版本控制
2. 生产环境建议使用环境变量
3. 配置文件应该有清晰的注释
4. 为配置项设置合理的默认值

