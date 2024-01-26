package rule

import (
	"log"
	"regexp"
	"strings"
)

// ExtractMyipInfo 提取myip.ipip.net响应中的值
func ExtractMyipInfo(responseBody string) (currentIP string, location string) {

	// 定义正则表达式
	regexIP := `((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`
	regexLocation := `来自：([^,]*)\n`

	// 编译正则表达式
	reIP, err := regexp.Compile(regexIP)
	if err != nil {
		log.Fatal("[-]正则表达式编译失败:", err)
	}
	reLt, err := regexp.Compile(regexLocation)
	if err != nil {
		log.Fatal("[-]正则表达式编译失败:", err)
	}

	// 查找匹配的字符串
	currentIP = reIP.FindString(responseBody)
	location = reLt.FindString(responseBody)
	location = strings.Replace(location, "来自：", "", -1)
	location = strings.Replace(location, "\n", "", -1)
	return currentIP, location
}

// CheckIpCN 检查IP归属地
func CheckIpCN(location string) bool {
	if strings.Contains(location, "中国") || strings.Contains(location, "\\u4e2d\\u56fd") {
		switch {
		case strings.Contains(location, "香港") || strings.Contains(location, "\\u9999\\u6e2f"):
			return false
		case strings.Contains(location, "台湾") || strings.Contains(location, "\\u53f0\\u6e7e"):
			return false
		default:
			return true
		}
	}
	return false
}
