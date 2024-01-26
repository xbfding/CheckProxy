package smux

import "CheckProxy/glider/proxy"

func init() {
	proxy.AddUsage("smux", `
Smux scheme:
  smux://host:port
`)
}
