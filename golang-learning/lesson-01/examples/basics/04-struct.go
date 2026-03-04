package main

import "fmt"

// 定义结构体
type Person struct {
	Name string
	Age  int
}

// 值接收者方法
func (p Person) GetInfo() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// 值接收者方法（不能修改结构体）
func (p Person) IncrementAgeWrong() {
	p.Age++ // 这不会修改原始结构体
}

// 指针接收者方法（可以修改结构体）
func (p *Person) IncrementAge() {
	p.Age++
}

// 指针接收者方法：修改信息
func (p *Person) ChangeName(name string) {
	p.Name = name
}

// 嵌入结构体
type Employee struct {
	Person // 匿名字段
	ID     string
}

// Employee自己的方法
func (e Employee) GetEmployeeInfo() string {
	return fmt.Sprintf("ID: %s, %s", e.ID, e.GetInfo())
}

// 接口示例
type Speaker interface {
	Speak() string
}

func (p Person) Speak() string {
	return fmt.Sprintf("Hi, I'm %s", p.Name)
}

// 多态示例
func introduce(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	fmt.Println("=== 结构体示例 ===")

	// 初始化结构体
	p1 := Person{
		Name: "Alice",
		Age:  25,
	}
	fmt.Println("p1:", p1)

	// 简短初始化
	p2 := Person{"Bob", 30}
	fmt.Println("p2:", p2)

	// 部分初始化
	p3 := Person{Name: "Charlie"}
	fmt.Println("p3:", p3)

	// 使用值接收者方法
	fmt.Println("\n值接收者方法:")
	fmt.Println(p1.GetInfo())
	fmt.Println(p2.GetInfo())

	// 使用指针接收者方法
	fmt.Println("\n指针接收者方法:")
	fmt.Printf("年龄: %d\n", p1.Age)
	p1.IncrementAge()
	fmt.Printf("增加年龄后: %d\n", p1.Age)

	// 尝试值接收者修改（不会成功）
	p1.IncrementAgeWrong()
	fmt.Printf("使用值接收者后: %d\n", p1.Age)

	// 修改name
	p1.ChangeName("Alice Smith")
	fmt.Println("修改姓名后:", p1.GetInfo())

	// 嵌入结构体
	fmt.Println("\n=== 嵌入结构体 ===")
	employee := Employee{
		Person: Person{Name: "David", Age: 28},
		ID:     "E001",
	}
	fmt.Println("employee:", employee)
	fmt.Println("直接访问嵌入的字段:", employee.Name, employee.Age)
	fmt.Println("调用嵌入结构体的方法:", employee.GetInfo())
	fmt.Println("调用自己的方法:", employee.GetEmployeeInfo())

	// 接口和多态
	fmt.Println("\n=== 接口和多态 ===")
	// introduce(p1)
	// introduce(employee)
}
