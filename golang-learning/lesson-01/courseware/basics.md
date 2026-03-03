# 第1期 - 基础语法部分课件

## 课程目标
- 掌握Go语言的基础语法
- 理解数据类型和控制流
- 学会函数和结构体的使用
- 完成基础编程练习

## 📚 示例代码索引

本课程配套的示例代码位于 `examples/basics/` 目录下：

| 示例文件 | 对应章节 | 主要内容 |
|---------|---------|---------|
| `01-types.go` | 第二部分：基础数据类型 | 基本类型、复合类型示例 |
| `02-control-flow.go` | 第三部分：控制流 | if/else, switch, for, defer, panic/recover |
| `03-function.go` | 第四部分：函数 | 函数定义、多返回值、闭包、柯里化 |
| `04-struct.go` | 第四部分：结构体 | 结构体、方法、嵌入、值/指针接收者 |
| `05-currying.go` | 第四部分：闭包详解 | 柯里化详细示例和应用 |
| `calculator.go` | 第五部分：实践练习 | 计算器实现（综合练习） |

**运行示例代码：**
```bash
# 进入示例目录
cd lesson-01/examples/basics

# 运行指定示例
go run 01-types.go
go run 03-function.go
go run 04-struct.go
# ...
```

---

## 第一部分：环境搭建（5分钟）

### 1. Go安装

**下载地址：**
- 官网：https://golang.org/dl/
- 中国镜像：https://golang.google.cn/dl/

**安装步骤：**
```bash
# macOS
brew install go

# Linux
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到PATH
export PATH=$PATH:/usr/local/go/bin
```

**验证安装：**
```bash
go version
```

### 2. GOPATH vs GOPROXY

**GOPATH（旧方式）：**
- Go 1.11之前必须设置
- 指定工作区位置
- 所有代码放在`$GOPATH/src`下

**Go Modules（推荐）：**
- Go 1.11+引入
- 不需要GOPATH
- 更好的依赖管理

**配置GOPROXY（使用代理加速）：**
```bash
# 七牛云代理（推荐）
go env -w GOPROXY=https://goproxy.cn,direct

# 阿里云代理
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 查看配置
go env | grep GOPROXY
```

### 3. 编辑器配置

**Ubuntu/Linux 文本编辑器：**
```bash
# 安装 vim（推荐）
sudo apt-get update
sudo apt-get install vim

# 使用 nano（更简单，已内置）
nano 文件名
# 保存：Ctrl + O
# 退出：Ctrl + X

# 使用图形界面编辑器
gedit 文件名
```

**VS Code配置（可选）：**
- 安装扩展：Go扩展（ID: golang.go）
- 配置设置：
```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "gofmt",
  "go.lintTool": "golangci-lint"
}
```

---

## 第二部分：基础数据类型

### 1. 基本类型

#### 1.1 布尔类型 (bool)
```go
var isTrue bool = true
var isFalse bool = false

// 布尔类型占用1字节内存
// 只有两个值：true 和 false
```

#### 1.2 整数类型 - 有符号整数
```go
var age int = 25              // int, 平台相关（32位系统32位，64位系统64位）
var count int8 = 100          // int8, -128 到 127 (1字节)
var number int16 = 30000      // int16, -32768 到 32767 (2字节)
var data int32 = 2000000000   // int32, -2^31 到 2^31-1 (4字节)
var big int64 = 9000000000000 // int64, -2^63 到 2^63-1 (8字节)
```

**有符号整数的表示原理：**
- **符号位**：最高位（最左边的bit）用于表示符号
  - `0` 表示正数或零
  - `1` 表示负数
- **范围计算**：n位有符号整数可以表示 2^n 个不同值
  - 负数：-2^(n-1) 到 -1
  - 非负数：0 到 2^(n-1)-1

**int8 示例（8位 = 1字节）：**
```
二进制表示         十进制值
10000000  →  -128 (最小负数)
10000001  →  -127
...
11111111  →    -1 (最大负数)
00000000  →     0 (零)
00000001  →     1 (最小正数)
...
01111111  →   127 (最大正数)
```
总共可以表示：256 个值 = 2^8
- 正数/零：128 个（0 到 127）
- 负数：128 个（-128 到 -1）

#### 1.3 整数类型 - 无符号整数
```go
var index uint = 100          // uint, 平台相关
var count uint8 = 255         // uint8, 0 到 255 (1字节)
var max uint16 = 65535        // uint16, 0 到 65535 (2字节)
var large uint32 = 4294967295 // uint32, 0 到 2^32-1 (4字节)
var huge uint64 = 99999999999 // uint64, 0 到 2^64-1 (8字节)

// 特殊类型别名
var b byte = 65               // byte, uint8的别名，用于存储字符ASCII码
var r rune = '中'             // rune, int32的别名，用于存储Unicode字符
```

**无符号整数的表示原理：**
- **无符号**：没有符号位，所有位都用于表示数值
- **范围**：n位无符号整数可以表示 2^n 个不同值，范围是 0 到 2^n-1
- **优势**：相比相同位数的有符号整数，可以表示更大的正数

**uint8 (byte) 示例（8位 = 1字节）：**
```
二进制表示         十进制值
00000000  →     0 (最小值)
00000001  →     1
...
01111111  →   127
10000000  →   128
...
11111110  →   254
11111111  →   255 (最大值)
```
总共可以表示：256 个值 = 2^8（0 到 255）

**对比总结：**
| 类型 | 位数 | 符号位 | 值域范围 | 总数量 |
|------|------|--------|----------|--------|
| int8 | 8 | 有 | -128 到 127 | 256 |
| uint8 (byte) | 8 | 无 | 0 到 255 | 256 |

**关键理解：如何区分正负数？**

相同位数的 `int8` 和 `uint8` 都是 8 位，通过**类型**区分如何解读这些位：

```go
var signed int8 = -128    // 二进制: 10000000
var unsigned uint8 = 128  // 二进制: 10000000
// 注意：它们的二进制表示完全相同！
```

**区别在于：**
- **int8（有符号）**：最高位是符号位
  - 最高位 = 0 → 正数或零（0 到 127）
  - 最高位 = 1 → 负数（-128 到 -1）
  
- **uint8（无符号）**：所有位都是数值位
  - 所有位按二进制直接计算数值
  - 范围：0 到 255

**编译器如何知道？**
- 通过**类型声明**：当你声明 `var x int8`，编译器就知道应该用有符号方式解读这8位
- 当你声明 `var y uint8`，编译器就知道应该用无符号方式解读这8位
- 因此**类型**是区分正负数的关键！

#### 1.4 浮点类型
```go
var price float32 = 99.99              // float32, 单精度浮点数（约7位十进制精度）
var precise float64 = 3.14159265359    // float64, 双精度浮点数（约15位十进制精度）
var pi = 3.14                          // 默认推断为float64

// 科学计数法
var scientific float64 = 1.23e-4       // 0.000123
```

**重要说明：Go语言没有 `double` 类型**

如果你来自C/C++/Java等语言，可能会寻找 `double` 类型，但在Go中：
- ❌ **没有 `double` 类型**
- ✅ **使用 `float64` 替代 `double`**
- ✅ **`float32` 对应 C/C++ 的 `float`**
- ✅ **`float64` 对应 C/C++ 的 `double`**

