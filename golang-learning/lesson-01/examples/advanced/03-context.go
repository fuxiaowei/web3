// Package main 提供Context上下文控制的完整示例集
//
// 本文件包含12个由浅入深的Context使用示例，涵盖：
//   - 基础用法：取消、超时、截止时间、传递值
//   - 进阶应用：级联取消、多Worker协同、Pipeline
//   - 实用场景：HTTP请求、数据库查询、错误处理
//   - 综合实战：任务管理系统、Worker Pool
//
// 学习建议：
//  1. 按顺序运行每个示例，理解基本概念
//  2. 修改参数（如超时时间），观察行为变化
//  3. 阅读课件 lesson-01/courseware/advanced.md 第三部分
//  4. 查看 README-Context.md 了解更多最佳实践
//
// 运行方式：
//
//	go run lesson-01/examples/advanced/03-context.go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============ 1. 可取消的Context ============
// cancellableDemo 演示如何使用WithCancel创建可取消的Context
// 关键点：
//  1. WithCancel返回context和cancel函数
//  2. 调用cancel()会关闭ctx.Done()通道
//  3. 所有监听ctx.Done()的goroutine都会收到信号
func cancellableDemo() {
	fmt.Println("=== 可取消的Context ===")

	// 创建一个可取消的context
	// Background()返回一个空的根context，永远不会被取消
	ctx, cancel := context.WithCancel(context.Background())

	// 启动工作goroutine
	go func() {
		for {
			select {
			// 监听取消信号
			// 当cancel()被调用时，ctx.Done()会被关闭，此case会触发
			case <-ctx.Done():
				fmt.Println("Goroutine收到取消信号:", ctx.Err())
				return
			default:
				// 如果没有收到取消信号，继续工作
				fmt.Println("工作中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 主goroutine工作2秒后取消
	time.Sleep(2 * time.Second)
	fmt.Println("发送取消信号")
	cancel() // 调用cancel函数，触发ctx.Done()

	// 等待goroutine退出
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

// ============ 2. 超时控制示例 ============
// timeoutContextDemo 演示如何使用WithTimeout设置操作超时时间
// 关键点：
//  1. WithTimeout会在指定时间后自动调用cancel
//  2. 超时后ctx.Done()会被关闭，ctx.Err()返回context.DeadlineExceeded
//  3. defer cancel()是最佳实践，即使超时也要调用cancel释放资源
func timeoutContextDemo() {
	fmt.Println("=== 超时Context ===")

	// 设置1秒超时
	// WithTimeout等价于WithDeadline(parent, time.Now().Add(timeout))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// 即使超时会自动取消，也应该显式调用cancel以释放资源
	defer cancel()

	ch := make(chan string)

	// 模拟一个耗时2秒的操作（故意超过超时时间）
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	// 使用select同时等待结果和超时信号
	select {
	case result := <-ch:
		// 如果在超时前收到结果，执行这里
		fmt.Println("收到结果:", result)
	case <-ctx.Done():
		// 超时后执行这里
		// ctx.Err()会返回context.DeadlineExceeded错误
		fmt.Println("操作超时:", ctx.Err())
	}
	fmt.Println()
}

// ============ 3. 截止时间示例 ============
// deadlineContextDemo 演示如何使用WithDeadline设置操作的截止时间
// 关键点：
//  1. WithDeadline使用绝对时间（而WithTimeout使用相对时间）
//  2. 可以通过ctx.Deadline()获取截止时间和剩余时间
//  3. 到达截止时间后，ctx.Done()会自动关闭
func deadlineContextDemo() {
	fmt.Println("=== 截止时间Context ===")

	// 设置3秒后的截止时间（绝对时间点）
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 检查剩余时间
	// Deadline()方法返回截止时间和一个布尔值（表示是否设置了截止时间）
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("截止时间: %v, 剩余: %v\n", d, time.Until(d))
	}

	// 模拟一个需要4秒的工作（故意超过截止时间）
	select {
	case <-time.After(4 * time.Second):
		// 如果工作在截止时间前完成
		fmt.Println("工作完成")
	case <-ctx.Done():
		// 超过截止时间后执行这里
		fmt.Println("已超过截止时间:", ctx.Err())
	}
	fmt.Println()
}

// ============ 4. 传递值示例 ============
// contextKey 自定义类型作为Context的key
// 使用自定义类型可以避免不同包之间的key冲突
type contextKey string

const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
)

// valueContextDemo 演示如何使用Context传递请求范围的数据
// 关键点：
//  1. 使用自定义类型作为key，避免冲突
//  2. WithValue返回新的context，不会修改原context（不可变性）
//  3. 只用于传递请求范围的元数据，不要滥用
//  4. Value查找是链式查找，性能为O(n)
func valueContextDemo() {
	fmt.Println("=== Context传递值 ===")

	// 创建携带值的context
	ctx := context.Background()
	// 每次WithValue都返回一个新的context，形成链式结构
	ctx = context.WithValue(ctx, requestIDKey, "req-123")
	ctx = context.WithValue(ctx, userIDKey, "user-456")

	// 将context传递给处理函数
	processRequest(ctx)
	fmt.Println()
}

// processRequest 从context中提取请求相关的值
func processRequest(ctx context.Context) {
	// 使用Value方法获取存储的值
	// 如果key不存在，返回nil
	if reqID := ctx.Value(requestIDKey); reqID != nil {
		fmt.Printf("Request ID: %v\n", reqID)
	}

	if userID := ctx.Value(userIDKey); userID != nil {
		fmt.Printf("User ID: %v\n", userID)
	}
}

// ============ 5. 级联取消（父取消，子也取消）============
// cascadeCancelDemo 演示Context的级联取消机制
// 关键点：
//  1. Context形成树形结构，父context取消时，所有子context也会被取消
//  2. 但子context取消不会影响父context
//  3. 这是实现优雅关闭的关键机制
func cascadeCancelDemo() {
	fmt.Println("=== 级联取消示例 ===")

	// 创建父context
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// 从父context派生出两个子context
	// 子context继承父context的取消信号
	childCtx1, cancel1 := context.WithCancel(parentCtx)
	defer cancel1()

	childCtx2, cancel2 := context.WithCancel(parentCtx)
	defer cancel2()

	// 启动两个worker，分别使用不同的子context
	go worker(childCtx1, "Worker 1")
	go worker(childCtx2, "Worker 2")

	// 让worker工作一段时间
	time.Sleep(1 * time.Second)

	// 取消父context
	// 关键：这会导致所有子context（childCtx1、childCtx2）也被取消
	fmt.Println("取消父context")
	parentCancel()

	// 等待worker退出
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

// worker 模拟一个工作goroutine，持续工作直到收到取消信号
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			// 收到取消信号后退出
			fmt.Printf("%s: 收到取消信号\n", name)
			return
		default:
			// 模拟工作
			fmt.Printf("%s: 工作中...\n", name)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

// ============ 6. 多个goroutine协同工作 ============
// multiWorkerDemo 演示如何使用Context协调多个goroutine
// 关键点：
//  1. 一个context可以控制多个goroutine
//  2. 结合WaitGroup确保所有goroutine都正常退出
//  3. 一次cancel()调用可以同时通知所有worker退出
func multiWorkerDemo() {
	fmt.Println("=== 多worker协同工作 ===")

	// 创建可取消的context
	ctx, cancel := context.WithCancel(context.Background())
	// 使用WaitGroup等待所有worker退出
	var wg sync.WaitGroup

	// 启动3个worker goroutine
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			count := 0 // 记录每个worker处理的任务数
			for {
				select {
				case <-ctx.Done():
					// 收到取消信号，打印统计信息并退出
					fmt.Printf("Worker %d: 退出 (处理了 %d 个任务)\n", id, count)
					return
				default:
					// 模拟处理任务
					count++
					fmt.Printf("Worker %d: 处理任务 #%d\n", id, count)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(i)
	}

	// 让worker工作2秒
	time.Sleep(2 * time.Second)
	fmt.Println("\n发送取消信号给所有worker...")
	// 一次cancel()调用会通知所有监听ctx.Done()的goroutine
	cancel()

	// 等待所有worker优雅退出
	// 这是实现优雅关闭的关键步骤
	wg.Wait()
	fmt.Println("所有worker已退出\n")
}

// ============ 7. Context在Pipeline中的应用 ============
// pipelineDemo 演示在数据流Pipeline中使用Context控制流程
// 关键点：
//  1. Pipeline的每个阶段都接收context，可以响应取消信号
//  2. 超时会导致整个Pipeline停止
//  3. 使用channel连接各个阶段，形成数据流
func pipelineDemo() {
	fmt.Println("=== Pipeline示例 ===")

	// 创建带超时的context，3秒后自动取消整个Pipeline
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Stage 1: 生成数据
	// 返回一个只读channel，用于向下游传递数据
	dataCh := generateData(ctx)

	// Stage 2: 处理数据
	// 接收上游的数据，处理后发送到下游
	processedCh := processData(ctx, dataCh)

	// Stage 3: 消费结果
	// 最终阶段，输出处理结果
	for result := range processedCh {
		fmt.Println("最终结果:", result)
	}

	fmt.Println("Pipeline完成\n")
}

// generateData Pipeline第一阶段：生成数据
// 返回只读channel，符合Go的channel使用惯例
func generateData(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		// 确保channel被关闭，通知下游没有更多数据
		defer close(ch)
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				// 收到取消信号（超时或手动取消），停止生成
				fmt.Println("生成器: 收到取消信号")
				return
			case ch <- i:
				// 发送数据到下游
				fmt.Println("生成器: 生成", i)
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
	return ch
}

// processData Pipeline第二阶段：处理数据
// 从input读取数据，处理后发送到返回的channel
func processData(ctx context.Context, input <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		// 持续从input读取数据
		for data := range input {
			select {
			case <-ctx.Done():
				// 收到取消信号，停止处理
				fmt.Println("处理器: 收到取消信号")
				return
			case ch <- data * 2:
				// 处理数据（这里简单地乘以2）并发送到下游
				fmt.Println("处理器: 处理", data, "->", data*2)
			}
		}
	}()
	return ch
}

// ============ 8. 模拟HTTP请求超时控制 ============
// httpRequestDemo 演示如何为HTTP请求设置超时
// 关键点：
//  1. 实际应用中使用http.NewRequestWithContext()
//  2. 超时可以防止慢请求阻塞程序
//  3. 使用ctx.Err()判断具体的错误类型
func httpRequestDemo() {
	fmt.Println("=== HTTP请求超时控制模拟 ===")

	// 创建2秒超时的context
	// 在真实场景中，这个context会传递给http.NewRequestWithContext()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 模拟HTTP请求
	// 使用缓冲channel避免goroutine泄漏
	result := make(chan string, 1)
	go func() {
		// 模拟一个耗时3秒的慢请求（故意超过超时时间）
		time.Sleep(3 * time.Second)
		result <- "HTTP响应数据"
	}()

	// 等待结果或超时
	select {
	case res := <-result:
		// 在超时前收到响应
		fmt.Println("请求成功:", res)
	case <-ctx.Done():
		// 超时或被取消
		// 检查具体的错误类型
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("❌ 请求超时")
		} else {
			fmt.Println("❌ 请求被取消")
		}
	}
	fmt.Println()
}

// ============ 9. 模拟数据库查询超时 ============
// databaseQueryDemo 演示如何为数据库查询设置超时
// 关键点：
//  1. 实际应用中使用db.QueryContext(ctx, query, args...)
//  2. 避免慢查询拖垮整个系统
//  3. Context是数据库操作的标准超时控制方式
func databaseQueryDemo() {
	fmt.Println("=== 数据库查询超时模拟 ===")

	// 创建1秒超时的context
	// 在真实场景中，这个context会传递给db.QueryContext()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 模拟数据库查询
	result := queryDatabase(ctx, "SELECT * FROM users WHERE active = true")

	if result != nil {
		fmt.Println("查询成功:", result)
	}
	fmt.Println()
}

// queryDatabase 模拟执行数据库查询
// 实际场景中应该使用: db.QueryContext(ctx, query, args...)
func queryDatabase(ctx context.Context, query string) []string {
	// 使用缓冲channel避免goroutine泄漏
	result := make(chan []string, 1)

	go func() {
		// 模拟慢查询（2秒），故意超过超时时间
		// 在真实场景中，这里会执行真正的数据库查询
		time.Sleep(2 * time.Second)
		result <- []string{"user1", "user2", "user3"}
	}()

	// 等待查询结果或超时
	select {
	case data := <-result:
		// 在超时前得到查询结果
		return data
	case <-ctx.Done():
		// 超时或被取消
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("❌ 数据库查询超时")
		}
		return nil
	}
}

// ============ 10. Context错误处理 ============
// contextErrorHandlingDemo 演示如何正确处理Context错误
// 关键点：
//  1. ctx.Err()返回错误类型：nil、Canceled或DeadlineExceeded
//  2. 应该区分不同的错误类型，做出相应处理
//  3. context.Canceled表示手动取消
//  4. context.DeadlineExceeded表示超时
func contextErrorHandlingDemo() {
	fmt.Println("=== Context错误处理 ===")

	// 测试1：超时场景
	ctx1, cancel1 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel1()
	time.Sleep(1 * time.Second) // 故意等待超过超时时间
	handleContextError(ctx1, "超时测试")

	// 测试2：手动取消场景
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2() // 立即取消
	handleContextError(ctx2, "取消测试")

	// 测试3：正常场景（未取消也未超时）
	ctx3 := context.Background()
	handleContextError(ctx3, "正常测试")

	fmt.Println()
}

// handleContextError 处理Context错误并输出相应信息
// 这是处理Context错误的标准模式
func handleContextError(ctx context.Context, label string) {
	// 获取context的错误状态
	err := ctx.Err()

	fmt.Printf("[%s] ", label)
	// 使用switch区分不同的错误类型
	switch err {
	case context.Canceled:
		// 手动调用cancel()导致的取消
		fmt.Println("操作被取消")
	case context.DeadlineExceeded:
		// 超时导致的取消
		fmt.Println("操作超时")
	case nil:
		// 没有错误，context仍然有效
		fmt.Println("Context仍然有效")
	default:
		// 其他未知错误（理论上不应该出现）
		fmt.Println("未知错误:", err)
	}
}

// ============ 11. 实战：任务管理系统 ============

// Task 表示一个任务
// 在实际应用中，可以扩展更多字段，如优先级、依赖关系等
type Task struct {
	ID       int           // 任务ID
	Duration time.Duration // 任务预计执行时间
}

// TaskManager 任务管理器
// 提供任务的添加、执行和管理功能
type TaskManager struct {
	tasks []Task     // 任务列表
	mu    sync.Mutex // 保护tasks的并发安全
}

// AddTask 添加任务到管理器
// 使用互斥锁保证并发安全
func (tm *TaskManager) AddTask(task Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks = append(tm.tasks, task)
}

// ExecuteTask 执行单个任务，支持超时和取消
// 返回error表示任务是否成功完成
func (tm *TaskManager) ExecuteTask(ctx context.Context, task Task) error {
	fmt.Printf("任务 %d: 开始执行 (预计耗时 %v)\n", task.ID, task.Duration)

	// 模拟任务执行
	// 使用select同时监听任务完成和取消信号
	select {
	case <-time.After(task.Duration):
		// 任务正常完成
		fmt.Printf("✅ 任务 %d: 执行完成\n", task.ID)
		return nil
	case <-ctx.Done():
		// 收到取消信号（超时或手动取消）
		fmt.Printf("❌ 任务 %d: 被取消 (%v)\n", task.ID, ctx.Err())
		return ctx.Err()
	}
}

// ExecuteAll 并发执行所有任务
// 使用WaitGroup等待所有任务完成
// 所有任务共享同一个context，可以统一取消
func (tm *TaskManager) ExecuteAll(ctx context.Context) {
	var wg sync.WaitGroup

	// 复制任务列表，避免长时间持有锁
	tm.mu.Lock()
	tasks := tm.tasks
	tm.mu.Unlock()

	// 为每个任务启动一个goroutine
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			// 将context传递给每个任务
			// 当context被取消时，所有任务都会收到信号
			tm.ExecuteTask(ctx, t)
		}(task)
	}

	// 等待所有任务完成（正常完成或被取消）
	wg.Wait()
	fmt.Println("所有任务处理完成")
}

