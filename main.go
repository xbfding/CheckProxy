package main

import (
	"CheckProxy/glider"
	"CheckProxy/pkg/sqlite"
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	h             bool
	dbFilePath    string
	proxyFilePath string
	FileName      string
	forward       []string
)

func GetFileName(FilePath string) (FileName string, FileSuffix string) {

	str := filepath.Base(FilePath)
	lastDot := strings.LastIndex(FilePath, ".")
	if lastDot == -1 {
		return FilePath, ""
	}

	FileName = str[:lastDot-2]
	FileSuffix = str[lastDot-1:]
	return FileName, FileSuffix

}

func usage() {
	fmt.Fprintf(os.Stderr, `代理检查器

Options:
`)
	// 自定义参数的排序顺序
	paramOrder := []string{"h", "dbFilePath", "proxyFilePath"}

	// 输出参数帮助信息，按照自定义的排序顺序
	for _, param := range paramOrder {
		flag.VisitAll(func(param string) func(*flag.Flag) {
			return func(f *flag.Flag) {
				if f.Name == param {
					fmt.Fprintf(os.Stderr, "  -%s  %s\n", f.Name, f.Usage)
				}
			}
		}(param))
	}
}

func init() {
	flag.BoolVar(&h, "h", false, "帮助")
	flag.StringVar(&dbFilePath, "db", "", "数据库路径")
	flag.StringVar(&proxyFilePath, "txt", "", "代理文件路径")
	flag.Usage = usage
}

func ExpCsv(timestamp string, fwdok []string) {
	//创建CSV文件
	dirPath := "export"
	csvFileName := filepath.Join(dirPath, fmt.Sprintf("xrayProxySub_%s.csv", timestamp))
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatal("[-]", err)
	}
	file, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("[-]", err)
	}
	defer file.Close()

	// 创建 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 自定义表头并写入 CSV 文件
	headers := []string{"glider代理配置"} // 自定义表头
	writer.Write(headers)

	// 遍历数据库查询结果并写入 CSV 文件
	for _, tmp := range fwdok {
		row := []string{tmp}
		writer.Write(row)
	}
}

func DB() {
	forward = sqlite.SelectSubTarget(dbFilePath)
}

func readfile() (fwds []string) {
	file, err := os.Open(proxyFilePath)
	if err != nil {
		log.Fatal("文件打开失败:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fwds = append(fwds, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("文件读取出错:", err)
	}
	return fwds
}

func TXT() {

	for _, tmp := range readfile() {
		tmp = strings.Replace(tmp, "forward=", "", 1)
		tmp = strings.Replace(tmp, "\"", "", -1)
		forward = append(forward, tmp)
	}
}

func startCheck() {
	if len(forward) == 0 {
		log.Fatal("[-]未获得代理!")
	}
	fmt.Println("[*]开始检测代理可用性!")
	fwdok := glider.Main(forward)
	ExpCsv(FileName, fwdok)
	fmt.Println("[*]已导出格式化代理表!")
}
func main() {
	flag.Parse()
	if h || flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	if dbFilePath != "" || proxyFilePath != "" {

		if dbFilePath != "" {
			FileName, _ = GetFileName(dbFilePath)
			DB()
		} else {
			FileName, _ = GetFileName(proxyFilePath)
			TXT()
		}
		startCheck()
	} else {
		println("请输入需检测的代码文件路径")
		flag.Usage()
		return
	}

}