**类型对应关系：**

| Go语言 | C/C++ | Java | 说明 |
|--------|-------|------|------|
| `float32` | `float` | `float` | 32位单精度浮点数 |
| `float64` | `double` | `double` | 64位双精度浮点数 |

**默认类型：**
```go
x := 3.14              // 默认为 float64（相当于 double）
var y float32 = 3.14 // 明确指定为 float32

// 在Go中，不带类型的浮点数字面量总是 float64
fmt.Printf("%T\n", 3.14)  // 输出：float64
```

**浮点数的二进制表示（IEEE 754标准）：**

浮点数使用IEEE 754标准表示，将数值分为三个部分：符号位、指数位和尾数位（有效数字位）。

**float32（32位单精度浮点数）：**
```
位分布：1位符号 + 8位指数 + 23位尾数 = 32位

[S][EEEEEEEE][MMMMMMMMMMMMMMMMMMMMMMM]
 ↑     ↑              ↑
符号   指数(8位)      尾数(23位)
```

- **符号位（1位）**：`0`表示正数，`1`表示负数
- **指数位（8位）**：表示2的幂次，范围 -126 到 +127（使用偏移量127，实际值 = 存储值 - 127）
- **尾数位（23位）**：存储有效数字（小数部分）

**示例：1.0 的 float32 表示**
```
二进制：0 01111111 00000000000000000000000
       ↑   ↑         ↑
      符号 指数      尾数
       
符号位 = 0（正数）
指数位 = 01111111 = 127 → 实际指数 = 127 - 127 = 0
尾数位 = 000...000（隐含1.0）
结果 = (-1)^0 × 1.0 × 2^0 = 1.0
```

**float64（64位双精度浮点数）：**
```
位分布：1位符号 + 11位指数 + 52位尾数 = 64位

[S][EEEEEEEEEEE][MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM]
 ↑      ↑                            ↑
符号   指数(11位)                   尾数(52位)
```

- **符号位（1位）**：`0`表示正数，`1`表示负数
- **指数位（11位）**：表示2的幂次，范围 -1022 到 +1023（使用偏移量1023）
- **尾数位（52位）**：存储有效数字（小数部分）

**对比表：**

| 类型 | 总位数 | 符号位 | 指数位 | 尾数位 | 精度 | 范围 |
|------|--------|--------|--------|--------|------|------|
| float32 | 32 | 1 | 8 | 23 | 约7位十进制 | ±3.4×10^38 |
| float64 | 64 | 1 | 11 | 52 | 约15-17位十进制 | ±1.7×10^308 |

**特殊值：**
- **零**：所有位都为0
  - float32: `0 00000000 00000000000000000000000` = +0.0
- **无穷大**：指数全为1，尾数全为0
  - 正无穷：`0 11111111 00000000000000000000000`
  - 负无穷：`1 11111111 00000000000000000000000`
- **NaN（Not a Number）**：指数全为1，尾数非0
  - 例如：`0/0`的结果

**为什么会有精度问题？**
```go
var x float32 = 0.1
var y float32 = 0.2
fmt.Println(x + y)  // 可能输出：0.30000001（不是精确的0.3）

// 原因：某些小数无法精确表示为二进制小数
// 类似 1/3 = 0.333... 在十进制中无法精确表示
```

**注意事项：**
- 浮点数运算可能有精度损失
- 比较浮点数应该使用误差范围，而不是直接 `==` 比较
- 金融计算建议使用 `decimal` 类型（需要第三方库）

#### 1.5 复数类型

**什么是复数？**
 **编程中基本很少用，了解即可**
复数是数学中的概念，由一个**实部**和一个**虚部**组成，形式为：`a + bi`
- **实部（Real Part）**：`a`，是实数部分（真实存在的数字）
- **虚部（Imaginary Part）**：`b`，是虚数部分（带有 `i` 的系数）
- **`i`**：虚数单位，满足 `i² = -1`

**生活中的类比：**
想象你在一个坐标系中：
- 实部：X轴上的位置（左右移动）
- 虚部：Y轴上的位置（上下移动）
- 复数：表示平面上的一个点

**Go语言中的复数类型：**
```go
var c1 complex64 = 1 + 2i              // complex64, 实部和虚部都是float32
var c2 complex128 = 3.14 + 6.28i       // complex128, 实部和虚部都是float64

// 使用complex函数创建
c3 := complex(1.0, 2.0)               // 实部1.0，虚部2.0

// 获取实部和虚部
realPart := real(c1)                   // 获取实部：1
imagPart := imag(c1)                   // 获取虚部：2
```

**复数的基本概念详解：**

**1. 实部（Real Part）**
- 实部是复数的"实数"部分
- 表示实际存在的数值
- 可以用 `real()` 函数获取

**2. 虚部（Imaginary Part）**
- 虚部是复数的"虚数"部分
- 是虚数单位 `i` 的系数
- 可以用 `imag()` 函数获取
- **注意**：`imag()` 返回的是系数，不包含 `i`

**示例说明：**
```go
c := 3 + 4i  // 这是一个复数

fmt.Println(real(c))  // 输出：3（实部）
fmt.Println(imag(c))  // 输出：4（虚部，注意是4，不是4i）

// 完整形式：3 + 4i
// 实部 = 3
// 虚部 = 4（系数）
```

**复数类型对比：**

| 类型 | 内存大小 | 实部类型 | 虚部类型 | 精度 |
|------|----------|----------|----------|------|
| `complex64` | 64位（8字节） | float32 | float32 | 单精度 |
| `complex128` | 128位（16字节） | float64 | float64 | 双精度 |

**复数运算示例：**
```go
// 创建复数
c1 := 3 + 4i        // 实部3，虚部4
c2 := 1 + 2i        // 实部1，虚部2

// 加法：(3+4i) + (1+2i) = 4 + 6i
sum := c1 + c2
fmt.Println(sum)             // 输出：(4+6i)
fmt.Println(real(sum))       // 输出：4
fmt.Println(imag(sum))       // 输出：6

// 乘法：(3+4i) × (1+2i) = 3 + 6i + 4i + 8i²
// = 3 + 10i - 8 = -5 + 10i（因为i² = -1）
product := c1 * c2
fmt.Println(product)         // 输出：(-5+10i)

// 求模（绝对值）
// |3+4i| = √(3² + 4²) = √(9 + 16) = 5
// Go中没有内置函数，需要自己计算
modulus := math.Sqrt(real(c1)*real(c1) + imag(c1)*imag(c1))
fmt.Println(modulus)         // 输出：5
```

**实际应用场景：**
- **信号处理**：音频、图像处理
- **工程计算**：电路分析、机械振动
- **科学计算**：量子力学、电磁学
- **图形学**：旋转、缩放等变换

**注意事项：**
- 在日常编程中，复数使用频率不高
- 主要用在科学计算和特定工程领域
- 如果不涉及这些领域，可以先了解，不需要深入掌握

