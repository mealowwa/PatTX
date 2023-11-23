package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"io"
)

func colorize(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return fmt.Sprintf("\x1b[32m%d\x1b[0m", statusCode) // 绿色
	case statusCode >= 400 && statusCode < 500:
		return fmt.Sprintf("\x1b[31m%d\x1b[0m", statusCode) // 红色
	case statusCode >= 500:
		return fmt.Sprintf("\x1b[31m%d\x1b[0m", statusCode) // 红色
	default:
		return fmt.Sprintf("%d", statusCode)
	}
}

func containsString(content, substring string) bool {
	return strings.Contains(content, substring)
}

func main() {
	// 使用bufio.NewReader读取标准输入
	scanner := bufio.NewScanner(os.Stdin)

	// 使用WaitGroup等待所有goroutines完成
	var wg sync.WaitGroup

	// 创建不跟随重定向的Client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 不跟随重定向
			return http.ErrUseLastResponse
		},
	}

	// 逐行读取输入并发起GET请求
	for scanner.Scan() {
		url := scanner.Text()

		// 增加WaitGroup计数
		wg.Add(1)
        
		// 启动goroutine处理请求
		go func(url string) {
			defer wg.Done()

            var statusCode int

			var respbody string

			for i := 0; i < 100; i++ {
				// 发起GET请求
				resp, err := client.Get(url)

				if err != nil {
					return
				}
				defer resp.Body.Close()

				// 获取状态码，并高亮显示
				statusCode = resp.StatusCode

				//fmt.Printf("%s: Status %s\n", url, colorize(statusCode))
				

				// 读取响应内容
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}

				if statusCode == 503 || statusCode == 429 {
					continue
				}

				respbody = string(body)
				// 检查响应内容中是否包含特定字符串
				if containsString(respbody, "Type the characters you see in this image") {

				} else {
					break // 如果不包含特定字符串，退出循环
				}

			}
 
			if statusCode == 200 {
				//fmt.Printf("%s: Status %s\n", url, colorize(statusCode))
				fmt.Println(url)
			} 

			
			
		}(url)
	}

	// 等待所有goroutines完成
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}