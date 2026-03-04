package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Person 演示 JSON 编解码
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// jsonDemo 展示标准库 encoding/json 的常用操作。
func jsonDemo() {
	fmt.Println("=== JSON 编解码 ===")

	// 准备要编码的结构体数据。
	p := Person{Name: "Alice", Age: 30}

	// 将结构体编码为 JSON。
	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println("编码结果:", string(data))

	// 也可以使用缩进格式输出，便于阅读。
	pretty, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("缩进格式:", string(pretty))

	var p2 Person
	if err := json.Unmarshal(data, &p2); err != nil {
		panic(err)
	}
	fmt.Println("解码结果:", p2)

	// 解码未知字段时，可以借助 map[string]any 捕获动态字段。
	dynamicJSON := `{"name":"Bob","age":25,"tags":["golang","stdlib"]}`
	var dynamic map[string]any
	if err := json.Unmarshal([]byte(dynamicJSON), &dynamic); err != nil {
		panic(err)
	}
	fmt.Println("动态解码:", dynamic)
	fmt.Println()
}

func fileDemo() {
	fmt.Println("=== 文件读写 ===")

	// runtime.Caller(0) 返回当前函数所在的源文件路径，从而可以定位到与该 Go 文件同级的目录。
	// 这样无论从哪里运行程序，都能在源代码所在目录进行读写，方便查看生成的示例文件。
	_, filename, _, ok := runtime.Caller(0)
	fmt.Println("filename:", filename)
	if !ok {
		panic("无法获取当前文件路径")
	}
	dir := filepath.Dir(filename)

	// 如果需要将文件写入其它目录，只要通过配置覆盖 dir 即可，这里示例读取环境变量 FILE_DEMO_DIR。
	// if customDir, ok := os.LookupEnv("FILE_DEMO_DIR"); ok && strings.TrimSpace(customDir) != "" {
	// 	dir = customDir
	// 	fmt.Println("使用自定义目录:", dir)
	// }

	inputPath := filepath.Join(dir, "go-advanced-input.txt")
	outputPath := filepath.Join(dir, "go-advanced-output.txt")

	// 使用 os.WriteFile 快速写入一个小文件。
	if err := os.WriteFile(inputPath, []byte("Hello, File!"), 0o644); err != nil {
		panic(err)
	}

	// 使用 ReadFile 读取全部内容。
	content, err := os.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("inputPath:", inputPath)
	fmt.Println("读取内容:", string(content))

	// 使用 bufio.Writer 按行写入更多数据。
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	lines := []string{
		"第一行: Hello, Output!",
		"第二行: 来自 bufio.Writer",
	}
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			panic(err)
		}
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}
	fmt.Println("写入文件:", outputPath)

	// 使用 bufio.Scanner 按行读取刚才写入的文件。
	file, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("按行读取输出文件内容:")
	for scanner.Scan() {
		fmt.Println("  >", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func timeDemo() {
	fmt.Println("=== 时间处理 ===")

	now := time.Now()
	fmt.Println("当前时间:", now)
	fmt.Println("格式化:", now.Format("2006-01-02 15:04:05"))

	// 将字符串解析为 time.Time。
	// 第一个参数是布局(layout)，必须写成 Go 独有的“参考时间” 2006-01-02 15:04:05 的格式，用来描述输入字符串的结构。
	// 第二个参数是需要被解析的具体时间字符串，返回值 t 是解析得到的 time.Time，err 表示解析过程中是否出错。
	t, err := time.Parse("2006-01-02", "2024-01-01")
	if err != nil {
		panic(err)
	}
	fmt.Println("解析的时间:", t)

	loc, err := time.LoadLocation("Asia/Shanghai")
	// 第一个参数同样是layout  这个layout是用来解析输入的不是用来格式化输出的
	// 第二个参数是要格式化的时间字符串
	ti, err := time.ParseInLocation("20060102", "20240101", loc)
	if err != nil {
		panic(err)
	}
	fmt.Println("解析的时间（本地时区）:", ti)

	//  参数是layout format函数的layout是用来格式化输出的
	fmt.Println(ti.Format("2006-01-02"))

	//总结一下，time.Parse函数用来字符串转时间类型，  time.Format是用来转换时间输出格式

	// 计算时间差，得到的是一个 duration。
	// t1, err := time.Parse("2006-01-01 14:00:00", "2025-11-17 14:00:00")
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-11-17 14:00:00", loc)
	if err != nil {
		panic(err)
	}
	duration := time.Since(t1)
	fmt.Println("距离 2025-11-17 14:00:00 已过:", duration)

	// 将时间转换到其他时区。
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	fmt.Println("洛杉矶时区现在:", now.In(location))

	// // 使用 ticker 进行简单的定时任务。
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()

	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		<-ticker.C
		// <-ticker.C 从 ticker 的 channel 中取出下一次触发的时间，若尚未触发会在此处阻塞等待。
		fmt.Printf("Ticker 第 %d 次触发\n", i)
	}
	fmt.Println("Ticker 演示结束")
	fmt.Println()
}

// hashDemo 展示 crypto 包的哈希算法。
func hashDemo() {
	fmt.Println("=== 哈希计算 ===")

	data := "Hello World"

	// 1. MD5 哈希（128位，不推荐用于安全场景）
	fmt.Println("1. MD5 哈希:")
	md5Hash := md5.New()
	md5Hash.Write([]byte(data))
	fmt.Printf("   MD5(\"%s\") = %x\n", data, md5Hash.Sum(nil))

	// 2. SHA1 哈希（160位，已被认为不够安全）
	fmt.Println("\n2. SHA1 哈希:")
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(data))
	fmt.Printf("   SHA1(\"%s\") = %x\n", data, sha1Hash.Sum(nil))

	// 3. SHA256 哈希（256位，常用且安全）
	fmt.Println("\n3. SHA256 哈希:")
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(data))
	fmt.Printf("   SHA256(\"%s\") = %x\n", data, sha256Hash.Sum(nil))

	// 4. SHA512 哈希（512位，更安全但更长）
	fmt.Println("\n4. SHA512 哈希:")
	sha512Hash := sha512.New()
	sha512Hash.Write([]byte(data))
	fmt.Printf("   SHA512(\"%s\") = %x\n", data, sha512Hash.Sum(nil))

	// 5. 使用 hex 编码展示不同的输出格式
	fmt.Println("\n5. 不同编码格式:")
	sha256Hash2 := sha256.Sum256([]byte(data))
	fmt.Printf("   十六进制: %x\n", sha256Hash2)
	fmt.Printf("   十六进制(大写): %X\n", sha256Hash2)
	fmt.Printf("   hex.EncodeToString: %s\n", hex.EncodeToString(sha256Hash2[:]))

	// 6. 演示相同数据产生相同哈希
	fmt.Println("\n6. 哈希一致性验证:")
	data1 := "password123"
	data2 := "password123"
	data3 := "password124" // 只有一个字符不同

	hash1 := sha256.Sum256([]byte(data1))
	hash2 := sha256.Sum256([]byte(data2))
	hash3 := sha256.Sum256([]byte(data3))

	fmt.Printf("   \"%s\" -> %x\n", data1, hash1)
	fmt.Printf("   \"%s\" -> %x\n", data2, hash2)
	fmt.Printf("   \"%s\" -> %x\n", data3, hash3)
	fmt.Printf("   data1 == data2: %v (哈希相同)\n", hash1 == hash2)
	fmt.Printf("   data1 == data3: %v (哈希不同，即使只差一个字符)\n", hash1 == hash3)

	// 7. HMAC - 带密钥的哈希消息认证码
	fmt.Println("\n7. HMAC 消息认证:")
	secret := []byte("my-secret-key")
	message := "Important message"

	// 使用 HMAC-SHA256
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	signature := h.Sum(nil)
	fmt.Printf("   消息: \"%s\"\n", message)
	fmt.Printf("   密钥: \"%s\"\n", string(secret))
	fmt.Printf("   HMAC-SHA256: %x\n", signature)

	// 验证 HMAC
	h2 := hmac.New(sha256.New, secret)
	h2.Write([]byte(message))
	expectedSignature := h2.Sum(nil)
	isValid := hmac.Equal(signature, expectedSignature)
	fmt.Printf("   签名验证: %v\n", isValid)

	// 使用错误的密钥验证
	wrongSecret := []byte("wrong-key")
	h3 := hmac.New(sha256.New, wrongSecret)
	h3.Write([]byte(message))
	wrongSignature := h3.Sum(nil)
	isValidWrong := hmac.Equal(signature, wrongSignature)
	fmt.Printf("   错误密钥验证: %v\n", isValidWrong)

	// 8. 文件内容哈希计算
	fmt.Println("\n8. 文件内容哈希:")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("无法获取当前文件路径")
	}
	dir := filepath.Dir(filename)
	testFile := filepath.Join(dir, "go-advanced-input.txt")

	// 先确保文件存在
	if err := os.WriteFile(testFile, []byte("File content for hashing"), 0o644); err != nil {
		fmt.Printf("   无法创建测试文件: %v\n", err)
	} else {
		fileHash, err := hashFile(testFile)
		if err != nil {
			fmt.Printf("   计算文件哈希失败: %v\n", err)
		} else {
			fmt.Printf("   文件: %s\n", filepath.Base(testFile))
			fmt.Printf("   SHA256: %s\n", fileHash)
		}
	}

	// 9. 多次写入累积哈希
	fmt.Println("\n9. 流式哈希计算（多次写入）:")
	streamHash := sha256.New()
	chunks := []string{"Hello", " ", "World", "!"}
	for i, chunk := range chunks {
		streamHash.Write([]byte(chunk))
		fmt.Printf("   写入第 %d 块: \"%s\"\n", i+1, chunk)
	}
	finalHash := streamHash.Sum(nil)
	fmt.Printf("   最终哈希: %x\n", finalHash)

	// 对比一次性哈希
	oneTimeHash := sha256.Sum256([]byte("Hello World!"))
	fmt.Printf("   一次性哈希: %x\n", oneTimeHash)
	fmt.Printf("   两种方式结果相同: %v\n", hex.EncodeToString(finalHash) == hex.EncodeToString(oneTimeHash[:]))

	fmt.Println()
}

// hashFile 计算文件的 SHA256 哈希值
func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	// jsonDemo()
	// fileDemo()
	// timeDemo()
	hashDemo()
}