#### 1.6 字符串类型 (string)
```go
var name string = "Golang"
var greeting = "Hello World"           // 类型推断

// 字符串特点
// - 不可变的字节序列
// - 使用UTF-8编码
// - 可以使用索引访问字节（不是字符）
str := "Hello"
firstByte := str[0]                    // 72 (H的ASCII码)

// 原始字符串字面量
raw := `这是一个
多行
字符串`                                 // 原样保留换行符和特殊字符

// 字符串拼接
str1 := "Hello"
str2 := "World"
combined := str1 + " " + str2          // "Hello World"
```

#### 1.7 指针类型 (pointer)
```go
var p *int                             // 指向int类型的指针
var x int = 10
p = &x                                 // 获取x的地址

// 零值
var ptr *int                           // nil（指针的零值是nil）
```

#### 1.8 基本类型总结

**基础类型一览表：**

| 类型 | 说明 | 值域 | 零值 |
|------|------|------|------|
| bool | 布尔类型 | true, false | false |
| int8 | 8位有符号整数 | -128 到 127 | 0 |
| int16 | 16位有符号整数 | -32,768 到 32,767 | 0 |
| int32 (rune) | 32位有符号整数 | -2^31 到 2^31-1 | 0 |
| int64 | 64位有符号整数 | -2^63 到 2^63-1 | 0 |
| int | 平台相关有符号整数 | -2^31 到 2^31-1 或 -2^63 到 2^63-1 | 0 |
| uint8 (byte) | 8位无符号整数 | 0 到 255 | 0 |
| uint16 | 16位无符号整数 | 0 到 65,535 | 0 |
| uint32 | 32位无符号整数 | 0 到 4,294,967,295 | 0 |
| uint64 | 64位无符号整数 | 0 到 2^64-1 | 0 |
| uint | 平台相关无符号整数 | 0 到 2^32-1 或 2^64-1 | 0 |
| float32 | 32位浮点数 | IEEE 754 单精度 | 0.0 |
| float64 | 64位浮点数 | IEEE 754 双精度 | 0.0 |
| complex64 | 64位复数 | float32 + float32*i | 0+0i |
| complex128 | 128位复数 | float64 + float64*i | 0+0i |
| string | 字符串 | UTF-8字节序列 | "" (空字符串) |
| *T | 指针 | 指向类型T的地址 | nil |

**类型转换：**

Go语言**不允许隐式类型转换**，必须显式转换。这与其他语言（如Java）不同。

**Go语言（必须显式转换）：**
```go
// ❌ Go不允许：隐式类型转换
var x int = 10
// var y float64 = x  // 编译错误！不能隐式转换

// ✅ Go要求：显式类型转换
var x int = 10
var y float64 = float64(x)             // 必须显式转换：int转float64

var a float64 = 3.14
var b int = int(a)                     // float64转int（会截断小数部分）

var str string = "123"
var num, _ = strconv.Atoi(str)         // string转int需要使用strconv包
```

**Java语言（允许隐式转换）示例：**

在Java中，某些类型转换是**隐式**的，编译器会自动进行转换：

```java
// ✅ Java允许：int自动隐式转换为double
int x = 10;
double y = x;          // 自动转换：int → double（无需显式转换）
System.out.println(y); // 输出：10.0

// ✅ Java允许：short隐式转换为int
short s = 100;
int i = s;             // 自动转换：short → int

// ✅ Java允许：float隐式转换为double
float f = 3.14f;
double d = f;          // 自动转换：float → double

// ✅ Java允许：char隐式转换为int
char c = 'A';
int code = c;          // 自动转换：char → int，输出：65
```

**对比总结：**

| 操作 | Go语言 | Java语言 |
|------|--------|----------|
| int → float64 | 必须显式：`float64(x)` | 可以隐式：`double y = x` |
| int → int64 | 必须显式：`int64(x)` | 可以隐式：`long y = x` |
| float32 → float64 | 必须显式：`float64(x)` | 可以隐式：`double y = x` |

**为什么Go不允许隐式转换？**

1. **安全性**：避免意外的精度损失或数据截断
2. **明确性**：代码意图更清晰，减少歧义
3. **可读性**：一眼就能看出发生了类型转换

**示例对比：为什么Go更安全？**

```java
// Java：隐式转换可能导致问题
int count = 1000000;
float price = count;  // 隐式转换：可能丢失精度
// price 可能不是精确的 1000000.0
```

```go
// Go：必须显式转换，提醒你可能有问题
var count int = 1000000
// var price float32 = count  // 编译错误！必须显式转换
var price float32 = float32(count)  // 明确表示：我知道我在转换
```

**Go的类型转换规则：**
- ✅ 同类型转换：`int32(x)` → 改变类型但保持精度
- ⚠️ 不同精度转换：`float64(x)` → `int(x)` 会截断小数
- ❌ 不能转换：不兼容的类型无法转换（需要特殊函数，如`strconv.Atoi`）

### 2. 复合类型

Array, Slice, Map

#### 2.1 数组（Array）

**什么是数组（Array）？**

数组（Array）是Go语言中的**固定长度**的序列数据结构，它由相同类型的元素组成，通过索引访问元素。

**数组的核心特征：**
1. **固定长度**：一旦声明，长度不能改变
2. **值类型**：数组作为值传递时，会复制整个数组的内容
3. **同构性**：数组中的元素必须是相同类型
4. **索引访问**：使用从0开始的整数索引访问元素

**代码示例：**
```go
// 声明数组
var arr [5]int                    // 声明长度为5的整数数组
arr := [5]int{1, 2, 3, 4, 5}     // 初始化的数组
arr2 := [...]int{1, 2, 3}        // 自动推断长度

// 数组特点
// - 长度固定
// - 值类型（复制会复制整个数组）
// - 索引从0开始
```

#### 2.2 切片（Slice）

**什么是切片（Slice）？**

切片（Slice）是Go语言中最常用的数据结构之一，它是**动态数组**的抽象。与固定长度的数组不同，切片可以根据需要自动增长或缩小。

**切片的核心特征：**
1. **动态大小**：长度可以根据需要改变
2. **引用类型**：切片本身是一个轻量级数据结构，包含指针、长度和容量
3. **底层数组**：切片实际上是对底层数组的引用，多个切片可以共享同一个底层数组
4. **自动扩容**：当容量不足时，Go会自动分配更大的底层数组

**切片的内部结构（理解）：**
```go
// 切片在内存中的表示（简化理解）
type slice struct {
    ptr    *T      // 指向底层数组的指针
    len    int     // 当前长度（元素个数）
    cap    int     // 容量（底层数组可以容纳的元素数）
}
```

**代码示例：**
```go
// 声明切片
var slice []int                    // 空切片
slice := []int{1, 2, 3, 4, 5}     // 初始化切片
slice2 := make([]int, 5)          // 使用make创建切片

// 切片操作
slice = append(slice, 6)          // 追加元素
slice = slice[1:3]                // 切片截取
length := len(slice)              // 获取长度
capacity := cap(slice)            // 获取容量

// 切片的本质
// - 对底层数组的引用
// - 包含长度和容量信息
// - 自动扩容机制
```

**切片的自动扩容机制详解：**

切片的容量（capacity）是指底层数组可以容纳的元素个数。当使用`append`向切片添加元素时，如果容量不足，Go会自动扩容。

