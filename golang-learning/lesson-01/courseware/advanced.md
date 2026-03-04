# 第1期 - 进阶特性部分课件

## 课程目标
- 掌握Go语言的接口和多态
- 理解并发编程的核心概念
- 学会使用Goroutine和Channel
- 掌握并发安全技术
- 完成并发编程实践

---

## 第一部分：接口和多态

### 1. 接口的定义和实现

在Go语言中，**接口（interface）**用来描述一组方法签名，任何类型只要实现了接口中声明的全部方法，就被视为满足该接口。Go不要求显式声明“implements”，这被称为**隐式实现**：类型的结构和功能自然地决定了它属于哪个接口。这样的设计让代码更松耦合，也便于我们在不改动已有类型的前提下，为它们赋予新的行为。

由于接口是基于行为的抽象，多种类型可以同时实现同一个接口。比如矩形、圆形都可以实现 `Shape` 接口，只要提供 `Area()` 和 `Perimeter()` 两个方法即可。调用方只关心接口，不关心具体类型，让我们轻松实现多态。

```go
// 定义接口
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 实现接口（隐式实现）
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// 多个类型实现同一接口
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// 接口作为参数
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", 
        s.Area(), s.Perimeter())
}
```

> 课堂提示：示例代码 `lesson-01/examples/advanced/01-interface.go` 在此基础上扩展了三角形实现、可直接运行的数据库示例以及 `WriterCloser` 组合接口的演示，运行 `go run lesson-01/examples/advanced/01-interface.go` 可以看到完整输出。

### 2. 空接口和类型断言

> 课堂提示：空接口 `interface{}` 在 Go 1.18 之前常被用来临时存放“任意类型”，如今更多用于需要和标准库互操作的场景（例如 `fmt` 包、`encoding/json`）以及处理确实不确定类型的数据。

- 空接口使用场景  
  - `fmt.Println`、`log.Print` 等变长参数需要接受任意类型。  
  - `encoding/json` 反序列化到 `map[string]interface{}`，让我们按键提取动态字段。  
  - 抽象事件总线、消息队列时，可以先用空接口承载载荷，后续配合断言或注册回调处理。
- 最佳实践提醒  
  - 提前约定空接口中可能出现的具体类型，配合文档或命名约束。  
  - 优先考虑具体类型或泛型（Go 1.18+），空接口是兜底方案。  
  - 使用类型断言或 `type switch` 时，一定要处理失败分支，避免 `panic`。

```go
// 空接口可以表示任何类型
var value interface{}

value = 42
value = "hello"
value = []int{1, 2, 3}

// 类型断言
func doSomething(v interface{}) {
    // 方式1：类型断言
    if str, ok := v.(string); ok {
        fmt.Println("是字符串:", str)
    }

    // 方式2：type switch
    switch data := v.(type) {
    case int:
        fmt.Println("整数:", data)
    case string:
        fmt.Println("字符串:", data)
    case []int:
        fmt.Println("整数切片:", data)
    default:
        fmt.Println("未知类型")
    }
}

func main() {
    // 常见用法：从 map[string]interface{} 中取值并断言具体类型
    payload := map[string]interface{}{
        "id":    1001,
        "name":  "golang",
        "extra": []string{"interface", "assertion"},
    }

    if id, ok := payload["id"].(int); ok {
        fmt.Println("ID:", id)
    }

    // 断言失败时要有兜底处理
    if tags, ok := payload["extra"].([]string); ok {
        fmt.Println("tags:", tags)
    } else {
        fmt.Println("extra 字段不是期望的 []string")
    }
}
```

### 3. 接口组合

```go
// 定义多个接口
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}

// 使用
type File struct {
    name string
}

func (f *File) Read(data []byte) (int, error) {
    // 实现读操作
    return 0, nil
}

func (f *File) Write(data []byte) (int, error) {
    // 实现写操作
    return 0, nil
}
```

### 4. 实战：实现多态性系统

```go
// 定义一个数据库接口
type Database interface {
    Connect() error
    Query(sql string) (interface{}, error)
    Close() error
}

// 实现MySQL数据库
type MySQL struct {
    connection string
}

func (m *MySQL) Connect() error {
    // 连接MySQL
    return nil
}

func (m *MySQL) Query(sql string) (interface{}, error) {
    // 执行查询
    return []string{"result1", "result2"}, nil
}

func (m *MySQL) Close() error {
    // 关闭连接
    return nil
}

// 实现PostgreSQL数据库
type PostgreSQL struct {
    connection string
}

func (p *PostgreSQL) Connect() error {
    // 连接PostgreSQL
    return nil
}

func (p *PostgreSQL) Query(sql string) (interface{}, error) {
    // 执行查询
    return []string{"result1", "result2"}, nil
}

func (p *PostgreSQL) Close() error {
    // 关闭连接
    return nil
}

// 多态使用
func ExecuteQuery(db Database, sql string) {
    db.Connect()
    defer db.Close()
    
    result, err := db.Query(sql)
    if err != nil {
        panic(err)
    }
    fmt.Println("Result:", result)
}
```

---

## 第二部分：Goroutine入门

### 实战演示：课程示例代码

**示例路径：** `lesson-01/examples/advanced/02-goroutine.go`

**运行方式：**
```bash
go run lesson-01/examples/advanced/02-goroutine.go
```

