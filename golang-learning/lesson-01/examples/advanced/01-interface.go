package main

import (
	"fmt"
	"math"
)

// 定义接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 实现Shape接口
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 实现Shape接口
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle 实现Shape接口
type Triangle struct {
	A, B, C float64
}

func (t Triangle) Area() float64 {
	// 海伦公式
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

type Animal interface {
	CallName() string
	Say() string
}

type Dog struct {
	Name string
}

func (d Dog) Say() string {
	return "wang"
}

func (d Dog) CallName() string {
	return d.Name
}

type Cat struct {
	Name string
}

func (d Cat) Say() string {
	return "miao"
}

func (d Cat) CallName() string {
	return d.Name
}

// PrintShapeInfo 多态函数
func PrintAnimalInfo(s Animal) {
	fmt.Println(s.CallName(), s.Say())
}

// PrintShapeInfo 多态函数
func PrintShapeInfo(s Shape, name string) {
	fmt.Printf("%s: Area=%.2f, Perimeter=%.2f\n",
		name, s.Area(), s.Perimeter())
}

// Database 接口与多态实现
type Database interface {
	Connect() error
	Query(sql string) ([]string, error)
	Close() error
}

type MySQL struct {
	connection string
}

func (m *MySQL) Connect() error {
	fmt.Printf("连接 MySQL: %s\n", m.connection)
	return nil
}

func (m *MySQL) Query(sql string) ([]string, error) {
	fmt.Printf("MySQL 执行查询: %s\n", sql)
	return []string{"mysql-result-1", "mysql-result-2"}, nil
}

func (m *MySQL) Close() error {
	fmt.Println("关闭 MySQL 连接")
	return nil
}

type PostgreSQL struct {
	connection string
}

func (p *PostgreSQL) Connect() error {
	fmt.Printf("连接 PostgreSQL: %s\n", p.connection)
	return nil
}

func (p *PostgreSQL) Query(sql string) ([]string, error) {
	fmt.Printf("PostgreSQL 执行查询: %s\n", sql)
	return []string{"postgres-result-1", "postgres-result-2"}, nil
}

func (p *PostgreSQL) Close() error {
	fmt.Println("关闭 PostgreSQL 连接")
	return nil
}

func executeQuery(db Database, sql string) {
	if err := db.Connect(); err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer db.Close()

	result, err := db.Query(sql)
	if err != nil {
		fmt.Println("查询失败:", err)
		return
	}
	fmt.Println("查询结果:", result)
}

// 空接口示例
func emptyInterfaceDemo(value interface{}) {
	fmt.Printf("\n值: %v, 类型: %T\n", value, value)

	// 类型断言
	if str, ok := value.(string); ok {
		fmt.Println("是字符串:", str)
	}

	// type switch
	switch v := value.(type) {
	case int:
		fmt.Println("整数:", v)
	case string:
		fmt.Println("字符串:", v)
	case []int:
		fmt.Println("整数切片:", v)
	case map[string]interface{}:
		fmt.Println("map[string]interface{}:", v)
	case map[string]int:
		fmt.Println("map[string]int{}:", v)
	default:
		fmt.Println("未知类型")
	}
}

// 接口组合
type Writer interface {
	Write([]byte) (int, error)
}

type Closer interface {
	Close() error
}

type WriterCloser interface {
	Writer
	Closer
}

// File 实现WriterCloser
type File struct {
	name string
}

func (f *File) Write(data []byte) (int, error) {
	fmt.Printf("写入到 %s: %s\n", f.name, string(data))
	return len(data), nil
}

func (f *File) Close() error {
	fmt.Printf("关闭 %s\n", f.name)
	return nil
}

func main() {
	fmt.Println("=== 接口和多态示例 ===\n")

	// 不同形状
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 5},
		Triangle{A: 3, B: 4, C: 5},
	}

	names := []string{"矩形", "圆形", "三角形"}

	for i, shape := range shapes {
		PrintShapeInfo(shape, names[i])
	}

	animals := []Animal{Dog{Name: "dog"}, Cat{Name: "cat"}}

	for _, animal := range animals {
		PrintAnimalInfo(animal)
	}

	// 空接口
	fmt.Println("\n=== 空接口示例 ===")
	emptyInterfaceDemo(42)
	emptyInterfaceDemo("hello")
	emptyInterfaceDemo([]int{1, 2, 3})
	emptyInterfaceDemo(map[string]int{"a": 1})

	// 接口组合
	fmt.Println("\n=== 接口组合示例 ===")
	file := &File{name: "test.txt"}
	var wc WriterCloser = file
	var wr Writer = file
	var wcr Closer = file

	data := []byte("Hello World")
	wc.Write(data)
	wc.Close()

	wr.Write(data)
	wcr.Close()

	// // 多态数据库示例
	fmt.Println("\n=== 数据库多态示例 ===")
	mysql := &MySQL{connection: "mysql://localhost:3306"}
	postgres := &PostgreSQL{connection: "postgres://localhost:5432"}
	// mock
	executeQuery(mysql, "SELECT * FROM users")
	executeQuery(postgres, "SELECT * FROM orders")
}
