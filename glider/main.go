package glider

import (
	"CheckProxy/glider/rule"
)

var (
	version = "0.16.3"
)

func Main(forwarder []string) (fwdok []string) {
	config := parseConfig(forwarder)
	// global rule proxy
	pxy := rule.NewProxy(config.Forwards, &config.Strategy, config.dbFilePath)
	// enable checkers
	return pxy.Check()
}
