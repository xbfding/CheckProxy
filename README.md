## 检查代理可用性

- 目前仅支持`http`和`socks5`协议

- 支持多线程检查，默认100线程，可以指定线程

## 使用说明

**代理文件格式：逐行排列即可**

例如：

```sh
socks5://127.0.0.1:7890
http://127.0.0.1:7890
```

**使用命令**

```shell
#快速开始
./CheckProxy -pfile ./proxyfile.txt

#指定线程数量
./CheckProxy -pfile ./proxyfile.txt -t 300
```

## 编译

```shell
#最小体积编译
go build -ldflags="-s -w" .\CheckProxy.go
```

