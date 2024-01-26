package glider

import (
	// comment out the services you don't need to make the compiled binary smaller.
	// _ "CheckProxy/glider/service/xxx"

	// comment out the protocols you don't need to make the compiled binary smaller.
	_ "CheckProxy/glider/proxy/http"
	_ "CheckProxy/glider/proxy/mixed"
	_ "CheckProxy/glider/proxy/obfs"
	_ "CheckProxy/glider/proxy/pxyproto"
	_ "CheckProxy/glider/proxy/reject"
	_ "CheckProxy/glider/proxy/smux"
	_ "CheckProxy/glider/proxy/socks4"
	_ "CheckProxy/glider/proxy/socks5"
	_ "CheckProxy/glider/proxy/ss"
	_ "CheckProxy/glider/proxy/ssh"
	_ "CheckProxy/glider/proxy/ssr"
	_ "CheckProxy/glider/proxy/tcp"
	_ "CheckProxy/glider/proxy/tls"
	_ "CheckProxy/glider/proxy/trojan"
	_ "CheckProxy/glider/proxy/udp"
	_ "CheckProxy/glider/proxy/vless"
	_ "CheckProxy/glider/proxy/vmess"
	_ "CheckProxy/glider/proxy/ws"
)
