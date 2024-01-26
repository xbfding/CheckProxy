package glider

import (
	"CheckProxy/glider/dns"
	"CheckProxy/glider/rule"
)

// Config is global config struct.
type Config struct {
	Verbose    bool
	LogFlags   int
	TCPBufSize int
	UDPBufSize int

	Listens []string

	Forwards []string
	Strategy rule.Strategy

	RuleFiles []string
	RulesDir  string

	DNS       string
	DNSConfig dns.Config

	rules []*rule.Config

	Services   []string
	dbFilePath string
}

func parseConfig(forwarder []string) *Config {
	conf := &Config{}

	conf.Forwards = forwarder

	conf.Strategy.Strategy = "rr"
	conf.Strategy.Check = "https://forge.speedtest.cn/api/location/info"
	conf.Strategy.CheckInterval = 60
	conf.Strategy.CheckTimeout = 3
	conf.Strategy.CheckTolerance = 100
	conf.Strategy.CheckLatencySamples = 3
	conf.Strategy.CheckDisabledOnly = false
	conf.Strategy.MaxFailures = 3
	conf.Strategy.DialTimeout = 3
	conf.Strategy.RelayTimeout = 0
	conf.Strategy.IntFace = ""

	return conf
}