// taskManagerDemo 任务管理系统的完整演示
// 这是一个综合实战示例，展示了：
//  1. 如何组织并发任务管理系统
//  2. 如何使用Context控制多个并发任务
//  3. 如何处理超时和取消
//  4. 如何实现优雅关闭
func taskManagerDemo() {
	fmt.Println("=== 任务管理系统示例 ===")

	// 创建任务管理器
	tm := &TaskManager{}

	// 添加5个随机耗时的任务
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 5; i++ {
		tm.AddTask(Task{
			ID:       i,
			Duration: time.Duration(rand.Intn(3)+1) * time.Second, // 1-3秒随机
		})
	}

	// 创建带超时的context（5秒超时）
	// 如果所有任务在5秒内未完成，会自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 可选：模拟手动取消
	// 取消注释下面的代码，可以看到手动取消的效果
	/*
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("\n⚠️  手动触发取消信号")
			cancel()
		}()
	*/

	// 执行所有任务
	fmt.Println("开始执行任务...")
	tm.ExecuteAll(ctx)

	// 检查context的最终状态，判断任务完成情况
	if ctx.Err() == context.DeadlineExceeded {
		// 超时导致任务未完成
		fmt.Println("\n❌ 任务执行超时")
	} else if ctx.Err() == context.Canceled {
		// 手动取消导致任务未完成
		fmt.Println("\n❌ 任务被手动取消")
	} else {
		// 所有任务在超时前正常完成
		fmt.Println("\n✅ 所有任务正常完成")
	}
	fmt.Println()
}

