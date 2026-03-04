package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 网页内容结构
type PageContent struct {
	URL  string
	Body []byte
	Err  error
}

// Web Crawler
type WebCrawler struct {
	maxConcurrency int
	timeout        time.Duration
}

func NewWebCrawler(maxConcurrency int, timeout time.Duration) *WebCrawler {
	return &WebCrawler{
		maxConcurrency: maxConcurrency,
		timeout:        timeout,
	}
}

// 抓取单个URL
func (wc *WebCrawler) fetchURL(url string) (*PageContent, error) {
	client := &http.Client{
		Timeout: wc.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body := make([]byte, 4096) // 简化示例
	n, _ := resp.Body.Read(body)

	return &PageContent{
		URL:  url,
		Body: body[:n],
	}, nil
}

// 并发抓取多个URL
func (wc *WebCrawler) CrawlURLs(urls []string) []*PageContent {
	// 创建worker pool
	jobs := make(chan string, len(urls))
	results := make(chan *PageContent, len(urls))
	var wg sync.WaitGroup

	// 启动worker
	for i := 0; i < wc.maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				content, err := wc.fetchURL(url)
				if err != nil {
					content = &PageContent{
						URL: url,
						Err: err,
					}
				}
				results <- content
			}
		}()
	}

	// 发送任务
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	// 等待完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	var contents []*PageContent
	for content := range results {
		contents = append(contents, content)
	}

	return contents
}

func main() {
	fmt.Println("=== 简单Web爬虫示例 ===\n")

	// 创建爬虫
	crawler := NewWebCrawler(3, 5*time.Second)

	// 要抓取的URL列表
	urls := []string{
		"https://www.example.com",
		"https://httpbin.org/get",
		"https://jsonplaceholder.typicode.com/posts/1",
	}

	fmt.Println("开始并发抓取...")
	start := time.Now()

	contents := crawler.CrawlURLs(urls)

	duration := time.Since(start)

	// 显示结果
	fmt.Printf("\n抓取完成，耗时: %v\n\n", duration)
	fmt.Println("结果:")
	for i, content := range contents {
		if content.Err != nil {
			fmt.Printf("  %d. URL: %s, 错误: %v\n", i+1, content.URL, content.Err)
		} else {
			fmt.Printf("  %d. URL: %s, 内容长度: %d bytes\n", i+1, content.URL, len(content.Body))
		}
	}
}