**示例覆盖的核心知识点：**
- `basicGoroutine()`：并发启动多个goroutine，体验调度顺序的不确定性。
- `waitGroupDemo()`：使用`sync.WaitGroup`等待协程完成，避免主协程提前退出。
- `channelDemo()` 与 `bufferedChannelDemo()`：演示无缓冲与有缓冲channel的发送、接收与关闭。
- `selectDemo()`、`timeoutDemo()`、`nonBlockingDemo()`：展示`select`在多路复用、超时控制与非阻塞通信中的用法。
- `loopSelectDemo()`、`quitChannelDemo()`、`closedChannelDemo()`：演示在循环监听中动态屏蔽已关闭channel、使用退出信号以及读取关闭后零值。
- `fairnessDemo()`：展示`select`面对多个就绪case时的公平性。

> 课堂建议：运行示例时，可以注释/解注部分函数调用，观察输出顺序的变化，加深对调度和同步机制的理解。

### 1. Goroutine的创建和使用

```go
// 基本的goroutine
func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // 启动goroutine
    go sayHello()
    
    // 等待goroutine完成（简单方式）
    time.Sleep(time.Second)
}

// 使用WaitGroup等待
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d\n", id)
        }(i)
    }
    
    wg.Wait()
}

// 协程泄漏示例（错误的做法）
func leakExample() {
    ch := make(chan int)
    
    go func() {
        ch <- 1 // 永远阻塞，因为没有人接收
    }()
    
    // 应该使用缓冲channel或者接收值
}
```

### 2. Channel通信机制

```go
// 无缓冲channel（同步）
func channelDemo() {
    ch := make(chan int)
    
    go func() {
        ch <- 42 // 发送数据
    }()
    
    value := <-ch // 接收数据
    fmt.Println(value)
}

// 缓冲channel（异步）
func bufferedChannelDemo() {
    ch := make(chan int, 3) // 缓冲区大小为3
    
    ch <- 1
    ch <- 2
    ch <- 3
    // ch <- 4 // 这里会阻塞，因为缓冲区满了
    
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}

// channel关闭
func closeChannelDemo() {
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)
    
    // 关闭后的channel仍然可以读取
    // 不能再发送数据
    for value := range ch {
        fmt.Println(value)
    }
}
```

> 示例源码中的 `closedChannelDemo()` 使用 `select` + `default` 的写法，帮助学员看到关闭 channel 后立即返回零值且 `ok=false` 的行为。

### 3. Select语句

#### 3.1 什么是Select？

`select` 是Go语言中用于处理多个channel操作的专用语句，类似于 `switch` 语句，但专门用于channel通信。

**核心概念：**
- `select` 允许goroutine同时等待多个channel操作
- 它会阻塞，直到其中一个case可以执行
- 如果有多个case同时就绪，会**随机选择**一个执行
- 可以用于实现超时、非阻塞操作等并发模式

**基本语法：**
```go
select {
case msg1 := <-ch1:
    // 处理ch1的消息
case msg2 := <-ch2:
    // 处理ch2的消息
case ch3 <- value:
    // 向ch3发送数据
default:
    // 如果所有case都阻塞，执行default（可选）
}
```

#### 3.2 Select的工作原理

**执行流程：**

1. **检查所有case**：`select` 会检查每个case中的channel操作是否可以立即执行
2. **选择就绪的case**：
   - 如果有多个case就绪，**随机选择**一个执行
   - 如果只有一个case就绪，执行该case
   - 如果所有case都阻塞，执行`default`（如果有）
3. **阻塞等待**：如果没有`default`且所有case都阻塞，`select`会阻塞，直到某个case就绪

**关键特性：**
- **随机选择**：当多个case同时就绪时，Go会随机选择一个，这保证了公平性
- **非阻塞**：使用`default`可以实现非阻塞操作
- **超时控制**：结合`time.After`可以实现超时机制

#### 3.3 Select的使用场景

**1. 多路复用（Multiplexing）**

同时监听多个channel，哪个先有数据就处理哪个：

```go
func selectDemo() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from ch2"
    }()
    
    // 随机选择一个就绪的channel
    select {
    case msg1 := <-ch1:
        fmt.Println(msg1)  // 可能输出："from ch1"
    case msg2 := <-ch2:
        fmt.Println(msg2)  // 可能输出："from ch2"
    }
    // 输出取决于哪个channel先有数据
}
```

**2. 超时控制（Timeout）**

为channel操作设置超时时间，避免无限等待：

```go
func timeoutDemo() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case msg := <-ch:
        fmt.Println("收到:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("超时了")  // 1秒后输出这个
    }
    // 因为发送需要2秒，但超时是1秒，所以会输出"超时了"
}
```

> 示例 `timeoutDemo()` 中发送端刻意延迟 2 秒，超时设置为 1 秒，现场演示时会输出“超时了”，便于说明 `time.After` 的用法。

**3. 非阻塞操作（Non-blocking）**

使用`default`实现非阻塞的channel操作：

```go
func nonBlockingDemo() {
    ch := make(chan int)
    
    // 非阻塞接收
    select {
    case value := <-ch:
        fmt.Println("收到:", value)
    default:
        fmt.Println("没有值可读（非阻塞）")  // 立即输出这个
    }
    
    // 非阻塞发送
    select {
    case ch <- 42:
        fmt.Println("发送成功")
    default:
        fmt.Println("channel已满，发送失败")
    }
}
```

