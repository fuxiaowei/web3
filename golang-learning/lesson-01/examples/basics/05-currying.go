package main

import (
	"fmt"
	"strings"
)

// ==================== 什么是柯里化 ====================
// 柯里化（Currying）是一种函数式编程技术：
// 将接受多个参数的函数转换为一系列接受单个参数的函数
// 通过闭包，每个函数都可以访问之前传入的参数

// ==================== 示例1: 简单的柯里化 ====================
// 普通函数：add(a, b int) int
// 柯里化后：curryAdd(a int) func(int) int

// 普通加法函数
func add(a, b int) int {
	return a + b
}

// 柯里化版本：先传入第一个参数，返回一个函数
func curryAdd(x int) func(int) int {
	// 闭包捕获了 x
	return func(y int) int {
		return x + y
	}
}

// ==================== 示例2: 多参数柯里化 ====================
// 普通函数：multiply(a, b, c int) int
// 柯里化后：curryMultiply(a int) func(int) func(int) int

func multiply(a, b, c int) int {
	return a * b * c
}

// 多参数柯里化：逐层返回函数
func curryMultiply(a int) func(int) func(int) int {
	return func(b int) func(int) int {
		return func(c int) int {
			return a * b * c
		}
	}
}

// ==================== 示例3: 实际应用 - 配置化函数 ====================

// 创建一个带有前缀的日志函数
func createLogger(prefix string) func(string) {
	return func(message string) {
		fmt.Printf("[%s] %s\n", prefix, message)
	}
}

// 创建一个带有乘法器的计算函数
func createMultiplier(factor int) func(int) int {
	return func(value int) int {
		return value * factor
	}
}

// ==================== 示例4: 字符串处理柯里化 ====================

// 创建一个字符串连接函数（固定前缀）
func stringWithPrefix(prefix string) func(string) string {
	return func(s string) string {
		return prefix + s
	}
}

// 创建一个字符串重复函数
func repeatString(n int) func(string) string {
	return func(s string) string {
		return strings.Repeat(s, n)
	}
}

// ==================== 示例5: 条件判断柯里化 ====================

// 创建一个比较函数
func greaterThan(threshold int) func(int) bool {
	return func(value int) bool {
		return value > threshold
	}
}

// 创建一个范围检查函数
func inRange(min, max int) func(int) bool {
	return func(value int) bool {
		return value >= min && value <= max
	}
}

// ==================== 示例6: 部分应用（Partial Application） ====================
// 注意：部分应用和柯里化略有不同
// - 柯里化：每次只接受一个参数
// - 部分应用：一次可以传入多个参数

// 部分应用示例
func partialAdd(a, b int) func(int) int {
	sum := a + b // 先计算前两个参数的和
	return func(c int) int {
		return sum + c
	}
}

func main() {
	fmt.Println("========== 柯里化（Currying）示例 ==========\n")

	// ==================== 示例1: 基本柯里化 ====================
	fmt.Println("--- 示例1: 基本柯里化 ---")
	// 普通函数调用
	fmt.Printf("普通函数: add(5, 3) = %d\n", add(5, 3))

	// 柯里化调用
	addFive := curryAdd(5)                             // 先传入第一个参数，返回一个函数 f(x) = 5 + y
	fmt.Printf("柯里化: addFive(3) = %d\n", addFive(3))   // 5 + 3
	fmt.Printf("柯里化: addFive(10) = %d\n", addFive(10)) // 5 + 10
	resultt := curryAdd(3)(5)                          // 3+y   3+5  8
	fmt.Println(resultt)

	// 可以重复使用同一个柯里化函数
	fmt.Printf("柯里化: addFive(100) = %d\n", addFive(100)) // 5 + 100
	fmt.Println()

	// // ==================== 示例2: 多参数柯里化 ====================
	fmt.Println("--- 示例2: 多参数柯里化 ---")
	fmt.Printf("普通函数: multiply(2, 3, 4) = %d\n", multiply(2, 3, 4))

	// 柯里化调用：逐层传入参数
	multiplyBy2 := curryMultiply(2)   // 2 * b * c
	multiplyBy2And3 := multiplyBy2(3) // 2 * 3 * c
	result := multiplyBy2And3(4)      // 2 * 3 * 4
	fmt.Printf("柯里化: curryMultiply(2)(3)(4) = %d\n", result)

	// 也可以链式调用
	result2 := curryMultiply(2)(3)(4)
	fmt.Printf("链式调用: curryMultiply(2)(3)(4) = %d\n", result2)
	fmt.Println()

	// // ==================== 示例3: 实际应用 ====================
	fmt.Println("--- 示例3: 实际应用 - 配置化函数 ---")

	// 创建不同前缀的日志函数
	infoLog := createLogger("INFO")
	warnLog := createLogger("WARN")
	errorLog := createLogger("ERROR")

	infoLog("应用程序启动")
	warnLog("内存使用率较高")
	errorLog("数据库连接失败")
	fmt.Println()

	// 创建不同的乘法器
	double := createMultiplier(2)             //  2 * x
	triple := createMultiplier(3)             //  3 * x
	fmt.Printf("double(5) = %d\n", double(5)) // 2 * 5
	fmt.Printf("triple(5) = %d\n", triple(5)) // 3 * 5
	fmt.Println()

	// // ==================== 示例4: 字符串处理 ====================
	fmt.Println("--- 示例4: 字符串处理柯里化 ---")

	withPrefix := stringWithPrefix("Hello, ") //  hello +  ''
	fmt.Println(withPrefix("World"))
	fmt.Println(withPrefix("Golang"))

	repeat3Times := repeatString(3)
	fmt.Printf("repeat3Times(\"Go\") = %s\n", repeat3Times("Go"))
	fmt.Println()

	// // ==================== 示例5: 条件判断 ====================
	// fmt.Println("--- 示例5: 条件判断柯里化 ---")

	// isGreaterThan10 := greaterThan(10)
	// fmt.Printf("isGreaterThan10(5) = %v\n", isGreaterThan10(5))
	// fmt.Printf("isGreaterThan10(15) = %v\n", isGreaterThan10(15))

	// isInValidRange := inRange(1, 100)
	// fmt.Printf("isInValidRange(50) = %v\n", isInValidRange(50))
	// fmt.Printf("isInValidRange(150) = %v\n", isInValidRange(150))
	// fmt.Println()

	// // ==================== 示例6: 部分应用 ====================
	// fmt.Println("--- 示例6: 部分应用 vs 柯里化 ---")
	// partialAdd10 := partialAdd(5, 5) // 先计算 5+5=10
	// fmt.Printf("partialAdd(5,5)(3) = %d\n", partialAdd10(3))
	// fmt.Println()

	// // ==================== 柯里化的优势 ====================
	// fmt.Println("========== 柯里化的优势 ==========")
	// fmt.Println("1. 函数复用：可以创建预配置的函数")
	// fmt.Println("2. 代码组织：将复杂函数拆分为多个步骤")
	// fmt.Println("3. 灵活性：可以按需组合不同的函数")
	// fmt.Println("4. 可读性：代码更清晰，意图更明确")
}