**扩容触发条件：**
```go
slice := make([]int, 0, 3)  // 长度0，容量3
fmt.Printf("初始: len=%d, cap=%d\n", len(slice), cap(slice))

slice = append(slice, 1)    // len=1, cap=3（未触发扩容）
slice = append(slice, 2)    // len=2, cap=3（未触发扩容）
slice = append(slice, 3)    // len=3, cap=3（未触发扩容）
slice = append(slice, 4)    // len=4, cap=6（触发扩容！）
fmt.Printf("扩容后: len=%d, cap=%d\n", len(slice), cap(slice))
```

**扩容规则（Go 1.18+）：**

实际的扩容策略会根据元素类型和当前容量动态调整：

1. **初始容量和翻倍策略**：
   - 空切片首次添加元素时，根据元素大小分配初始容量（通常为 4 或 8）
   - 小切片（容量 < 256）大致按翻倍策略扩容：
     - 容量 0 → 首次扩容到 4 或 8（取决于元素大小）
     - 容量 4 → 扩容后容量 8
     - 容量 8 → 扩容后容量 16
     - 容量 16 → 扩容后容量 32

2. **大切片（容量 ≥ 256）**：每次增加约 25%（向上取整）
   - 容量 256 → 扩容后容量约为 336
   - 容量 512 → 扩容后容量约为 672

**注意**：具体的扩容策略可能会根据Go版本和元素大小有所调整，但总体原则是：小切片翻倍，大切片按比例增长。

**扩容示例演示：**
```go
func demonstrateSliceGrowth() {
    var s []int
    fmt.Println("开始扩容演示：")
    
    for i := 0; i < 20; i++ {
        oldCap := cap(s)
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

// 实际输出示例：
// 添加元素 0: 长度=1, 容量 0 → 4 (扩容！)
// 添加元素 1: 长度=2, 容量=4 (未扩容)
// 添加元素 2: 长度=3, 容量=4 (未扩容)
// 添加元素 3: 长度=4, 容量=4 (未扩容)
// 添加元素 4: 长度=5, 容量 4 → 8 (扩容！)
// 添加元素 8: 长度=9, 容量 8 → 16 (扩容！)
// 添加元素 16: 长度=17, 容量 16 → 32 (扩容！)
// ...
```

**扩容的内部机制：**

1. **分配新数组**：当容量不足时，Go会分配一个新的更大的底层数组
2. **复制数据**：将旧数组的所有元素复制到新数组
3. **更新切片**：切片指向新的底层数组，长度增加，容量扩大
4. **旧数组回收**：旧的底层数组由垃圾回收器回收

**重要注意事项：**

⚠️ **扩容会创建新的底层数组！**
```go
s1 := []int{1, 2, 3}
s2 := s1              // s2和s1共享底层数组
s2[0] = 100
fmt.Println(s1)       // 输出：[100 2 3]（共享底层数组）

s1 = append(s1, 4)    // 触发扩容，创建新数组
s1[0] = 999
fmt.Println(s1)       // 输出：[999 2 3 4]（新数组）
fmt.Println(s2)       // 输出：[100 2 3]（仍指向旧数组）
```

**性能优化建议：**

1. **预估容量**：如果知道大概需要多少元素，使用`make([]T, 0, capacity)`预设容量
   ```go
   // ❌ 性能差：频繁扩容
   var s []int
   for i := 0; i < 1000; i++ {
       s = append(s, i)  // 可能触发多次扩容
   }
   
   // ✅ 性能好：预设容量
   s := make([]int, 0, 1000)
   for i := 0; i < 1000; i++ {
       s = append(s, i)  // 只扩容一次（或0次）
   }
   ```

2. **一次性追加多个元素**：
   ```go
   s = append(s, 1, 2, 3, 4, 5)  // 一次性追加，可能只扩容一次
   ```

**扩容总结表（int类型示例）：**

| 当前容量 | 追加后触发扩容 | 新容量 |
|---------|---------------|--------|
| 0 | 1个元素 | 4（初始分配） |
| 4 | 5个元素 | 8 |
| 8 | 9个元素 | 16 |
| 16 | 17个元素 | 32 |
| 32 | 33个元素 | 64 |
| 256 | 257个元素 | 约 336 |

**注意**：不同元素类型的初始容量可能不同，扩容策略也会根据元素大小和内存对齐进行调整。

**关键理解：**
- 长度（len）：切片实际包含的元素个数
- 容量（cap）：底层数组可容纳的元素个数
- 扩容是自动的：当`append`导致长度超过容量时自动触发
- 扩容有成本：需要分配新内存和复制数据，应尽量减少扩容次数

#### 2.3 映射（Map）

**什么是映射（Map）？**

映射（Map）是Go语言中的**键值对**（key-value）集合，类似于其他语言中的字典（dictionary）或哈希表（hash table）。它提供了一种快速查找数据的方式，通过唯一的键来访问对应的值。

**映射的核心特征：**
1. **键值对存储**：每个元素都是键（key）和值（value）的配对
2. **引用类型**：映射是引用类型，传递时不会复制所有数据
3. **无序性**：Go语言的映射是无序的，遍历顺序不固定
4. **快速查找**：通过键快速访问值，时间复杂度接近O(1)
5. **键唯一性**：同一个映射中，每个键只能出现一次

**代码示例：**
```go
// 声明映射
var m map[string]int              // 使用前需要初始化
m := make(map[string]int)         // 使用make初始化
m2 := map[string]int{              // 直接初始化
    "apple":  5,
    "banana": 3,
}

// 映射操作
m["apple"] = 10                   // 添加或修改
value := m["apple"]               // 读取
delete(m, "banana")               // 删除
value, ok := m["orange"]          // 检查key是否存在

// 遍历映射
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
```

---

### 3. 指针（Pointer）详解

#### 3.1 什么是指针？

指针是存储另一个变量内存地址的变量。可以理解为"指向另一个变量的箭头"。

```go
// 声明指针
var p *int          // 指向int类型的指针
x := 10
p = &x              // 获取x的地址（& 是取地址运算符）

// 使用指针
fmt.Println(*p)     // 解引用，输出10（* 是解引用运算符）
*p = 20             // 通过指针修改x的值
fmt.Println(x)      // 输出20，因为x的值被修改了
```

#### 3.2 `*` 和 `&` 的区别详解

这是指针操作中最关键的两个运算符，初学者经常混淆。记住它们的区别非常重要！

**`&` (取地址运算符 - Address Of)**
- **作用**：获取变量的内存地址
- **位置**：放在变量前面
- **返回**：指针类型（地址）
- **记忆方法**：`&` 看起来像箭头指向变量，表示"取地址"

```go
x := 10
fmt.Printf("x的值: %d\n", x)           // 输出：x的值: 10
fmt.Printf("x的地址: %p\n", &x)        // 输出：x的地址: 0xc000012260

p := &x  // p是指针，存储x的地址
fmt.Printf("p的值（地址）: %p\n", p)    // 输出：p的值（地址）: 0xc000012260
```

**`*` (解引用运算符 - Dereference)**
- **作用**：通过指针获取或修改指向的值
- **位置**：放在指针前面
- **返回**：指针指向的值
- **记忆方法**：`*` 表示"指向的内容"