**4. 循环监听（Loop with Select）**

在循环中使用`select`持续监听多个channel：

```go
func loopSelect() {
    ch1 := make(chan int)
    ch2 := make(chan string)
    done := make(chan bool)
    
    go func() {
        for i := 0; i < 5; i++ {
            ch1 <- i
            time.Sleep(100 * time.Millisecond)
        }
        close(ch1)
    }()
    
    go func() {
        for i := 0; i < 3; i++ {
            ch2 <- fmt.Sprintf("msg-%d", i)
            time.Sleep(150 * time.Millisecond)
        }
        close(ch2)
    }()
    
    // 持续监听，直到所有channel都关闭
    for {
        select {
        case val, ok := <-ch1:
            if !ok {
                ch1 = nil  // 关闭的channel设为nil，select会忽略它
                continue
            }
            fmt.Println("ch1:", val)
        case msg, ok := <-ch2:
            if !ok {
                ch2 = nil
                continue
            }
            fmt.Println("ch2:", msg)
        default:
            // 如果没有数据，可以做其他事情
            if ch1 == nil && ch2 == nil {
                fmt.Println("所有channel已关闭")
                return
            }
        }
    }
}
```

**5. 退出信号（Quit Channel）**

使用`select`处理退出信号：

```go
func quitChannelDemo() {
    jobs := make(chan int)
    quit := make(chan bool)
    
    // 工作goroutine
    go func() {
        for {
            select {
            case job := <-jobs:
                fmt.Printf("处理任务: %d\n", job)
            case <-quit:
                fmt.Println("收到退出信号")
                return
            }
        }
    }()
    
    // 发送任务
    for i := 0; i < 5; i++ {
        jobs <- i
    }
    
    // 发送退出信号
    quit <- true
    time.Sleep(100 * time.Millisecond)
}
```

#### 3.4 Select的执行机制详解

**1. 随机选择的公平性**

当多个case同时就绪时，Go会**随机选择**一个，这保证了公平性：

```go
func fairnessDemo() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    // 同时向两个channel发送数据
    go func() {
        for i := 0; i < 10; i++ {
            ch1 <- i
        }
        close(ch1)
    }()
    
    go func() {
        for i := 0; i < 10; i++ {
            ch2 <- i * 10
        }
        close(ch2)
    }()
    
    // 持续接收，随机选择就绪的channel
    for i := 0; i < 20; i++ {
        select {
        case val, ok := <-ch1:
            if !ok {
                ch1 = nil
                continue
            }
            fmt.Printf("ch1: %d\n", val)
        case val, ok := <-ch2:
            if !ok {
                ch2 = nil
                continue
            }
            fmt.Printf("ch2: %d\n", val)
        }
    }
    // 输出顺序是随机的，体现了公平性
}
```

> 源码 `fairnessDemo()` 会同时关闭两个 channel 并持续 `select`，运行后可以观察不同运行下的输出顺序，便于说明公平性。

**2. Select的阻塞行为**

```go
// 没有default：会阻塞
select {
case <-ch1:
    // ...
case <-ch2:
    // ...
}
// 如果ch1和ch2都阻塞，select会一直等待

// 有default：不会阻塞
select {
case <-ch1:
    // ...
case <-ch2:
    // ...
default:
    // 立即执行这个
}
// 如果ch1和ch2都阻塞，立即执行default
```

**3. 关闭channel的处理**

当channel关闭后：
- 从关闭的channel接收数据，会立即返回零值，`ok`为`false`
- 向关闭的channel发送数据，会panic
- 在select中，关闭的channel可以设置为`nil`，select会忽略它

```go
func closedChannelDemo() {
    ch := make(chan int)
    close(ch)
    
    select {
    case val, ok := <-ch:
        fmt.Printf("val: %d, ok: %v\n", val, ok)  // val: 0, ok: false
    default:
        fmt.Println("default")
    }
}
```

#### 3.5 Select的常见模式

**模式1：超时模式**
```go
select {
case result := <-ch:
    // 处理结果
case <-time.After(5 * time.Second):
    // 超时处理
    return errors.New("操作超时")
}
```

**模式2：取消模式**
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

select {
case result := <-ch:
    // 处理结果
case <-ctx.Done():
    // 取消操作
    return ctx.Err()
}
```

**模式3：优先级模式**
```go
select {
case urgent := <-urgentChan:
    // 优先处理紧急任务
    handleUrgent(urgent)
case normal := <-normalChan:
    // 处理普通任务
    handleNormal(normal)
default:
    // 没有任务时，可以执行其他操作
    doIdleWork()
}
```

**模式4：心跳模式**
```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        // 定期执行的任务
        heartbeat()
    case msg := <-messages:
        // 处理消息
        handleMessage(msg)
    }
}
```

#### 3.6 Select的注意事项

**1. 避免空的select**
```go
// ❌ 错误：空的select会永远阻塞
select {}

// ✅ 正确：至少有一个case
select {
case <-ch:
    // ...
}
```

**2. 处理关闭的channel**
```go
// ✅ 正确：检查channel是否关闭
select {
case val, ok := <-ch:
    if !ok {
        ch = nil  // 设置为nil，select会忽略它
        continue
    }
    // 处理val
}
```

**3. 避免在select中发送到nil channel**
```go
var ch chan int  // nil channel

