package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func param() string {
	// 读取 param3.txt 文件的内容
	paramFile, err := os.Open("params.txt")
	if err != nil {
		panic(err)
	}
	defer paramFile.Close()

	var parames []string
	scanner := bufio.NewScanner(paramFile)
	for scanner.Scan() {
		parames = append(parames, scanner.Text())
	}

	totalLines := len(parames)
	// fmt.Println(numLines)
	paramList := []string{}
	num := 1
	count := 100

	for _, par := range parames {
		if len(par) < 40 {
			parSub := par + "=1%22%3e%3cz0%3e" + "x" + fmt.Sprint(num) + "x"
			num = (num % count) + 1
			paramList = append(paramList, parSub+"&")
		}
	}

	for _, par := range parames {
		if len(par) < 40 {
			parSub := par + "=1%22z0" + "x" + fmt.Sprint(num) + "x"
			num = (num % count) + 1
			paramList = append(paramList, parSub+"&")
		}
	}

	for _, par := range parames {
		if len(par) < 40 {
			parSub := par + "=1%27z1" + "x" + fmt.Sprint(num) + "x"
			num = (num % count) + 1
			paramList = append(paramList, parSub+"&")
		}
	}

	for _, par := range parames {
		if len(par) < 40 {
			parSub := par + "=1%3cz2" + "x" + fmt.Sprint(num) + "x"
			num = (num % count) + 1
			paramList = append(paramList, parSub+"&")
		}
	}

	paramNum := 0
	totalNum := 0

	tempFile, err := os.CreateTemp("", "temp_chuanlian_*.txt")
	if err != nil {
		panic(err)
	}
	//defer os.Remove(tempFile.Name())

	for _, i := range paramList {
		// fmt.Println("i:" + i)
		paramNum++
		totalNum++

		if paramNum < count && totalNum != totalLines {
			tempFile.WriteString(i)
		}

		if paramNum == count {
			i = strings.TrimRight(i, "&")
			tempFile.WriteString(i + "\n")
			paramNum = 0
		}

		if totalNum == totalLines {
			tempFile.WriteString(i + "\n")
			totalNum = 0
			paramNum = 0
		}
	}

	// fmt.Println("参数串联完毕！")
	return tempFile.Name()
}

func main() {
	tempChuanlian := param()

	urlsFile, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
	defer urlsFile.Close()

	var urls []string
	scanner := bufio.NewScanner(urlsFile)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	paramFile, err := os.Open(tempChuanlian)
	if err != nil {
		panic(err)
	}
	defer paramFile.Close()

	var params []string
	scanner = bufio.NewScanner(paramFile)
	for scanner.Scan() {
		params = append(params, scanner.Text())
	}

	var combinations []string
	for _, url := range urls {
		for _, param := range params {
			combinations = append(combinations, url+"?"+param)
		}
	}

	outputFile, err := os.Create("out.log")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	outputFile.WriteString(strings.Join(combinations, "\n"))

	fmt.Println("Combination script executed successfully. Output saved in out.log.")
}