// ============ 12. Worker Pool with Context ============
// workerPoolDemo 演示带Context控制的Worker Pool模式
// 这是并发编程中最常用的模式之一
// 关键点：
//  1. 固定数量的worker goroutine，避免goroutine爆炸
//  2. 使用channel作为任务队列
//  3. Context用于优雅关闭所有worker
//  4. WaitGroup确保所有worker都退出后再关闭results channel
func workerPoolDemo() {
	fmt.Println("=== Worker Pool示例 ===")

	// 创建可取消的context，用于控制所有worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建任务队列和结果队列
	// 使用缓冲channel可以减少阻塞
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// 启动3个worker goroutine（固定数量）
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go poolWorker(ctx, w, jobs, results, &wg)
	}

	// 发送任务到队列
	go func() {
		for i := 1; i <= 8; i++ {
			jobs <- i
		}
		close(jobs) // 关闭jobs channel表示没有更多任务
	}()

	// 模拟在2秒后取消所有worker（可选）
	// 取消注释可以看到取消效果
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n⚠️  发送取消信号...")
		cancel()
	}()

	// 等待所有worker完成后关闭results channel
	// 这是一个重要的模式：确保发送方关闭channel
	go func() {
		wg.Wait()      // 等待所有worker退出
		close(results) // 关闭results，通知接收方没有更多结果
	}()

	// 持续接收并打印结果
	// 当results被关闭且清空后，for-range会自动退出
	for result := range results {
		fmt.Printf("收到结果: %d\n", result)
	}

	fmt.Println("Worker Pool完成\n")
}