// ⚠️ 注意：向nil channel发送会永远阻塞
select {
case ch <- 42:  // 这会永远阻塞
    fmt.Println("发送成功")
default:
    fmt.Println("不会执行到这里")
}
```

**4. 在循环中使用select**
```go
// ✅ 正确：使用break退出select，继续循环
for {
    select {
    case <-ch:
        // 处理
        break  // 只退出select，不退出for循环
    }
}

// ✅ 正确：使用标签退出外层循环
loop:
for {
    select {
    case <-ch:
        break loop  // 退出外层for循环
    }
}
```

**5. Select的性能考虑**
- `select`本身性能开销很小
- 但频繁使用`default`可能浪费CPU（忙等待）
- 在循环中使用`select`时，可以考虑添加适当的延迟

#### 3.7 实战示例

```go
// 基本select
func selectDemo() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from ch1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from ch2"
    }()
    
    // 随机选择一个就绪的channel
    select {
    case msg1 := <-ch1:
        fmt.Println(msg1)
    case msg2 := <-ch2:
        fmt.Println(msg2)
    }
}

// select with timeout
func timeoutDemo() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case msg := <-ch:
        fmt.Println("收到:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("超时了")
    }
}

// select with default（非阻塞）
func nonBlockingDemo() {
    ch := make(chan int)
    
    select {
    case value := <-ch:
        fmt.Println(value)
    default:
        fmt.Println("没有值可读")
    }
}
```

### 4. 实战：生产者消费者模式

```go
func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int, id int) {
    for value := range ch {
        fmt.Printf("Consumer %d received: %d\n", id, value)
        time.Sleep(50 * time.Millisecond)
    }
}

func main() {
    ch := make(chan int, 5)
    
    go producer(ch)
    
    // 多个消费者
    for i := 0; i < 3; i++ {
        go consumer(ch, i)
    }
    
    time.Sleep(2 * time.Second)
}
```

**拓展示例：使用 `goroutine` + `select` 管理停止信号**

```go
func producer(ctx context.Context, id int, queue chan<- int) {
    defer log.Printf("producer %d exit", id)

    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        item := rand.Intn(1000)

        select {
        case queue <- item:
            log.Printf("producer %d -> %d\n", id, item)
        case <-ctx.Done():
            return
        }

        time.Sleep(300 * time.Millisecond)
    }
}

func consumer(ctx context.Context, id int, queue <-chan int) {
    defer log.Printf("consumer %d exit", id)

    for {
        select {
        case item := <-queue:
            log.Printf("consumer %d <- %d\n", id, item)
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    queue := make(chan int, 5)
    ctx, cancel := context.WithCancel(context.Background())

    for i := 0; i < 2; i++ {
        go producer(ctx, i+1, queue)
    }

    for i := 0; i < 3; i++ {
        go consumer(ctx, i+1, queue)
    }

    time.Sleep(3 * time.Second)
    cancel()
    time.Sleep(500 * time.Millisecond)

    fmt.Println("done")
}
```

> 要点：`queue` 用作缓冲队列，`select` 既能处理数据通路，又能监听 `ctx.Done()` 实现优雅退出，避免 goroutine 泄漏。

---

## 第三部分：并发安全

### 1. Mutex和RWMutex

```go
// 使用Mutex保护共享资源
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (sc *SafeCounter) Increment() {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.count++
}

func (sc *SafeCounter) GetCount() int {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    return sc.count
}

// 使用RWMutex（读多写少场景）
type SafeData struct {
    mu     sync.RWMutex
    data   map[string]int
}

func (sd *SafeData) Read(key string) int {
    sd.mu.RLock() // 读锁
    defer sd.mu.RUnlock()
    return sd.data[key]
}

func (sd *SafeData) Write(key string, value int) {
    sd.mu.Lock() // 写锁
    defer sd.mu.Unlock()
    sd.data[key] = value
}
```

### 2. WaitGroup同步

```go
func waitGroupDemo() {
    var wg sync.WaitGroup
    mu := sync.Mutex{}
    sum := 0
    
    // 启动多个goroutine
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // 执行一些工作
            time.Sleep(100 * time.Millisecond)
            
            // 更新共享变量
            mu.Lock()
            sum += id
            mu.Unlock()
        }(i)
    }
    
    wg.Wait() // 等待所有goroutine完成
    fmt.Println("Sum:", sum)
}
```

### 3. Context上下文控制

#### 3.1 什么是Context？

`context.Context` 是Go语言标准库中用于**跨goroutine传递取消信号、超时、截止时间和请求范围值**的标准方式。它是并发控制的核心工具。

**核心概念：**
- Context在多个goroutine之间传播控制信号
- Context是不可变的，每次派生都会创建新的Context
- Context是线程安全的，可以安全地在多个goroutine中使用
- Context形成树形结构，父Context取消时，所有子Context也会被取消

**主要用途：**
1. **取消控制**：取消长时间运行的操作
2. **超时控制**：为操作设置超时时间
3. **截止时间**：设置操作必须完成的最后期限
4. **传递值**：在请求范围内传递元数据（如trace ID、用户ID等）

#### 3.2 Context的类型

Go提供了几种创建Context的方法：

**1. context.Background()**

根Context，通常用于main函数、初始化或测试中。它永远不会被取消、没有值、没有截止时间。

```go
func main() {
    ctx := context.Background()
    // 作为根context使用
}
```

**2. context.TODO()**

当不确定使用哪个Context时使用，通常是占位符，表示"稍后会替换成真正的Context"。

```go
func someFunction() {
    ctx := context.TODO()  // 待完善
    // ...
}
```

**3. context.WithCancel(parent)**

创建可取消的Context，返回Context和cancel函数。调用cancel函数会取消该Context及其所有子Context。

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // 确保释放资源
```