```go
x := 10
p := &x

fmt.Printf("p存储的地址: %p\n", p)      // 输出：地址
fmt.Printf("*p获取的值: %d\n", *p)      // 输出：*p获取的值: 10

*p = 20  // 通过*修改指针指向的值
fmt.Printf("x现在的值: %d\n", x)        // 输出：x现在的值: 20
```

**对比记忆表：**

| 运算符 | 名称 | 作用 | 使用位置 | 示例 | 结果 |
|--------|------|------|----------|------|------|
| `&` | 取地址 | 获取变量地址 | 变量前 | `&x` | 地址（指针） |
| `*` | 解引用 | 获取/修改指向的值 | 指针前 | `*p` | 值 |

**完整示例对比：**

```go
x := 10              // 定义一个int变量
fmt.Println("原始变量 x =", x)  // 输出：原始变量 x = 10

// 使用 & 获取地址
p := &x              // p是指针，存储x的地址
fmt.Printf("&x 获取地址: %p\n", &x)    // 输出：地址，如 0xc000012260
fmt.Printf("p 存储的地址: %p\n", p)    // 输出：相同的地址

// 使用 * 解引用获取值
fmt.Printf("*p 获取值: %d\n", *p)      // 输出：*p 获取值: 10

// 使用 * 修改值
*p = 20              // 通过指针修改x的值
fmt.Println("修改后 x =", x)           // 输出：修改后 x = 20
fmt.Printf("*p 现在的值: %d\n", *p)    // 输出：*p 现在的值: 20
```

**关键理解：**
```
变量 x  → 值: 10, 地址: 0xc000012260
  |
  | 使用 & 获取地址
  ↓
指针 p  → 存储: 0xc000012260
  |
  | 使用 * 解引用
  ↓
访问值  → 10
```

**常见错误：**
```go
x := 10
p := &x

// ❌ 错误：不能对值使用解引用
// var y = *x  // 编译错误：x是值，不是指针

// ✅ 正确：对指针使用解引用
var y = *p     // y = 10

// ❌ 错误：不能对指针使用取地址后再解引用
// var z = *&p  // 这会得到指针本身

// ✅ 正确：地址和解引用的组合
var z = *&x    // z = 10（先取地址，再解引用，得到原值）
```

#### 3.3 指针的基本操作

**1. 声明和初始化**
```go
var p *int          // 声明一个指向int的指针（值为nil）
x := 10
p = &x              // 将x的地址赋值给p

// 或者直接声明并初始化
x := 10
p := &x             // 类型推断，p的类型是*int
```

**2. 获取和修改值**
```go
x := 10
p := &x

fmt.Println(p)       // 输出地址，如：0xc000012260
fmt.Println(*p)     // 解引用获取值，输出：10

*p = 20             // 通过指针修改x的值
fmt.Println(x)      // 输出：20（原变量也被修改了）
```

**3. nil指针**
```go
var p *int          // 未初始化的指针值为nil
fmt.Println(p)      // 输出：<nil>

if p == nil {
    fmt.Println("指针是nil")
}

// 注意：解引用nil指针会导致panic
// *p = 10  // 这会panic！
```

**4. 指针作为函数参数**
```go
// 值传递（不会修改原变量）
func doubleByValue(x int) {
    x = x * 2
}

// 指针传递（会修改原变量）
func doubleByPointer(x *int) {
    *x = *x * 2
}

x := 10
doubleByValue(x)
fmt.Println(x)      // 输出：10（没有改变）

doubleByPointer(&x)
fmt.Println(x)      // 输出：20（被修改了）
```

**5. 指针的指针（多级指针）**
```go
x := 10
p := &x             // p是指向int的指针
pp := &p            // pp是指向指针的指针（*int）

fmt.Println(x)      // 10
fmt.Println(*p)     // 10
fmt.Println(**pp)   // 10（需要两次解引用）
```

#### 3.4 指针的常见用途

**1. 修改函数外部变量**
```go
func swap(a, b *int) {
    temp := *a
    *a = *b
    *b = temp
}

x, y := 10, 20
swap(&x, &y)
fmt.Println(x, y)   // 输出：20 10
```

**2. 避免大对象复制**
```go
type LargeStruct struct {
    data [1000000]int
}

// 传值会复制整个结构体（性能差）
func processByValue(s LargeStruct) {
    // ...
}

// 传指针只复制8字节的地址（性能好）
func processByPointer(s *LargeStruct) {
    // ...
}
```

**3. 返回函数内部创建的变量地址**
```go
func newInt() *int {
    x := 10
    return &x       // Go的垃圾回收器会管理这个内存
}

p := newInt()
fmt.Println(*p)     // 输出：10
```

**4. 与结构体结合使用**
```go
type Person struct {
    Name string
    Age  int
}

func updateAge(p *Person, age int) {
    p.Age = age     // 可以使用 . 直接访问，无需 (*p).Age
}

person := Person{Name: "Alice", Age: 25}
updateAge(&person, 30)
fmt.Println(person.Age)  // 输出：30
```

#### 3.5 指针的注意事项

1. **不能对指针进行算术运算**
   ```go
   p := &x
   // p++  // 错误！Go不支持指针算术
   ```

2. **空指针检查**
   ```go
   if p != nil {
       fmt.Println(*p)
   }
   ```

3. **指针与切片的区别**
   - 切片本身就是一个引用类型
   - 传递切片时不需要使用指针
   - 但如果需要修改切片本身（如重新分配），则需要传递指针

**指针的作用总结：**
- ✅ 避免大对象复制，提高性能
- ✅ 修改函数外部变量
- ✅ 实现引用传递语义
- ✅ 返回函数内部创建的变量地址
- ✅ 实现复杂的数据结构和算法

---

## 第三部分：控制流

### 1.1 if/else

```go
// 基本语法
if condition {
    // do something
} else {
    // do something else
}

// 条件表达式不需要括号
if age >= 18 {
    fmt.Println("已成年")
}

// 支持初始化语句
if num := 10; num > 5 {
    fmt.Println("num大于5")
}

// 多条件判断
if score >= 90 {
    fmt.Println("优秀")
} else if score >= 80 {
    fmt.Println("良好")
} else if score >= 60 {
    fmt.Println("及格")
} else {
    fmt.Println("不及格")
}
```

### 1.2 switch

```go
// 基本switch
day := 3
switch day {
case 1:
    fmt.Println("Monday")
case 2:
    fmt.Println("Tuesday")
case 3:
    fmt.Println("Wednesday")
default:
    fmt.Println("Other day")
}

// 多值匹配
switch day {
case 1, 7:
    fmt.Println("周末")
case 2, 3, 4, 5, 6:
    fmt.Println("工作日")
}

// 条件表达式
switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
default:
    fmt.Println("C")
}

// 无fallthrough
// Go的switch默认不穿透，需要使用break显式中断
switch {
case true:
    fmt.Println("true")
    break  // 显式中断
}
```

### 1.3 for循环

```go
// 基本for循环
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// 类似while循环
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}

// 死循环
for {
    // 需要在循环体内使用break跳出
    if condition {
        break
    }
}

// range遍历
arr := []int{1, 2, 3, 4, 5}
for index, value := range arr {
    fmt.Printf("index: %d, value: %d\n", index, value)
}

// 只遍历值
for _, value := range arr {
    fmt.Println(value)
}

// 遍历字符串（按rune）
str := "Hello 世界"
for _, char := range str {
    fmt.Printf("%c ", char)
}

// 遍历映射
m := map[string]int{"a": 1, "b": 2}
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
```

