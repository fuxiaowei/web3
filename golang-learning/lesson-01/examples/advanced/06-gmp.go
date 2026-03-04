package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// gmpDemo 演示通过调整 GOMAXPROCS 与 GODEBUG 查看 GMP 调度信息。
func gmpDemo() {
	fmt.Println("=== GMP 调度示例 ===")
	fmt.Printf("逻辑CPU数量: %d\n", runtime.NumCPU())

	// 	如果 n > 0：将 GOMAXPROCS 设置为 n，并返回修改前的旧值
	//  如果 n = 0：不修改 GOMAXPROCS，只是查询并返回当前值
	prev := runtime.GOMAXPROCS(0)
	fmt.Printf("当前 GOMAXPROCS: %d\n", prev)

	// 通过设置 GOMAXPROCS 观察调度变化
	runtime.GOMAXPROCS(2)
	fmt.Printf("调整后的 GOMAXPROCS: %d\n\n", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	taskCount := 8

	wg.Add(taskCount)
	start := time.Now()

	for i := 0; i < taskCount; i++ {
		go func(id int) {
			defer wg.Done()

			// 模拟 CPU 密集型工作
			sum := 0
			for j := 0; j < 50_000_000; j++ {
				sum += j % (id + 1)
			}
			fmt.Printf("Goroutine %d 完成，结果=%d\n", id, sum)
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("\n所有任务完成，耗时：%s\n", elapsed)
	fmt.Println("提示：运行时可配合命令 `GODEBUG=schedtrace=1000,scheddetail=1 go run 06-gmp.go` 观察调度日志。")

	// 恢复原始的 GOMAXPROCS，避免影响其他程序
	runtime.GOMAXPROCS(prev)
}

func main() {
	gmpDemo()
}
