package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func removeDuplication_map(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

// 提取myip.ipip.net响应中的值
func extractMyipInfo(responseBody string) (string, string) {
	var currentIP, location string

	// 定义正则表达式
	regexIP := `((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`
	regexLocation := `来自于：([^,]*)\n`

	// 编译正则表达式
	reIP, err := regexp.Compile(regexIP)
	if err != nil {
		log.Fatal("正则表达式编译失败:", err)
	}
	reLt, err := regexp.Compile(regexLocation)
	if err != nil {
		log.Fatal("正则表达式编译失败:", err)
	}

	// 查找匹配的字符串
	currentIP = reIP.FindString(responseBody)
	location = reLt.FindString(responseBody)
	location = strings.Replace(location, "来自于：", "", -1)
	location = strings.Replace(location, "\n", "", -1)
	return currentIP, location
}

func testProxy(proxyAddr string, trueIP string) (proxyTrue string, TrueIP string, proxyLocation string) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("解析代理地址出错:", err)
		return
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	resp, err := client.Get("http://myip.ipip.net")
	if err != nil {
		fmt.Println("[-]发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[-]读取响应内容失败:", err)
		return
	}

	bodyString := string(bodyBytes)
	currentIP, location := extractMyipInfo(bodyString)
	if currentIP == trueIP {
		fmt.Printf("[-]小心蜜罐!!![%s]", proxyAddr)
	}
	fmt.Printf("[+]可用代理：[%s],ip为：[%s],归属地：[%s]\n", proxyAddr, currentIP, location)
	return proxyAddr, currentIP, location

}

func CheckIpCN(location string) bool {
	switch {
	case strings.Contains(location, "中国"):
		return true
	default:
		return false
	}
}

// 主控制台
func MainlandProxyMultiThread(urls []string, numThreads int) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, numThreads)
	results := make(chan string)

	// 数据去重
	urls = removeDuplication_map(urls)
	println("去重后数据量：", len(urls))

	/**
	      // 规范代理格式
	  	var socks5Proxy []string
	  	for _, item := range urls {
	  		item = fmt.Sprintf("%s%s", "socks5://", item)
	  		socks5Proxy = append(socks5Proxy, item)
	  	}
	  **/
	resp, err := http.Get("http://myip.ipip.net")
	if err != nil {
		fmt.Printf("get.myip.err:%s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	currentIP, _ := extractMyipInfo(string(body))

	// 启动一个协程用于显示进度
	go func() {
		completed := 0
		for range results {
			completed++
			fmt.Printf("\033[2K完成进度: %d/%d\r", completed, len(urls))
		}
	}()

	// 遍历URL列表
	for _, Proxy := range urls {
		wg.Add(1)

		// 启动一个协程处理每个URL
		go func(url, proxy string) { // 将 Proxy 作为参数传递给闭包
			defer wg.Done()

			// 获取信号量
			sem <- struct{}{}
			defer func() { <-sem }() // 释放信号量

			// 发送GET请求并获取响应内容
			proxyTrue, TrueIP, proxyLocation := testProxy(proxy, currentIP)
			if proxyTrue != "" {
				if CheckIpCN(proxyLocation) {
					writeStr := fmt.Sprintf("%s        ip为：[%s],归属地：[%s]\n", proxyTrue, TrueIP, proxyLocation)
					file, pwdErr := os.OpenFile("MainlandProxy_out.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if pwdErr != nil {
						fmt.Printf("[-]MainlandProxy文件打开错误：%v", err)
					}
					defer file.Close()
					_, err := file.WriteString(writeStr)
					if err != nil {
						fmt.Printf("[-]MainlandProxy文件写入错误：%v", err)
					}
				} else {
					writeStr := fmt.Sprintf("%s        ip为：[%s],归属地：[%s]\n", proxyTrue, TrueIP, proxyLocation)
					file, pwdErr := os.OpenFile("OtherProxy_out.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if pwdErr != nil {
						fmt.Printf("[-]OtherProxy文件打开错误：%v", err)
					}
					defer file.Close()
					_, err := file.WriteString(writeStr)
					if err != nil {
						fmt.Printf("[-]OtherProxy文件写入错误：%v", err)
					}
				}

			}

			// 向结果通道发送信号
			results <- url
		}(Proxy, Proxy) // 将 Proxy 作为参数传递给 goroutine
	}

	// 等待所有协程完成
	wg.Wait()
	close(results)

	fmt.Println("[*]检查代理的任务完成啦！！！")
}

func main() {
	CheckProxy := `

	   ________              __   ____                       
	  / ____/ /_  ___  _____/ /__/ __ \_________  _  ____  __
	 / /   / __ \/ _ \/ ___/ //_/ /_/ / ___/ __ \| |/_/ / / /
	/ /___/ / / /  __/ /__/ ,< / ____/ /  / /_/ />  </ /_/ / 
	\____/_/ /_/\___/\___/_/|_/_/   /_/   \____/_/|_|\__, /  
	                                                /____/
	V0.1.0
	`
	fmt.Println(CheckProxy)

	var (
		proxyFilePath string
		numThreads    int
	)

	flag.StringVar(&proxyFilePath, "pfile", "", "需要检查的代理文件路径")
	flag.IntVar(&numThreads, "t", 100, "线程数量,没有特殊需求可以不指定")
	// 解析命令行参数
	fmt.Println("\n例如：./CheckProxy -pfile ./proxyfile.txt\n")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	proxys, err := os.ReadFile(proxyFilePath)
	if err != nil {
		fmt.Println("[-]文件打开错误：", err)
	}
	// 使用 bufio 包创建一个字符串切片来存储每一行的内容，去掉行尾的换行符
	var proxysString []string
	scanner := bufio.NewScanner(strings.NewReader(string(proxys)))
	for scanner.Scan() {
		line := scanner.Text()
		proxysString = append(proxysString, line)
	}

	// 检查扫描时是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时发生错误:", err)
		return
	}

	MainlandProxyMultiThread(proxysString, numThreads)
}