### 1.4 defer语句

`defer` 是Go语言中用于延迟执行函数调用的关键字。理解 `defer` 的执行时机非常重要。

#### defer的延迟执行时机

**核心概念：defer延迟到函数返回之前执行**

但更准确地说，`defer` 语句会在以下三种情况下执行：

1. **函数正常返回前**：执行完 `return` 语句之后，在函数真正返回之前
2. **发生panic时**：即使函数因为 `panic` 崩溃，defer 也会执行
3. **函数执行到最后一行时**：函数自然结束时（没有明确的 return）

**执行顺序：**
```
函数开始
  ↓
执行普通语句
  ↓
遇到 defer（注册到延迟执行栈，但不执行）
  ↓
继续执行其他语句
  ↓
函数准备返回（return/panic/结束）
  ↓
【关键时机】开始执行 defer 语句（按 LIFO 顺序）
  ↓
函数真正返回
```

#### 详细示例说明

**1. defer 在 return 之后执行**

```go
func example() int {
    defer fmt.Println("World")    // 会在 return 之后执行
    fmt.Println("Hello")
    return 42                     // 先执行 return
}
// 输出：
// Hello
// World        ← 在 return 之后才执行

func main() {
    result := example()
    fmt.Println("返回值:", result)
}
// 输出：
// Hello
// World
// 返回值: 42
```

**2. defer 在 panic 之后也会执行**

```go
func panicExample() {
    defer fmt.Println("清理工作")    // 即使 panic 也会执行
    fmt.Println("开始执行")
    panic("发生错误")              // 触发 panic
    fmt.Println("这行不会执行")
}
// 输出：
// 开始执行
// 清理工作      ← panic 后也会执行
// panic: 发生错误
```

**3. defer 执行顺序：LIFO（后进先出，Last In First Out）**

```go
func multipleDefer() {
    fmt.Println("函数开始")
    defer fmt.Println("1号defer")   // 第1个注册，最后执行
    defer fmt.Println("2号defer")   // 第2个注册，倒数第2执行
    defer fmt.Println("3号defer")   // 第3个注册，第1个执行
    fmt.Println("函数结束")
}
// 输出：
// 函数开始
// 函数结束
// 3号defer      ← 最后注册的，最先执行
// 2号defer
// 1号defer      ← 最先注册的，最后执行
```

**4. defer 捕获变量值的时机（立即计算）**

`defer` 语句中的**参数值会在 defer 语句执行时立即计算并捕获**，而不是在 defer 真正执行时计算。

```go
func deferValue() {
    i := 0
    defer fmt.Println("defer1:", i)    // i的值是0（立即捕获）
    i++
    defer fmt.Println("defer2:", i)    // i的值是1（立即捕获）
    i++
    fmt.Println("函数内:", i)          // 输出：函数内: 2
    return
}
// 输出：
// 函数内: 2
// defer2: 1      ← defer2 最后注册，最先执行（但值是1）
// defer1: 0      ← defer1 先注册，最后执行（值是0）
```

**如果需要捕获最终值，使用闭包：**

```go
func deferClosure() {
    i := 0
    defer func() {
        fmt.Println("defer闭包:", i)    // 捕获变量本身，而不是值
    }()
    i++
    i++
    fmt.Println("函数内:", i)           // 输出：函数内: 2
    return
}
// 输出：
// 函数内: 2
// defer闭包: 2   ← 闭包捕获的是变量的引用，所以是最终值
```

#### defer的常见用途

**1. 释放资源（最常见）**

```go
func readFile() error {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer file.Close()  // 确保无论函数如何返回，文件都会被关闭
    
    // 处理文件...
    // 即使中间发生错误或panic，文件也会被关闭
    return nil
}
```

**2. 解锁互斥锁**

```go
func safeOperation() {
    mu.Lock()
    defer mu.Unlock()  // 确保锁一定会被释放
    
    // 执行需要加锁的操作
}
```

**3. 恢复panic（配合recover使用）**

```go
func safeOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("捕获panic: %v\n", r)
        }
    }()
    
    // 可能触发panic的代码
    panic("发生错误")
}
```

**4. 性能测量**

```go
func slowOperation() {
    defer func() {
        fmt.Println("操作完成")
    }()
    
    start := time.Now()
    defer func() {
        fmt.Printf("耗时: %v\n", time.Since(start))
    }()
    
    // 执行耗时操作
}
```

#### defer的执行时机总结

| 情况 | defer 是否执行 | 执行顺序 |
|------|---------------|---------|
| 函数正常返回 | ✅ 是 | return 语句之后 |
| 函数自然结束 | ✅ 是 | 最后一行语句之后 |
| 发生 panic | ✅ 是 | panic 之后 |
| 多个 defer | ✅ 全部执行 | LIFO（后进先出） |

**重要理解：**
- `defer` 不是"异步"执行，而是在**函数返回流程中的特定时机**执行
- `defer` 语句**注册**到延迟执行栈时，就会计算参数值（立即计算值）
- 如果需要捕获变量的最终值，使用**闭包**而不是直接传递参数
- `defer` 保证了资源清理代码一定会执行，即使发生 panic

### 1.5 panic和recover

```go
// panic触发panic
func riskyOperation() {
    panic("something went wrong")
}

// panic会导致程序崩溃
// 输出：
// panic: something went wrong
// goroutine 1 [running]:
// main.riskyOperation()
// ...

// recover捕获panic
func safeOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("捕获panic: %v\n", r)
        }
    }()
    
    panic("发生错误")
}

// recover只能在defer中使用
func main() {
    safeOperation()  // 输出：捕获panic: 发生错误
    fmt.Println("程序继续执行")
}
```

---

## 第四部分：函数和结构体

> 💡 **示例代码：** `examples/basics/03-function.go`  
> 运行示例：`cd examples/basics && go run 03-function.go`

### 1.1 函数定义

```go
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
    return  // 自动返回命名的返回值
}

// 可变参数
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// 使用
result := sum(1, 2, 3, 4, 5)
```

**完整示例：** 查看 `examples/basics/03-function.go` 了解更多函数用法。

### 1.2 闭包（Closure）

#### 什么是闭包？

**闭包（Closure）** 是一个函数以及与其相关的引用环境组合而成的实体。更简单地说，闭包就是一个函数，它可以访问并"记住"其外部作用域中的变量，即使外部函数已经执行完毕。

**闭包的核心特征：**
1. **函数嵌套**：闭包是定义在函数内部的函数
2. **捕获外部变量**：内部函数可以访问外部函数的局部变量
3. **变量持久化**：即使外部函数执行完毕，内部函数仍然可以访问这些变量
4. **每次调用独立**：每次调用外部函数都会创建一个新的闭包实例

#### 闭包的工作原理

**没有闭包的普通函数：**
```go
func normalFunction() {
    count := 0
    count++
    fmt.Println(count)  // 每次调用都是1
}

normalFunction()  // 输出：1
normalFunction()  // 输出：1
// 每次调用时，count 都会重新创建并初始化为0
```

