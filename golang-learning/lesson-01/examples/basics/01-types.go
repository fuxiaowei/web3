package main

import (
	"fmt"
	"unsafe"
)

// 基础类型示例
func basicTypesDemo() {
	fmt.Println("=== 基础类型示例 ===")

	// 布尔类型
	var isTrue bool = true
	fmt.Printf("bool: %t\n", isTrue)

	// 整数类型 - 有符号
	var age int = 25
	fmt.Printf("int: %d\n", age)

	// 不同长度的整数
	var count int8 = 127                      // 8位整数，范围: -128 到 127
	var number int16 = 30000                  // 16位整数，范围: -32768 到 32767
	var data int32 = 2000000000               // 32位整数，范围: -2^31 到 2^31-1
	var bigNumber int64 = 9223372036854775807 // 64位整数，范围: -2^63 到 2^63-1
	fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n", count, number, data, bigNumber)

	// 无符号整数
	var u8 uint8 = 255                    // 0 到 255
	var u16 uint16 = 65535                // 0 到 65535
	var u32 uint32 = 4294967295           // 0 到 2^32-1
	var u64 uint64 = 18446744073709551615 // 0  to 2^64-1
	fmt.Printf("uint8: %d, uint16: %d, uint32: %d, uint64: %d\n", u8, u16, u32, u64)

	// 类型别名
	var b byte = 65  // byte 是 uint8 的别名
	var r rune = '中' // rune 是 int32 的别名
	fmt.Printf("byte: %d (%c), rune: %d (%c)\n", b, b, r, r)

	// 显示类型占用的内存大小
	fmt.Printf("\n类型大小:\n")
	fmt.Printf("int8 size: %d bytes\n", unsafe.Sizeof(count))
	fmt.Printf("int16 size: %d bytes\n", unsafe.Sizeof(number))
	fmt.Printf("int32 size: %d bytes\n", unsafe.Sizeof(data))
	fmt.Printf("int64 size: %d bytes\n", unsafe.Sizeof(bigNumber))
	fmt.Printf("bool size: %d bytes\n", unsafe.Sizeof(isTrue))

	// 浮点型
	var price float32 = 99.99
	var pi float64 = 3.14159265359
	fmt.Printf("\nfloat32: %.2f, float64: %.10f\n", price, pi)
	fmt.Printf("float32 size: %d bytes, float64 size: %d bytes\n", unsafe.Sizeof(price), unsafe.Sizeof(pi))

	// 字符串
	var name string = "Golang"
	var greeting = "Hello World" // 类型推断
	fmt.Printf("string: %s\n", name)
	fmt.Printf("类型推断: %s\n", greeting)

	// 字符串操作
	str := "Hello 世界"
	fmt.Printf("字符串长度: %d (字节数)\n", len(str)) // 注意：len返回字节数，不是字符数
	firstByte := str[0]
	fmt.Printf("第一个字节: %d ('%c')\n", firstByte, firstByte)

	// 原始字符串字面量
	raw := `这是
一个多行
字符串`
	fmt.Println("\n原始字符串:")
	fmt.Println(raw)

	// 字符（使用rune）
	var char rune = '中'
	fmt.Printf("rune: %c\n", char)

	// 复数类型
	var c1 complex64 = 1 + 2i
	var c2 complex128 = 3.14 + 6.28i
	c3 := complex(5.0, 10.0) // 使用complex函数创建

	fmt.Println("\n复数类型:")
	fmt.Printf("complex64: %v, 实部: %.1f, 虚部: %.1f\n", c1, real(c1), imag(c1))
	fmt.Printf("complex128: %.2f, 实部: %.2f, 虚部: %.2f\n", c2, real(c2), imag(c2))
	fmt.Printf("使用complex函数: %v\n", c3)
	fmt.Printf("complex64 size: %d bytes, complex128 size: %d bytes\n", unsafe.Sizeof(c1), unsafe.Sizeof(c2))

	// 演示：为什么 int8 的范围是 -128 到 127？
	fmt.Printf("\n计算: 2的8次方 = %d\n", 1<<8) // 使用位运算计算 2^8
	fmt.Printf("int8 可以表示 %d 个不同的值\n", 1<<8)
	fmt.Printf("int8 范围: %d 到 %d\n", -1<<7, 1<<7-1)

	// 演示int8的边界值
	fmt.Println("\nint8 边界值演示:")
	var minInt8 int8 = -128
	var maxInt8 int8 = 127
	fmt.Printf("最小值: %d\n", minInt8)
	fmt.Printf("最大值: %d\n", maxInt8)

	// 注意：如果超出范围会溢出
	// var overflow int8 = 128 // 这行代码会编译错误：溢出
	// var underflow int8 = -129 // 这行代码会编译错误：溢出

	// 演示：相同二进制在不同类型下的不同解读
	fmt.Println("\n=== 类型区分正负的关键演示 ===")

	// 这两个值的二进制表示完全相同！
	var signed int8 = -128
	var unsigned uint8 = 128

	fmt.Printf("有符号 int8 的 -128：\n")
	fmt.Printf("  十进制: %d\n", signed)
	fmt.Printf("  二进制: %08b\n", uint8(signed))

	fmt.Printf("\n无符号 uint8 的 128：\n")
	fmt.Printf("  十进制: %d\n", unsigned)
	fmt.Printf("  二进制: %08b\n", unsigned)

	fmt.Println("\n关键理解：")
	fmt.Println("  • 二进制表示相同：都是 10000000")
	fmt.Println("  • 类型决定了如何解读这个二进制")
	fmt.Println("  • int8 的最高位被看作符号位（1表示负数）")
	fmt.Println("  • uint8 的所有位都被看作数值（没有符号位）")
	fmt.Println("  • 因此编译器通过类型声明知道应该用哪种方式解读")
}

