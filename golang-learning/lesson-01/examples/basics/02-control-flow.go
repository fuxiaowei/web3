package main

import (
	"fmt"
	"os"
)

// if/else示例
func ifDemo() {
	fmt.Println("=== if/else示例 ===")

	age := 20
	if age >= 18 {
		fmt.Println("已成年")
	} else {
		fmt.Println("未成年")
	}

	// // 支持初始化语句
	if num := 10; num > 5 {
		fmt.Printf("num(%d)大于5\n", num)
	}

	// // 多条件判断
	scores := []int{95, 85, 70, 50}
	for index, score := range scores {
		fmt.Printf("索引%d: %d\n", index, score)
		if score >= 90 {
			fmt.Printf("分数%d: 优秀\n", score)
		} else if score >= 80 {
			fmt.Printf("分数%d: 良好\n", score)
		} else if score >= 60 {
			fmt.Printf("分数%d: 及格\n", score)
		} else {
			fmt.Printf("分数%d: 不及格\n", score)
		}
	}
}

// switch示例
func switchDemo() {
	fmt.Println("\n=== switch示例 ===")

	// 基本switch
	days := []int{1, 3, 5, 7, 10}
	for index, day := range days {
		fmt.Printf("索引%d: %d\n", index, day)
		switch day {
		case 1:
			fmt.Println("Monday")
		case 2:
			fmt.Println("Tuesday")
		case 3:
			fmt.Println("Wednesday")
		case 4:
			fmt.Println("Thursday")
		case 5:
			fmt.Println("Friday")
		case 6, 7:
			fmt.Println("周末")
		default:
			fmt.Printf("Day %d: 未知\n", day)
		}
	}

	// 条件表达式
	score := 85
	switch {
	case score >= 90:
		fmt.Println("等级: A")
	case score >= 80:
		fmt.Println("等级: B")
	case score >= 70:
		fmt.Println("等级: C")
	case score >= 60:
		fmt.Println("等级: D")
	default:
		fmt.Println("等级: E")
	}
}

// switch fallthrough 示例
func switchFallthroughDemo() {
	fmt.Println("\n=== switch fallthrough 示例 ===")

	// 1. 默认行为：不穿透（即使没有break）
	fmt.Println("1. 默认行为：不穿透")
	num := 1
	switch num {
	case 1:
		fmt.Println("  匹配到 case 1")
		// 没有 break，但也不会执行 case 2
	case 2:
		fmt.Println("  匹配到 case 2")
	default:
		fmt.Println("  默认分支")
	}
	// 输出：只输出 "匹配到 case 1"，不会执行 case 2

	// 2. 即使添加 break 也是多余的（显式中断）
	fmt.Println("\n2. 显式 break（虽然多余，但可以增加代码可读性）")
	num = 1
	switch num {
	case 1:
		fmt.Println("  匹配到 case 1")
		break // 显式中断，但实际不需要，因为默认就不穿透
	case 2:
		fmt.Println("  匹配到 case 2")
	default:
		fmt.Println("  默认分支")
	}

	// 3. 使用 fallthrough 实现穿透
	fmt.Println("\n3. 使用 fallthrough 实现穿透")
	num = 1
	switch num {
	case 1:
		fmt.Println("  匹配到 case 1")
		fallthrough // 继续执行下一个 case（穿透）
	case 2:
		fmt.Println("  执行 case 2（因为 fallthrough）")
		fallthrough
	case 3:
		fmt.Println("  匹配到 case 3")
	default:
		fmt.Println("  默认分支")
	}
	// 输出：会输出 "匹配到 case 1" 和 "执行 case 2（因为 fallthrough）"

	// 4. fallthrough 连续穿透
	fmt.Println("\n4. fallthrough 连续穿透")
	num = 1
	switch num {
	case 1:
		fmt.Println("  匹配到 case 1")
		fallthrough // 穿透到 case 2
	case 2:
		fmt.Println("  执行 case 2（因为 fallthrough）")
		fallthrough // 继续穿透到 case 3
	case 3:
		fmt.Println("  执行 case 3（因为 fallthrough）")
		// 这里没有 fallthrough，所以不会执行 default
	default:
		fmt.Println("  默认分支")
	}
	// 输出：会输出 case 1, 2, 3，但不会输出 default

	// 5. 实际应用场景：分级判断（显示所有满足的等级）
	fmt.Println("\n5. 实际应用场景：分级判断（显示所有满足的等级）")
	scores := []int{95, 85, 75, 65}
	for _, score := range scores {
		fmt.Printf("  分数 %d: ", score)
		first := true
		switch {
		case score >= 90:
			fmt.Print("优秀")
			first = false
			fallthrough // 优秀也属于良好和及格
		case score >= 100:
			if !first {
				fmt.Print("、良好")
			} else {
				fmt.Print("良好")
			}
			first = false
			fallthrough // 良好也属于及格
		case score >= 60:
			if !first {
				fmt.Print("、及格")
			} else {
				fmt.Print("及格")
			}
			// 这里不 fallthrough，所以不会输出"不及格"
		default:
			fmt.Print("不及格")
		}
		fmt.Println()
	}

	// 6. 对比：不使用 fallthrough 的正常分级判断
	fmt.Println("\n6. 对比：不使用 fallthrough 的正常分级判断")
	scores2 := []int{95, 85, 75, 55}
	for _, score := range scores2 {
		fmt.Printf("  分数 %d: ", score)
		switch {
		case score >= 90:
			fmt.Println("优秀")
		case score >= 80:
			fmt.Println("良好")
		case score >= 60:
			fmt.Println("及格")
		default:
			fmt.Println("不及格")
		}
	}
}

