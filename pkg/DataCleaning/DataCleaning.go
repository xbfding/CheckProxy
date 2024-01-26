package DataCleaning

import (
	"fmt"
	"strings"
)

// DataGL 添加http协议头
func DataGL(targets []string) []string {
	var dataQX []string
	for _, item := range targets {
		if strings.HasPrefix(item, "http://") {
			dataQX = append(dataQX, item)
		} else if strings.HasPrefix(item, "https://") {
			dataQX = append(dataQX, item)
		} else {
			targe := fmt.Sprintf("%s%s", "http://", item)
			dataQX = append(dataQX, targe)
		}
	}
	return dataQX
}

// RemoveDuplicationMap 数据去重
func RemoveDuplicationMap(arr []string) []string {
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
