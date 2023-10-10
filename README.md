## 检查代理可用性

目前仅支持`http`和`socks5`协议

## 使用说明

代理文件格式：逐行排列即可

```sh
socks5://127.0.0.1:7890
http://127.0.0.1:7890
```

使用命令

```shell
./CheckProxy -pfile ./proxyfile.txt
```

## 编译

```shell
#最小体积编译
go build -ldflags="-s -w" .\CheckProxy.go
```