**4. context.WithTimeout(parent, timeout)**

创建有超时时间的Context，超时后自动取消。等价于`WithDeadline(parent, time.Now().Add(timeout))`。

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // 即使未超时也应该调用，释放资源
```

**5. context.WithDeadline(parent, deadline)**

创建有截止时间的Context，到达截止时间后自动取消。

```go
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

**6. context.WithValue(parent, key, value)**

创建携带键值对的Context，用于传递请求范围的数据。**不要滥用！**

```go
type contextKey string
const userIDKey contextKey = "userID"

ctx := context.WithValue(context.Background(), userIDKey, "user123")
```

#### 3.3 Context的基本使用

**1. 取消控制示例**

```go
func cancellableDemo() {
	fmt.Println("=== 可取消的Context ===")
	
	ctx, cancel := context.WithCancel(context.Background())
	
	// 启动工作goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine收到取消信号:", ctx.Err())
				return
			default:
				fmt.Println("工作中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	// 工作2秒后取消
	time.Sleep(2 * time.Second)
	fmt.Println("发送取消信号")
	cancel()
	
	// 等待goroutine退出
	time.Sleep(500 * time.Millisecond)
}
```

**2. 超时控制示例**

```go
func timeoutContextDemo() {
	fmt.Println("=== 超时Context ===")
	
	// 设置1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	ch := make(chan string)
	
	// 模拟一个耗时2秒的操作
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()
	
	select {
	case result := <-ch:
		fmt.Println("收到结果:", result)
	case <-ctx.Done():
		fmt.Println("操作超时:", ctx.Err())  // 输出: context deadline exceeded
	}
}
```

**3. 截止时间示例**

```go
func deadlineContextDemo() {
	fmt.Println("=== 截止时间Context ===")
	
	// 设置3秒后的截止时间
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	
	// 检查剩余时间
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("截止时间: %v, 剩余: %v\n", d, time.Until(d))
	}
	
	// 等待超过截止时间
	time.Sleep(4 * time.Second)
	
	select {
	case <-ctx.Done():
		fmt.Println("已超过截止时间:", ctx.Err())
	default:
		fmt.Println("未超时")
	}
}
```

**4. 传递值示例**

```go
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

func valueContextDemo() {
	fmt.Println("=== Context传递值 ===")
	
	// 创建携带值的context
	ctx := context.Background()
	ctx = context.WithValue(ctx, requestIDKey, "req-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")
	
	// 在函数中读取值
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	if reqID := ctx.Value(requestIDKey); reqID != nil {
		fmt.Printf("Request ID: %v\n", reqID)
	}
	
	if userID := ctx.Value(userIDKey); userID != nil {
		fmt.Printf("User ID: %v\n", userID)
	}
}
```

#### 3.4 Context的高级使用场景

**场景1：级联取消（父取消，子也取消）**

```go
func cascadeCancelDemo() {
	fmt.Println("=== 级联取消示例 ===")
	
	// 创建父context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()
	
	// 创建子context
	childCtx1, cancel1 := context.WithCancel(parentCtx)
	defer cancel1()
	
	childCtx2, cancel2 := context.WithCancel(parentCtx)
	defer cancel2()
	
	// 启动子goroutine
	go worker(childCtx1, "Worker 1")
	go worker(childCtx2, "Worker 2")
	
	time.Sleep(1 * time.Second)
	
	// 取消父context，所有子context也会被取消
	fmt.Println("取消父context")
	parentCancel()
	
	time.Sleep(500 * time.Millisecond)
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: 收到取消信号\n", name)
			return
		default:
			fmt.Printf("%s: 工作中...\n", name)
			time.Sleep(300 * time.Millisecond)
		}
	}
}
```

**场景2：HTTP请求超时控制**

```go
func httpRequestDemo() {
	// 创建5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("请求超时")
		} else {
			fmt.Println("请求失败:", err)
		}
		return
	}
	defer resp.Body.Close()
	
	fmt.Println("请求成功:", resp.StatusCode)
}
```

**场景3：数据库查询超时**

```go
func databaseQueryDemo(db *sql.DB) {
	// 创建3秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// 执行查询
	rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE active = ?", true)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("数据库查询超时")
		} else {
			fmt.Println("查询失败:", err)
		}
		return
	}
	defer rows.Close()
	
	// 处理结果
	for rows.Next() {
		// ...
	}
}
```

**场景4：多个goroutine协同工作**

```go
func multiWorkerDemo() {
	fmt.Println("=== 多worker协同工作 ===")
	
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	
	// 启动多个worker
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d: 退出\n", id)
					return
				default:
					fmt.Printf("Worker %d: 处理任务\n", id)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(i)
	}
	
	// 工作2秒后取消所有worker
	time.Sleep(2 * time.Second)
	fmt.Println("发送取消信号给所有worker")
	cancel()
	
	// 等待所有worker退出
	wg.Wait()
	fmt.Println("所有worker已退出")
}
```