**使用闭包：**
```go
// 外部函数
func makeCounter() func() int {
    count := 0  // 这个变量被内部函数"捕获"了
    
    // 内部函数（闭包）
    return func() int {
        count++      // 访问并修改外部变量
        return count
    }
}

// 使用闭包
counter := makeCounter()     // 创建闭包实例
fmt.Println(counter())       // 1 - count从0变成1
fmt.Println(counter())       // 2 - count继续增加
fmt.Println(counter())       // 3 - count被"记住"了

// 创建另一个独立的闭包实例
counter2 := makeCounter()    // 新的闭包，新的count
fmt.Println(counter2())      // 1 - 独立的计数器
fmt.Println(counter())       // 4 - 第一个计数器继续
```

**关键理解：**
```
makeCounter() 被调用时：
  1. 创建局部变量 count = 0
  2. 创建匿名函数（闭包）
  3. 闭包"捕获"了 count 变量的引用
  4. 返回这个闭包
  
即使 makeCounter() 执行完毕：
  - count 变量不会被销毁
  - 闭包仍然持有 count 的引用
  - 每次调用闭包，都可以访问和修改 count
```

#### 闭包的内存模型（简化理解）

```
栈内存                      堆内存
┌──────────────┐          ┌──────────────┐
│ counter      │ ──────→  │ 闭包函数     │
│ (函数引用)   │          │   ↓          │
└──────────────┘          │ count = 3    │
                          │ (被捕获的变量)│
                          └──────────────┘

- counter 是栈上的变量，指向堆上的闭包
- 闭包函数和它捕获的 count 变量都在堆上
- 只要 counter 还在使用，count 就不会被回收
```

#### 闭包示例：计数器

```go
// 闭包：函数可以访问外部变量
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 使用
counter := makeCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
```

#### 闭包的高级用法：柯里化（Currying）

**什么是柯里化？**
柯里化是一种将接受多个参数的函数转换为一系列接受单个参数的函数的技术。

```go
// 普通函数：接受两个参数
func addNormal(x, y int) int {
    return x + y
}

// 柯里化版本：先传入一个参数，返回一个接受第二个参数的函数
func add(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// 使用柯里化
addFive := add(5)        // 创建一个"加5"的函数
result := addFive(3)     // 8（5 + 3）
result2 := addFive(10)   // 15（5 + 10）

// 也可以链式调用
result3 := add(5)(3)     // 8
```

#### 闭包的实际应用场景

**1. 函数工厂（创建配置化的函数）**
```go
// 创建带有特定前缀的日志函数
func createLogger(prefix string) func(string) {
    return func(message string) {
        fmt.Printf("[%s] %s\n", prefix, message)
    }
}

infoLog := createLogger("INFO")
errorLog := createLogger("ERROR")

infoLog("应用启动")     // 输出: [INFO] 应用启动
errorLog("连接失败")    // 输出: [ERROR] 连接失败
```

**2. 封装私有状态**
```go
func createAccount(initialBalance int) (deposit, withdraw, getBalance func(int) int) {
    balance := initialBalance  // 私有变量，外部无法直接访问
    
    deposit = func(amount int) int {
        balance += amount
        return balance
    }
    
    withdraw = func(amount int) int {
        if amount <= balance {
            balance -= amount
        }
        return balance
    }
    
    getBalance = func(_ int) int {
        return balance
    }
    
    return
}

// 使用
deposit, withdraw, getBalance := createAccount(100)
deposit(50)                    // 余额: 150
withdraw(30)                   // 余额: 120
fmt.Println(getBalance(0))     // 输出: 120
// balance 变量无法从外部直接访问，实现了数据封装
```

**3. 延迟执行和惰性求值**
```go
func fibonacci() func() int {
    a, b := 0, 1
    return func() int {
        result := a
        a, b = b, a+b
        return result
    }
}

fib := fibonacci()
fmt.Println(fib())  // 0
fmt.Println(fib())  // 1
fmt.Println(fib())  // 1
fmt.Println(fib())  // 2
fmt.Println(fib())  // 3
fmt.Println(fib())  // 5
```

**4. 回调函数和事件处理**
```go
func createHandler(name string) func() {
    count := 0
    return func() {
        count++
        fmt.Printf("%s 被调用了 %d 次\n", name, count)
    }
}

buttonClick := createHandler("按钮点击")
buttonClick()  // 按钮点击 被调用了 1 次
buttonClick()  // 按钮点击 被调用了 2 次
```

#### 闭包的注意事项

**1. 闭包会延长变量的生命周期**
```go
func createClosures() []func() {
    var funcs []func()
    for i := 0; i < 3; i++ {
        // ❌ 错误：所有闭包共享同一个 i
        funcs = append(funcs, func() {
            fmt.Println(i)
        })
    }
    return funcs
}

closures := createClosures()
for _, f := range closures {
    f()  // 输出：3, 3, 3（不是0, 1, 2）
}

// ✅ 正确：为每个闭包创建独立的变量
func createClosuresCorrect() []func() {
    var funcs []func()
    for i := 0; i < 3; i++ {
        i := i  // 创建新变量（遮蔽外部的i）
        funcs = append(funcs, func() {
            fmt.Println(i)
        })
    }
    return funcs
}
```

**2. 避免闭包导致的内存泄漏**
```go
// ⚠️ 注意：如果不再需要闭包，应该释放引用
var globalClosure func()

func setupClosure() {
    largeData := make([]byte, 1000000)  // 1MB数据
    globalClosure = func() {
        fmt.Println(len(largeData))  // 闭包捕获了 largeData
    }
    // 即使 setupClosure 执行完毕，largeData 也不会被回收
}

// 解决方案：不再需要时，清空引用
globalClosure = nil  // 释放闭包，largeData 才能被回收
```

#### 闭包 vs 结构体方法

闭包和结构体方法在某些场景下可以互换使用：

```go
// 使用闭包
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 使用结构体
type Counter struct {
    count int
}

func (c *Counter) Increment() int {
    c.count++
    return c.count
}

// 两者都能实现相同的功能
closure := makeCounter()
fmt.Println(closure())  // 1

counter := &Counter{}
fmt.Println(counter.Increment())  // 1
```

**选择建议：**
- **闭包**：简单的状态封装、函数式编程风格、临时使用
- **结构体**：复杂的数据结构、需要多个方法、面向对象风格

#### 闭包总结

**优点：**
- ✅ 数据封装：创建私有变量，外部无法直接访问
- ✅ 函数工厂：根据参数创建定制化的函数
- ✅ 代码简洁：不需要定义额外的结构体
- ✅ 函数式编程：支持高阶函数和柯里化

**缺点：**
- ⚠️ 内存占用：闭包会延长变量的生命周期
- ⚠️ 调试困难：变量作用域不如结构体直观
- ⚠️ 容易出错：循环中创建闭包时要特别注意

**关键要点：**
1. 闭包 = 函数 + 它所捕获的环境变量
2. 闭包可以"记住"外部变量的值，即使外部函数已经返回
3. 每次调用外部函数都会创建一个新的、独立的闭包实例
4. 闭包是实现函数式编程的重要工具