// 数组示例
func arrayDemo() {
	fmt.Println("\n=== 数组示例 ===")

	// 声明数组
	var arr [5]int
	fmt.Printf("声明但不初始化: %v\n", arr)

	// 初始化数组
	arr = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("初始化: %v\n", arr)
	fmt.Println(arr[0])
	fmt.Println(arr[1])
	fmt.Println(arr[2])

	// 自动推断长度
	arr2 := [...]int{1, 2, 3}
	fmt.Printf("自动长度: %v, 长度: %d\n", arr2, len(arr2))

}

// 切片示例
func sliceDemo() {
	fmt.Println("\n=== 切片示例 ===")

	// 声明切片, 跟声明一个数组的区别，就是 声明数组要有初始化的长度，而切片不需要
	var slice []int
	fmt.Printf("空切片: %v, nil: %t\n", slice, slice == nil)

	// 使用make创建切片
	slice = make([]int, 3, 5) // 长度3，容量5 什么叫容量5，就是初始化的slice，也就是这个动态数组的当前最大容量是5个，如果超过会自动扩容
	fmt.Printf("make创建: %v, 长度: %d, 容量: %d\n", slice, len(slice), cap(slice))

	// // 直接初始化
	slice = []int{1, 2, 3, 4, 5}
	fmt.Printf("初始化: %v\n", slice)

	// // 追加元素
	slice = append(slice, 6)
	fmt.Printf("追加后: %v, 容量: %d\n", slice, cap(slice))

	// // 切片截取
	subSlice := slice[1:3]
	fmt.Printf("切片[1:3]: %v\n", subSlice)

	// // 切片共享底层数组， 切片后的subSlice，和slice共享同一个底层数组，所以修改subSlice会影响slice
	subSlice[0] = 999
	fmt.Printf("修改subSlice后: %v\n", slice)
}