**场景5：Context在Pipeline中的应用**

```go
func pipelineDemo() {
	fmt.Println("=== Pipeline示例 ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// Stage 1: 生成数据
	dataCh := generateData(ctx)
	
	// Stage 2: 处理数据
	processedCh := processData(ctx, dataCh)
	
	// Stage 3: 输出结果
	for result := range processedCh {
		fmt.Println("最终结果:", result)
	}
	
	fmt.Println("Pipeline完成")
}

func generateData(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("生成器: 收到取消信号")
				return
			case ch <- i:
				fmt.Println("生成器: 生成", i)
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	return ch
}

func processData(ctx context.Context, input <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for data := range input {
			select {
			case <-ctx.Done():
				fmt.Println("处理器: 收到取消信号")
				return
			case ch <- data * 2:
				fmt.Println("处理器: 处理", data, "->", data*2)
			}
		}
	}()
	return ch
}
```

#### 3.5 Context的最佳实践

**1. 传递Context的规则**

✅ **正确做法：**
```go
// 将Context作为第一个参数传递
func DoSomething(ctx context.Context, arg1 string, arg2 int) error {
	// ...
}

// 方法中也应遵循相同规则
func (s *Service) Process(ctx context.Context, data string) error {
	// ...
}
```

❌ **错误做法：**
```go
// 不要把Context放在struct中
type Service struct {
	ctx context.Context  // 不推荐
	// ...
}

// 不要把Context作为最后一个参数
func DoSomething(arg1 string, arg2 int, ctx context.Context) error {
	// ...
}
```

**2. 不要存储Context**

❌ **错误做法：**
```go
type Server struct {
	ctx context.Context  // 不要这样做
}

func (s *Server) Start() {
	s.ctx = context.Background()  // 不要这样做
}
```

✅ **正确做法：**
```go
type Server struct {
	// 不存储context
}

func (s *Server) Start(ctx context.Context) {
	// 直接使用传入的context
}
```

**3. 不要传递nil Context**

❌ **错误做法：**
```go
func DoWork(ctx context.Context) {
	// ...
}

// 调用时
DoWork(nil)  // 不要传nil
```

✅ **正确做法：**
```go
// 如果不确定用什么context，使用context.TODO()
DoWork(context.TODO())

// 或者使用context.Background()
DoWork(context.Background())
```

**4. Context.Value使用原则**

只用于传递请求范围的数据，不要用于传递可选参数。

✅ **适合使用Context.Value的场景：**
- Request ID（请求追踪）
- User ID（用户认证）
- Trace ID（分布式追踪）
- Request-scoped metadata

❌ **不适合使用Context.Value的场景：**
- 函数的可选参数
- 配置信息
- 依赖注入

```go
// ✅ 正确：传递请求范围的元数据
type contextKey string
const requestIDKey contextKey = "requestID"

func handler(w http.ResponseWriter, r *http.Request) {
	requestID := generateRequestID()
	ctx := context.WithValue(r.Context(), requestIDKey, requestID)
	
	processRequest(ctx, r)
}

// ❌ 错误：用Context传递配置
ctx := context.WithValue(context.Background(), "config", config)  // 不要这样做
```

**5. 总是调用cancel函数**

```go
// ✅ 正确：使用defer确保cancel被调用
func DoWork() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()  // 即使未超时，也要调用cancel释放资源
	
	// ...
}

// ❌ 错误：忘记调用cancel
func DoWork() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 忘记调用cancel，会造成资源泄漏
	
	// ...
}
```

**6. 检查Context是否已取消**

```go
func DoWork(ctx context.Context) error {
	// 在开始工作前检查
	if err := ctx.Err(); err != nil {
		return err
	}
	
	// 在长时间操作中定期检查
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// 执行工作
		doSomething(i)
	}
	
	return nil
}
```

**7. Context错误处理**

```go
func handleContextError(ctx context.Context) {
	err := ctx.Err()
	
	switch err {
	case context.Canceled:
		fmt.Println("操作被取消")
	case context.DeadlineExceeded:
		fmt.Println("操作超时")
	case nil:
		fmt.Println("Context仍然有效")
	default:
		fmt.Println("未知错误:", err)
	}
}
```

#### 3.6 Context实战示例

**完整示例：构建一个支持取消和超时的任务系统**

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	ID       int
	Duration time.Duration
}

// TaskManager 任务管理器
type TaskManager struct {
	tasks []Task
	mu    sync.Mutex
}

// AddTask 添加任务
func (tm *TaskManager) AddTask(task Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks = append(tm.tasks, task)
}

// ExecuteTask 执行单个任务
func (tm *TaskManager) ExecuteTask(ctx context.Context, task Task) error {
	fmt.Printf("任务 %d: 开始执行\n", task.ID)
	
	// 模拟任务执行
	select {
	case <-time.After(task.Duration):
		fmt.Printf("任务 %d: 执行完成\n", task.ID)
		return nil
	case <-ctx.Done():
		fmt.Printf("任务 %d: 被取消 (%v)\n", task.ID, ctx.Err())
		return ctx.Err()
	}
}

// ExecuteAll 执行所有任务
func (tm *TaskManager) ExecuteAll(ctx context.Context) {
	var wg sync.WaitGroup
	
	tm.mu.Lock()
	tasks := tm.tasks
	tm.mu.Unlock()
	
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			tm.ExecuteTask(ctx, t)
		}(task)
	}
	
	wg.Wait()
	fmt.Println("所有任务处理完成")
}