// for循环示例
func forLoopDemo() {
	fmt.Println("\n=== for循环示例 ===")

	// 基本for循环
	fmt.Println("基本for循环:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %d ", i)
	}
	fmt.Println()

	// // 类似while循环
	fmt.Println("类似while循环:")
	i := 0
	for i < 5 {
		fmt.Printf("  %d ", i)
		i++
	}
	fmt.Println()

	// // 死循环（需要break）
	fmt.Println("循环直到条件满足:")
	i = 0
	for {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d ", i)
		i++
	}
	fmt.Println()

	// // range遍历数组/切片
	fmt.Println("range遍历数组:")
	arr := []int{10, 20, 30, 40, 50}
	for index, value := range arr {
		fmt.Printf("  索引[%d] = %d\n", index, value)
	}

	// // 只遍历值
	fmt.Println("只遍历值:")
	for _, value := range arr {
		fmt.Printf("  %d ", value)
	}
	fmt.Println()

	// // 遍历字符串
	fmt.Println("遍历字符串:")
	str := "Hello 世界"
	for i, char := range str {
		fmt.Printf("  [%d] %c\n", i, char)
	}
}

// defer示例
func deferDemo() {
	fmt.Println("\n=== defer示例 ===")

	// 基本defer
	// fmt.Println("1. 基本defer（延迟到函数返回前）:")
	// defer fmt.Println("World")
	// fmt.Println("Hello")

	// 多个defer执行顺序（LIFO）
	// fmt.Println("\n2. 多个defer执行顺序（LIFO，后进先出）:")
	// fmt.Println("函数开始")
	// defer fmt.Println("1号defer - 最先注册，最后执行")
	// defer fmt.Println("2号defer - 中间注册")
	// defer fmt.Println("3号defer - 最后注册，最先执行")
	// fmt.Println("函数结束")

	// // defer在return之后执行
	// fmt.Println("\n3. defer在return之后执行:")
	// fmt.Println(returnWithDefer())

	// // defer捕获变量值的时机
	// fmt.Println("\n4. defer捕获变量值的时机（立即计算）:")
	// deferValueDemo()

	// // defer闭包捕获最终值
	// fmt.Println("\n5. defer闭包捕获最终值:")
	// deferClosureDemo()

	// // defer在panic之后也会执行
	// fmt.Println("\n6. defer在panic之后也会执行:")
	// deferWithPanic()

	// // defer用于资源清理
	fmt.Println("\n7. defer用于资源清理:")
	if err := readFile("02-control-flow.go"); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}
	if err := readFile("nonexistent.txt"); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}
}

// 演示defer在return之后执行
func returnWithDefer() int {
	defer fmt.Println("  defer: 在return之后执行")
	fmt.Println("  return: 先准备返回值")
	return 42
}

// 演示defer捕获变量值的时机
func deferValueDemo() {
	i := 0
	defer fmt.Println("  defer1: i =", i) // 立即捕获i=0
	i++
	defer fmt.Println("  defer2: i =", i) // 立即捕获i=1
	i++
	fmt.Println("  函数内: i =", i) // 输出最终值2
	// defer执行顺序：defer2, defer1（LIFO）
	// defer2输出: i = 1
	// defer1输出: i = 0
}

// 演示defer闭包捕获最终值
func deferClosureDemo() {
	i := 0
	i++
	i++
	defer func() {
		fmt.Println("  defer闭包: i =", i) // 捕获变量引用，输出最终值2
	}()
	fmt.Println("  函数内: i =", i) // 输出: 2
	// defer执行时，闭包中捕获的是i的引用，所以输出最终值2
}

// 演示defer在panic之后也会执行
func deferWithPanic() {
	defer func() {
		fmt.Println("  清理工作: defer在panic之后执行")
	}()
	fmt.Println("  开始执行")
	fmt.Println("  即将panic...")
	// 注意：这里不真正panic，否则会终止程序
	// 实际应用中可以在defer中使用recover捕获panic
	fmt.Println("  (演示结束，实际panic会执行defer)")
}

func readFile(filename string) error {
	defer func() {
		fmt.Printf("  清理资源: %s\n", filename)
	}()

	fmt.Printf("  准备打开: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	defer fmt.Printf("  关闭文件: %s\n", filename)

	fmt.Printf("  成功打开: %s\n", filename)
	return nil
}

// panic和recover示例
func panicRecoverDemo() {
	fmt.Println("\n=== panic和recover示例 ===")

	fmt.Println("触发panic但使用recover捕获:")
	safeOperation()

	fmt.Println("\n正常执行panic:")
	riskyOperation()
}

func safeOperation() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  捕获panic: %v\n", r)
			fmt.Println("  程序继续执行")
		}
	}()

	fmt.Println("  即将触发panic")
	panic("发生错误")
	fmt.Println("  这行不会执行")
}

func riskyOperation() {
	fmt.Println("这会导致程序崩溃")
	panic("致命错误")
	fmt.Println("这行不会执行")
}

func main() {
	// ifDemo()
	// switchDemo()
	// switchFallthroughDemo() // 演示 fallthrough 行为
	// forLoopDemo()
	// deferDemo()
	panicRecoverDemo()
	fmt.Println("\n程序结束")
}