// poolWorker Worker Pool中的单个worker
// 参数：
//   - ctx: 用于接收取消信号
//   - id: worker的唯一标识
//   - jobs: 任务队列（只读）
//   - results: 结果队列（只写）
//   - wg: 用于通知主程序worker已退出
func poolWorker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	// 确保退出时通知WaitGroup
	defer wg.Done()

	for {
		select {
		// 尝试从任务队列接收任务
		case job, ok := <-jobs:
			if !ok {
				// jobs channel已关闭，没有更多任务
				fmt.Printf("Worker %d: 任务队列已关闭，退出\n", id)
				return
			}
			// 处理任务
			fmt.Printf("Worker %d: 处理任务 %d\n", id, job)
			time.Sleep(500 * time.Millisecond) // 模拟任务处理

			// 发送结果，同时监听取消信号
			select {
			case results <- job * 2:
				fmt.Printf("Worker %d: 完成任务 %d -> %d\n", id, job, job*2)
			case <-ctx.Done():
				// 在发送结果时收到取消信号
				fmt.Printf("Worker %d: 收到取消信号，丢弃结果\n", id)
				return
			}

		// 监听取消信号
		case <-ctx.Done():
			fmt.Printf("Worker %d: 收到取消信号，退出\n", id)
			return
		}
	}
}

