package main

import (
	"errors"
	"fmt"
)

// 基本函数
func add(a int, b int) int {
	return a + b
}

// 类型简写
func multiply(a, b int) int {
	return a * b
}

// 多返回值
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

// 命名返回值
func calculate(a, b int) (sum, product int) {
	sum = a + b
	product = a * b
	return // 自动返回命名的返回值 return  sum, product
}

// 可变参数
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// 闭包示例
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 闭包示例：柯里化
func curryAdd(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// 高阶函数：函数作为参数
func applyOperation(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

func main() {
	fmt.Println("=== 函数示例 ===")

	// 基本函数
	fmt.Println("add(10, 20):", add(10, 20))
	fmt.Println("multiply(10, 20):", multiply(10, 20))

	// 多返回值
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("divide(10, 2) = %d\n", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("错误:", err)
	}

	// 命名返回值
	s, p := calculate(10, 20)
	fmt.Printf("calculate(10, 20): sum=%d, product=%d\n", s, p)

	// 可变参数
	fmt.Println("sum(1, 2, 3, 4, 5):", sum(1, 2, 3, 4, 5))
	fmt.Println("sum(1, 2, 3):", sum(1, 2, 3))
	fmt.Println("sum():", sum())

	// 闭包
	fmt.Println("\n=== 闭包示例 ===")
	counter := makeCounter()
	fmt.Println("counter:", counter())
	fmt.Println("counter:", counter())
	fmt.Println("counter:", counter())

	counter2 := makeCounter()
	fmt.Println("counter2:", counter2())

	// // 柯里化
	addFive := curryAdd(5)
	fmt.Println("addFive(3):", addFive(3))
	fmt.Println("addFive(10):", addFive(10))

	// // 高阶函数
	// fmt.Println("\n=== 高阶函数 ===")
	// fmt.Println("applyOperation(10, 20, add):", applyOperation(10, 20, add))
	// fmt.Println("applyOperation(10, 20, multiply):", applyOperation(10, 20, multiply))
	// fmt.Println("applyOperation(100, 5, func(a, b int) int { return a / b }):",
	// 	applyOperation(100, 5, func(a, b int) int { return a / b }))
}