func main() {
	// 创建任务管理器
	tm := &TaskManager{}
	
	// 添加任务
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 5; i++ {
		tm.AddTask(Task{
			ID:       i,
			Duration: time.Duration(rand.Intn(3)+1) * time.Second,
		})
	}
	
	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	
	// 启动监听取消信号的goroutine
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n⚠️  手动触发取消信号")
		cancel()
	}()
	
	// 执行所有任务
	fmt.Println("开始执行任务...")
	tm.ExecuteAll(ctx)
	
	// 检查context状态
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("\n❌ 任务执行超时")
	} else if ctx.Err() == context.Canceled {
		fmt.Println("\n❌ 任务被手动取消")
	} else {
		fmt.Println("\n✅ 所有任务正常完成")
	}
}
```

#### 3.7 Context常见问题

**Q1: Context.Value会有性能问题吗？**

A: Context.Value的查找是链式查找，性能是O(n)，n是context链的深度。因此：
- 不要嵌套太深
- 不要频繁调用Value
- 考虑在函数入口处一次性提取所需的值

**Q2: 为什么不能在struct中存储Context？**

A: Context的生命周期通常绑定到请求，存储在struct中会导致：
- 生命周期混乱
- 难以控制取消信号的传播
- 可能导致goroutine泄漏

**Q3: Context.WithValue的key为什么要用自定义类型？**

A: 使用自定义类型作为key可以避免不同包之间的key冲突。

```go
// ✅ 正确
type contextKey string
const myKey contextKey = "myKey"

// ❌ 不推荐
const myKey = "myKey"  // string类型容易冲突
```

**Q4: 如何在多层函数调用中传递Context？**

A: 将Context作为第一个参数逐层传递：

```go
func Handler(ctx context.Context) {
	// 传递给下一层
	Service(ctx)
}

func Service(ctx context.Context) {
	// 传递给下一层
	Repository(ctx)
}

func Repository(ctx context.Context) {
	// 使用context
}
```

> **配套示例：** `lesson-01/examples/advanced/03-context.go` 包含了本节所有示例的可运行代码。运行 `go run lesson-01/examples/advanced/03-context.go` 可以看到完整的演示效果。

### 4. 实战：并发安全的计数器

```go
type SafeCounter struct {
    mu       sync.RWMutex
    counters map[string]int
}

func NewSafeCounter() *SafeCounter {
    return &SafeCounter{
        counters: make(map[string]int),
    }
}

func (sc *SafeCounter) Increment(key string) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.counters[key]++
}

func (sc *SafeCounter) Get(key string) int {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.counters[key]
}

func (sc *SafeCounter) GetAll() map[string]int {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    
    result := make(map[string]int)
    for k, v := range sc.counters {
        result[k] = v
    }
    return result
}
```

---

## 第四部分：GMP调度机制

### 1. Goroutine背后的执行模型

Go 运行时采用 **G-M-P（Goroutine、Machine、Processor）调度模型** 支撑高并发。三者职责分别是：

- `G`（Goroutine）：用户级协程，包含栈、状态等上下文，是被调度的基本单元。
- `M`（Machine）：映射到操作系统线程，用于真正执行 `G`。
- `P`（Processor）：逻辑处理器，持有可运行 `G` 的本地队列，并维护内存分配缓存。

只有当 `G` 绑定到 `P`，再由拥有该 `P` 的某个 `M` 执行时，协程才会真正运行，形成 `G → P → M` 的执行链路。[参考](https://go.cyub.vip/gmp/gmp-model/)

### 2. 调度关键机制

- **本地队列优先**：`M` 首先从自身绑定的 `P` 的本地队列中取 `G`，避免全局锁。
- **Work Stealing**：本地队列为空时，`M` 会从其他 `P` 窃取一半就绪 `G`，保证均衡。
- **全局队列兜底**：全局可运行队列确保没有 `P` 被饿死。
- **自旋复用线程**：没有可运行 `G` 时，`M` 会短暂自旋等待，减少频繁创建/销毁线程。
- **Hand Off 机制**：当 `G` 因系统调用阻塞，`P` 会解绑并交给其他空闲 `M`，维持整体吞吐。

### 3. 实用调试建议

- 通过 `GODEBUG=schedtrace=1000` 观察 `gomaxprocs`、`threads`、`runqueue` 等调度指标。
- 使用 `go tool trace` 获取更详细的时间线视图，分析 `G` 的生命周期。
- 结合 `WaitGroup`、`context` 管理 `G` 的退出和取消，可以更好地配合调度器。[参考](https://go.cyub.vip/gmp/gmp-model/)

**配套示例：** `lesson-01/examples/advanced/06-gmp.go` 会创建一批 CPU 密集型 goroutine，并调整 `GOMAXPROCS`。可搭配命令  
`GODEBUG=schedtrace=1000,scheddetail=1 go run lesson-01/examples/advanced/06-gmp.go` 现场观察调度日志。

课堂演示建议：结合 `runtime.GOMAXPROCS`、`schedtrace` 输出以及示例 `loopSelectDemo()`、`workerPool` 等场景，让学员直观看到调度模型对任务分发与线程利用率的影响。

#### 日志解析示例

课堂可先运行 `06-gmp.go` 并捕获 `schedtrace` 日志（示例日志存放在 `lesson-01/examples/advanced/gmp.log`）。引导学员按以下顺序阅读：

1. **时间与概要行**（如 `SCHED 1000ms`）：关注 `gomaxprocs`、`threads`、`idleprocs`、`runqueue` 等全局指标，判断当前是否存在排队或空闲。
2. **P 列表**：观察 `status`、`runqsize`、`m` 等字段，讨论为什么某些 `P` 空闲、某些 `P` 上有 backlog，以及 `syscalltick` 激增意味着什么。
3. **M 列表**：确认线程是否处于 `spinning` 或 `blocked`，说明 Goroutine 进入阻塞时 runtime 如何复用/挂起 M。
4. **G 列表**：结合 `status=1/2/4`，让学员识别 `sync.WaitGroup.Wait`、`chan send/receive` 等典型阻塞场景，理解 G 与 M、P 的绑定关系。

延伸问题：调低 `schedtrace` 的采样间隔（如 `schedtrace=100`）或调节 `GOMAXPROCS`，让学员比较日志变化，巩固 Work Stealing 和 Hand Off 的概念。[参考](https://go.cyub.vip/gmp/gmp-model/)

---

## 第五部分：并发编程最佳实践（5分钟）

### 1. 常见错误和陷阱

```go
// ❌ 错误1：没有等待goroutine
go someFunc() // goroutine还没执行完，main就结束了