// 映射示例
func mapDemo() {
	fmt.Println("\n=== 映射示例 ===")

	// 声明映射
	var m map[string]int
	fmt.Printf("声明但不初始化: %v, nil: %t\n", m, m == nil)

	// 使用make初始化
	m = make(map[string]int)
	m["apple"] = 5
	m["banana"] = 3
	fmt.Printf("添加后: %v\n", m)

	// // // 直接初始化
	m2 := map[string]int{
		"apple":  10,
		"banana": 5,
		"orange": 8,
	}
	fmt.Printf("直接初始化: %v\n", m2)

	// // // 读取值
	value := m2["apple"]
	fmt.Printf("apple的值: %d\n", value)
	value2 := m2["orange"]
	fmt.Printf("orange的值: %d\n", value2)

	m2["orange"] = 100
	fmt.Printf("orange的值: %d\n", m2["orange"])

	// // 检查key是否存在
	value, ok := m2["grape"]
	fmt.Printf("grape存在: %t, 值: %d\n", ok, value)

	valueA, okA := m2["apple"]
	fmt.Printf("apple存在: %t, 值: %d\n", okA, valueA)

	// delete(m2, "banana")
	// fmt.Printf("删除后: %v\n", m2)

	// // 遍历映射
	fmt.Println("遍历映射:")
	for key, value := range m2 {
		fmt.Printf("  %s: %d\n", key, value)
	}
}

// 指针示例
func pointerDemo() {
	fmt.Println("\n=== 指针示例 ===")

	x := 10
	fmt.Printf("x的值为: %d\n", x)

	// 获取地址
	p := &x
	fmt.Printf("x的地址: %p\n", p)
	fmt.Printf("指针的值: %d\n", *p)

	m := p
	fmt.Printf("m: %p\n", m)
	fmt.Printf("m指向的值: %d\n", *m)

	// // 通过指针修改值
	*p = 20
	fmt.Printf("修改后x的值为: %d\n", x)

	// // 指针作为函数参数
	increment(&x)
	fmt.Printf("函数修改后x的值为: %d\n", x)

	// // nil指针
	var ptr *int
	fmt.Printf("nil指针: %v\n", ptr)
	// 解引用nil指针会导致panic
	// fmt.Println(*ptr) // 这行会panic

	// // ===== 值传递 vs 指针传递 =====
	fmt.Println("\n=== 值传递 vs 指针传递 ===")

	// 值传递示例
	a := 10
	fmt.Printf("调用前 a = %d\n", a)
	valuePass(a)
	fmt.Printf("调用后 a = %d (值未改变)\n", a)

	// 指针传递示例
	b := 10
	fmt.Printf("\n调用前 b = %d\n", b)
	pointerPass(&b)
	fmt.Printf("调用后 b = %d (值已改变)\n", b)

	// // 详细说明
	fmt.Println("\n关键区别：")
	fmt.Println("1. 值传递：函数接收的是值的副本，修改副本不影响原值")
	fmt.Println("2. 指针传递：函数接收的是地址，通过地址可以直接修改原值")
}

func increment(p *int) {
	*p++ // 修改指针指向的值
}

// 值传递：函数接收的是值的副本
func valuePass(num int) {
	fmt.Printf("  函数内接收到的值: %d\n", num)
	num = 100 // 修改副本，不影响原值
	fmt.Printf("  函数内修改后: %d\n", num)
}

// 指针传递：函数接收的是地址
func pointerPass(num *int) {
	fmt.Printf("  函数内接收到的地址: %p\n", num)
	fmt.Printf("  函数内接收到的值: %d\n", *num)
	*num = 100 // 通过指针修改原值
	fmt.Printf("  函数内修改后: %d\n", *num)
}

func main() {
	// basicTypesDemo()
	// arrayDemo()
	// sliceDemo()
	// mapDemo()
	pointerDemo()

	// demonstrateSliceGrowth()
}

func demonstrateSliceGrowth() {
	var s []int
	fmt.Println("开始扩容演示：")

	for i := 0; i < 20; i++ {
		oldCap := cap(s)
		fmt.Print(oldCap)
		s = append(s, i)
		newCap := cap(s)

		if newCap != oldCap {
			fmt.Printf("添加元素 %d: 长度=%d, 容量 %d → %d (扩容！)\n",
				i, len(s), oldCap, newCap)
		} else {
			fmt.Printf("添加元素 %d: 长度=%d, 容量=%d (未扩容)\n",
				i, len(s), cap(s))
		}
	}
}
