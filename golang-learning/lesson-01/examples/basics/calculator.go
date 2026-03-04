package main

import (
	"errors"
	"fmt"
)

// Calculator 结构体
type Calculator struct {
	history []string
}

// NewCalculator 创建新的计算器
func NewCalculator() *Calculator {
	return &Calculator{
		history: make([]string, 0),
	}
}

// Add 加法
func (c *Calculator) Add(a, b float64) (float64, error) {
	result := a + b
	c.addToHistory(fmt.Sprintf("%.2f + %.2f = %.2f", a, b, result))
	return result, nil
}

// Subtract 减法
func (c *Calculator) Subtract(a, b float64) (float64, error) {
	result := a - b
	c.addToHistory(fmt.Sprintf("%.2f - %.2f = %.2f", a, b, result))
	return result, nil
}

// Multiply 乘法
func (c *Calculator) Multiply(a, b float64) (float64, error) {
	result := a * b
	c.addToHistory(fmt.Sprintf("%.2f × %.2f = %.2f", a, b, result))
	return result, nil
}

// Divide 除法
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	result := a / b
	c.addToHistory(fmt.Sprintf("%.2f ÷ %.2f = %.2f", a, b, result))
	return result, nil
}

// Sum 多个数字相加
func (c *Calculator) Sum(numbers ...float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("至少需要输入一个数字")
	}

	var total float64
	for _, num := range numbers {
		total += num
	}

	c.addToHistory(fmt.Sprintf("Sum of %d numbers = %.2f", len(numbers), total))
	return total, nil
}

// Average 求平均值
func (c *Calculator) Average(numbers ...float64) (float64, error) {
	sum, err := c.Sum(numbers...)
	if err != nil {
		return 0, err
	}

	average := sum / float64(len(numbers))
	c.addToHistory(fmt.Sprintf("Average of %d numbers = %.2f", len(numbers), average))
	return average, nil
}

// addToHistory 添加到历史记录
func (c *Calculator) addToHistory(entry string) {
	c.history = append(c.history, entry)
}

// GetHistory 获取历史记录
func (c *Calculator) GetHistory() []string {
	return c.history
}

// ClearHistory 清空历史记录
func (c *Calculator) ClearHistory() {
	c.history = c.history[:0]
}

func main() {
	calc := NewCalculator()

	fmt.Println("=== 简单计算器示例 ===\n")

	// 基本运算
	fmt.Println("基本运算:")
	result, _ := calc.Add(10, 20)
	fmt.Printf("10 + 20 = %.2f\n", result)

	result, _ = calc.Subtract(30, 15)
	fmt.Printf("30 - 15 = %.2f\n", result)

	result, _ = calc.Multiply(5, 6)
	fmt.Printf("5 × 6 = %.2f\n", result)

	result, _ = calc.Divide(100, 4)
	fmt.Printf("100 ÷ 4 = %.2f\n", result)

	// 除法错误处理
	result, err := calc.Divide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	// 多个数字运算
	fmt.Println("\n多个数字运算:")
	result, _ = calc.Sum(1, 2, 3, 4, 5)
	fmt.Printf("Sum(1,2,3,4,5) = %.2f\n", result)

	result, _ = calc.Average(10, 20, 30, 40, 50)
	fmt.Printf("Average(10,20,30,40,50) = %.2f\n", result)

	// 显示历史记录
	fmt.Println("\n计算历史:")
	for i, entry := range calc.GetHistory() {
		fmt.Printf("  %d. %s\n", i+1, entry)
	}

	// 清空历史
	calc.ClearHistory()
	fmt.Println("\n历史已清空")
	fmt.Printf("历史记录数: %d\n", len(calc.GetHistory()))
}