// ✅ 正确
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    someFunc()
}()
wg.Wait()

// ❌ 错误2：闭包捕获循环变量
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i) // 可能输出都是3
    }()
}

// ✅ 正确
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println(id) // 输出正确的值
    }(i)
}

// ❌ 错误3：channel没有关闭
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    // 忘记close(ch)
}()

// ❌ 错误4：在goroutine中不使用recover
go func() {
    panic("oops") // 会导致整个程序崩溃
}()

// ✅ 正确
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("捕获panic:", r)
        }
    }()
    // 可能panic的代码
}()
```

### 2. 性能优化建议

```go
// 1. 避免频繁创建goroutine
// ❌ 不好：在循环中频繁创建
for _, item := range items {
    go func(it Item) {
        process(it)
    }(item)
}

// ✅ 更好：使用worker pool
func workerPool(items []Item, numWorkers int) {
    jobs := make(chan Item, len(items))
    
    // 启动worker
    for w := 0; w < numWorkers; w++ {
        go func() {
            for item := range jobs {
                process(item)
            }
        }()
    }
    
    // 发送任务
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // 等待完成...
}

// 2. 使用缓冲channel提高性能
ch := make(chan int, 100) // 根据实际情况设置缓冲区大小

// 3. 合理使用context超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## 第六部分：标准库与项目实战

### 1. 常用标准库

```go
// JSON编解码
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func jsonDemo() {
    p := Person{Name: "Alice", Age: 30}
    
    // 编码
    data, _ := json.Marshal(p)
    fmt.Println(string(data))
    
    // 解码
    var p2 Person
    json.Unmarshal(data, &p2)
    fmt.Println(p2)
}

// 文件操作
func fileDemo() {
    // 读取文件
    data, err := ioutil.ReadFile("file.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
    
    // 写入文件
    err = ioutil.WriteFile("output.txt", []byte("Hello"), 0644)
    if err != nil {
        panic(err)
    }
}

// 时间处理
func timeDemo() {
    now := time.Now()
    fmt.Println("当前时间:", now)
    fmt.Println("格式化:", now.Format("2006-01-02 15:04:05"))
    
    // 解析时间
    t, _ := time.Parse("2006-01-02", "2024-01-01")
    fmt.Println("解析的时间:", t)
    
    // 计算时间差
    duration := time.Now().Sub(t)
    fmt.Println("时间差:", duration)
}

// 加密哈希
func hashDemo() {
    data := "Hello World"
    
    // SHA256
    h := sha256.New()
    h.Write([]byte(data))
    fmt.Printf("SHA256: %x\n", h.Sum(nil))
    
    // MD5
    m := md5.New()
    m.Write([]byte(data))
    fmt.Printf("MD5: %x\n", m.Sum(nil))
}
```

### 2. 项目实战：构建一个简单的Web爬虫

参见实战代码：`examples/practice/web-crawler.go`

---

## 课后作业

1. **接口练习：**
   - 实现一个支付系统
   - 支持多种支付方式（支付宝、微信、银行卡）
   - 使用接口实现多态

2. **并发练习：**
   - 实现一个并发安全的日志系统
   - 使用goroutine异步写入日志
   - 使用WaitGroup确保所有日志都被写入

3. **综合练习：**
   - 实现一个简单的任务调度器
   - 支持并发执行多个任务
   - 使用context控制任务超时

---

## 总结

通过本部分学习，你应该掌握：

1. ✅ 接口的定义和实现
2. ✅ 空接口和类型断言
3. ✅ 接口组合和多态
4. ✅ Goroutine的创建和使用
5. ✅ Channel的通信机制
6. ✅ Select语句的使用
7. ✅ GMP调度机制原理
8. ✅ 并发安全（Mutex、RWMutex）
9. ✅ WaitGroup同步
10. ✅ Context上下文控制
11. ✅ 并发编程的最佳实践
12. ✅ 常用标准库的使用

**恭喜！** 你已经掌握了Go语言的核心语法和并发编程！
