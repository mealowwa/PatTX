package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"net/url"
)

func main() {
	// 请将文件路径替换为你的文件路径
	filePath := "xss.log"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// 判断每行是否包含 "URL:" 字符串
		if strings.Contains(line, "URL:") {
			// 提取 URL 部分
			urlStart := strings.Index(line, "URL:") + len("URL:")
			urlPart := strings.TrimSpace(line[urlStart:])

			parsedURL, err := url.Parse(urlPart)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}

			path := parsedURL.Path

			// 判断每行是否包含 "Param:" 字符串
			if strings.Contains(line, "Param:") {
				// 提取 Param 部分
				paramStart := strings.Index(line, "Param:") + len("Param:")
				paramPart := strings.TrimSpace(line[paramStart:])
				// 将参数附加到路径上
				urlWithParam := fmt.Sprintf("%s  %s", path, paramPart)
				fmt.Println(urlWithParam)
			} else {
				
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