// ============ 主函数 ============
// main 运行所有Context示例
// 使用说明：
//   - 默认运行所有示例，可以注释掉不需要的
//   - 每个示例都是独立的，可以单独运行
//   - 建议按顺序学习，由浅入深
func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Context上下文控制 - 完整示例集合      ║")
	fmt.Println("╚════════════════════════════════════════╝\n")

	// 提示：可以注释掉不需要运行的示例
	// 建议初学者按顺序逐个运行，理解每个示例后再运行下一个

	// ========== 基础示例 ==========
	// 这4个示例展示Context的基本用法
	cancellableDemo()     // 1. 可取消的Context - 最基础的用法
	timeoutContextDemo()  // 2. 超时控制 - 最常用的场景
	deadlineContextDemo() // 3. 截止时间 - 绝对时间控制
	valueContextDemo()    // 4. 传递值 - 传递请求范围数据

	// // ========== 进阶示例 ==========
	// // 这3个示例展示Context的高级应用
	cascadeCancelDemo() // 5. 级联取消 - 父子Context关系
	multiWorkerDemo()   // 6. 多worker协同 - 统一控制多个goroutine
	pipelineDemo()      // 7. Pipeline应用 - 数据流处理

	// // ========== 实用场景 ==========
	// // 这3个示例展示Context在实际开发中的应用
	httpRequestDemo()          // 8. HTTP请求超时 - Web开发必备
	databaseQueryDemo()        // 9. 数据库查询超时 - 后端开发必备
	contextErrorHandlingDemo() // 10. 错误处理 - 正确处理Context错误

	// // ========== 综合实战 ==========
	// // 这2个示例是完整的实战项目，综合运用前面的所有知识
	taskManagerDemo() // 11. 任务管理系统 - 综合实战示例
	workerPoolDemo()  // 12. Worker Pool - 并发编程经典模式

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║         所有示例执行完成！             ║")
	fmt.Println("╚════════════════════════════════════════╝")
}
