## 检查代理可用性

- 内置glider,支持多个协议
- 支持多线程检查，默认30线程保证稳定性，失败重试3次
- 导出检测成功的大陆代理

## 使用说明

*代理文件格式为glider配置文件格式，文末附上。目前有clash转glider配置文件脚本[yamlToGliderConfig](https://github.com/xbfding/yamlToGliderConfig)，将json每行排列输入到脚本即可，后续有时间会优化*

**使用命令**

```shell
#快速开始
./CheckProxy -txt ./proxyfile.txt
```

**帮助**

```shell
Options:
  -h  帮助
  -db  数据库路径
  -txt  代理文件路径
```

## 编译

```shell
#最小体积编译
go build -ldflags="-s -w"
```

## 计划任务

- [ ] 增加导出海外代理
- [ ] 导出时附加IP定位以及延迟

## 附录

**代理格式如下**

```shell
# FORWARDERS
# ----------
# Forwarders, we can setup multiple forwarders.
# forward=SCHEME#OPTIONS

# FORWARDER OPTIONS
# priority: set the priority of that forwarder, default:0
# interface: set local interface or ip address used to connect remote server

# Socks5 proxy as forwarder
# forward=socks5://192.168.1.10:1080

# Socks5 proxy as forwarder with priority 100
# forward=socks5://192.168.1.10:1080#priority=100

# Socks5 proxy as forwarder with priority 100 and use `eth0` as source interface
# forward=socks5://192.168.1.10:1080#priority=100&interface=eth0

# Socks5 proxy as forwarder with priority 100 and use `192.168.1.100` as source ip
# forward=socks5://192.168.1.10:1080#priority=100&interface=192.168.1.100

# SS proxy as forwarder
# forward=ss://method:pass@1.1.1.1:8443

# SSR proxy as forwarder
# forward=ssr://method:pass@1.1.1.1:8443?protocol=auth_aes128_md5&protocol_param=xxx&obfs=tls1.2_ticket_auth&obfs_param=yyy

# ssh forwarder
# forward=ssh://user[:pass]@host:port[?key=keypath&timeout=SECONDS]
# forward=ssh://root:pass@host:port
# forward=ssh://root@host:port?key=/path/to/keyfile
# forward=ssh://root@host:port?key=/path/to/keyfile&timeout=5

# http proxy as forwarder
# forward=http://1.1.1.1:8080

# trojan as forwarder
# forward=trojan://PASSWORD@1.1.1.1:8080[?serverName=SERVERNAME][&skipVerify=true]

# trojanc as forwarder
# forward=trojanc://PASSWORD@1.1.1.1:8080

# vless forwarder
# forward=vless://5a146038-0b56-4e95-b1dc-5c6f5a32cd98@1.1.1.1:443

# vmess with aead auth
# forward=vmess://5a146038-0b56-4e95-b1dc-5c6f5a32cd98@1.1.1.1:443

# vmess with md5 auth (by setting alterID)
# forward=vmess://5a146038-0b56-4e95-b1dc-5c6f5a32cd98@1.1.1.1:443?alterID=2

# vmess over tls
# forward=tls://server.com:443,vmess://5a146038-0b56-4e95-b1dc-5c6f5a32cd98

# vmess over websocket
# forward=ws://1.1.1.1:80/path?host=server.com,vmess://chacha20-poly1305:5a146038-0b56-4e95-b1dc-5c6f5a32cd98

# vmess over ws over tls
# forward=tls://server.com:443,ws://,vmess://5a146038-0b56-4e95-b1dc-5c6f5a32cd98
# forward=tls://server.com:443,ws://@/path,vmess://5a146038-0b56-4e95-b1dc-5c6f5a32cd98

# ss over tls
# forward=tls://server.com:443,ss://AEAD_CHACHA20_POLY1305:pass@

# ss over kcp
# forward=kcp://aes:key@127.0.0.1:8444?dataShards=10&parityShards=3&mode=fast,ss://AEAD_CHACHA20_POLY1305:pass@

# ss with simple-obfs
# forward=simple-obfs://1.1.1.1:443?type=tls&host=apple.com,ss://AEAD_CHACHA20_POLY1305:pass@

# socks5 over unix domain socket
# forward=unix:///dev/shm/socket,socks5://
```