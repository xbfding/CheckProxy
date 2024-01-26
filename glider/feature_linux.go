package glider

import (
	// comment out the services you don't need to make the compiled binary smaller.
	_ "CheckProxy/glider/service/dhcpd"

	// comment out the protocols you don't need to make the compiled binary smaller.
	_ "CheckProxy/glider/proxy/redir"
	_ "CheckProxy/glider/proxy/tproxy"
	_ "CheckProxy/glider/proxy/unix"
	_ "CheckProxy/glider/proxy/vsock"
)