**完整示例：**
- 闭包基础示例：`examples/basics/03-function.go`
- 柯里化详细示例：`examples/basics/05-currying.go`

### 1.3 结构体和方法

> 💡 **示例代码：** `examples/basics/04-struct.go`  
> 运行示例：`cd examples/basics && go run 04-struct.go`

```go
// 定义结构体
type Person struct {
    Name string
    Age  int
}

// 初始化结构体
p1 := Person{
    Name: "Alice",
    Age:  25,
}

// 简短初始化
p2 := Person{"Bob", 30}

// 部分初始化
p3 := Person{Name: "Charlie"}

// 定义方法（值接收者）
func (p Person) GetInfo() string {
    return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// 定义方法（指针接收者，可以修改结构体）
func (p *Person) IncrementAge() {
    p.Age++
}

// 使用
person := Person{Name: "Alice", Age: 25}
fmt.Println(person.GetInfo())
person.IncrementAge()
fmt.Println(person.Age)

// 嵌入结构体（类似继承）
type Employee struct {
	Person        // 匿名字段
	ID     string
}

employee := Employee{
	Person: Person{Name: "David", Age: 28},
	ID:     "E001",
}
fmt.Println(employee.Name)  // 可以直接访问嵌入的字段
```

#### 1.3.1 值接收者 vs 指针接收者

这是结构体方法中非常重要的概念，需要深入理解。

**值接收者（Value Receiver）：**
```go
// 值接收者：接收结构体的副本
func (p Person) GetInfo() string {
    return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// ❌ 值接收者无法修改原始结构体
func (p Person) IncrementAgeWrong() {
    p.Age++  // 只修改了副本，不会影响原始结构体
}
```

**指针接收者（Pointer Receiver）：**
```go
// 指针接收者：接收结构体的指针
func (p *Person) IncrementAge() {
    p.Age++  // 可以修改原始结构体
}

func (p *Person) ChangeName(name string) {
    p.Name = name  // 可以修改原始结构体
}
```

**使用对比：**
```go
person := Person{Name: "Alice", Age: 25}

// 值接收者
fmt.Println(person.GetInfo())      // 正常工作，只读取数据

// 值接收者尝试修改（不会成功）
person.IncrementAgeWrong()
fmt.Println(person.Age)            // 还是 25，没有改变

// 指针接收者修改（会成功）
person.IncrementAge()
fmt.Println(person.Age)            // 变成 26，修改成功

person.ChangeName("Alice Smith")
fmt.Println(person.Name)           // "Alice Smith"，修改成功
```

**选择建议：**
| 使用场景 | 推荐接收者类型 | 原因 |
|---------|--------------|------|
| 只读取数据 | 值接收者 | 安全，不会意外修改 |
| 需要修改结构体 | 指针接收者 | 必须使用指针才能修改 |
| 大型结构体 | 指针接收者 | 避免复制大量数据 |
| 小型结构体（只读） | 值接收者 | 更简单直观 |

**完整示例：** 查看 `examples/basics/04-struct.go` 了解值接收者和指针接收者的区别。

---

## 第五部分：实践练习（5分钟）

> 💡 **实践示例：** `examples/basics/calculator.go`  
> 运行示例：`cd examples/basics && go run calculator.go`

### 练习：简单计算器程序

**需求：**
- 实现加减乘除运算
- 支持多个数字运算
- 返回结果和错误信息

**学习建议：**
1. 先自己尝试实现
2. 遇到困难时查看参考实现
3. 理解代码后，尝试添加新功能

**参考实现：**
- 查看 `examples/basics/calculator.go`
- 该示例综合运用了：函数、结构体、方法、错误处理等知识点

**扩展练习：**
- 添加更多运算符（如求余、幂运算）
- 支持括号运算
- 实现运算历史记录功能

---

## 课后作业

1. **基础练习：**
   - 实现一个学生管理系统
- 定义 `Student` 结构体
- 实现增删改查功能
- 使用切片存储学生列表

2. **进阶练习：**
   - 实现一个简单的银行账户系统
- 定义 `Account` 结构体
- 实现存款、取款、查询余额功能
- 处理账户余额不足的情况

---

## 总结

通过本部分学习，你应该掌握：

1. ✅ Go语言的基本类型和复合类型
2. ✅ 指针的使用方法
3. ✅ 控制流的各种语法
4. ✅ defer、panic、recover的用法
5. ✅ 函数的定义和多种返回值
6. ✅ 闭包的概念和用法
7. ✅ 结构体和方法定义

## 🔍 示例代码总结

学习完本部分内容后，建议按顺序运行以下示例代码加深理解：

### 基础示例（必学）
```bash
cd examples/basics

# 1. 数据类型示例
go run 01-types.go

# 2. 控制流示例
go run 02-control-flow.go

# 3. 函数示例（包含闭包）
go run 03-function.go

# 4. 结构体示例（包含方法）
go run 04-struct.go
```

### 进阶示例（可选）
```bash
# 5. 柯里化详细示例
go run 05-currying.go

# 6. 综合练习：计算器
go run calculator.go
```

## 📖 学习路径建议

**第一步：理解概念**
1. 阅读本课件的每个章节
2. 理解示例代码的原理

**第二步：动手实践**
1. 运行 `examples/basics/` 中的所有示例
2. 修改示例代码，观察结果变化
3. 尝试自己编写类似的代码

**第三步：完成练习**
1. 实现课后作业
2. 参考 `calculator.go` 的实现思路
3. 尝试扩展更多功能

**第四步：巩固提高**
1. 理解闭包和柯里化的应用场景
2. 掌握值接收者和指针接收者的区别
3. 练习错误处理的最佳实践

## 🎓 课后作业提交清单

完成以下任务后，你就掌握了本课程的核心内容：

- [ ] 运行并理解所有 `examples/basics/` 中的示例代码
- [ ] 实现学生管理系统（使用结构体、切片、方法）
- [ ] 实现银行账户系统（使用指针、错误处理）
- [ ] 编写至少 3 个使用闭包的实际例子
- [ ] 理解并能解释值接收者和指针接收者的区别
- [ ] 能够处理 panic 并使用 defer 清理资源

## 📌 重点知识点回顾

### 必须掌握
- ✅ **基本类型**：int, float64, string, bool
- ✅ **切片**：动态数组，make, append, len, cap
- ✅ **映射**：map, make, delete, 检查键是否存在
- ✅ **指针**：&取地址, *解引用, nil指针
- ✅ **函数**：多返回值, error处理
- ✅ **结构体**：定义、初始化、方法
- ✅ **defer**：资源清理、执行顺序（LIFO）

### 理解掌握
- ✅ **闭包**：捕获外部变量、函数工厂
- ✅ **接收者**：值接收者 vs 指针接收者
- ✅ **panic/recover**：异常处理机制
- ✅ **类型转换**：显式转换、strconv包

### 扩展了解
- 📘 **柯里化**：函数式编程技术
- 📘 **类型别名**：type关键字
- 📘 **嵌入结构体**：组合优于继承

**下一部分：** 我们将学习接口、并发编程和更多高级特性！

> 💡 **提示：** 在进入下一部分之前，请确保你已经理解并能够独立编写本部分的所有知识点。
> 进阶部分会大量使用这些基础知识，打好基础非常重要！
