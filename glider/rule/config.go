package rule

// Config is config of rule.
type Config struct {
	RulePath string

	Forward  []string
	Strategy Strategy

	DNSServers []string
	IPSet      string

	Domain []string
	IP     []string
	CIDR   []string
}

// Strategy configurations.
type Strategy struct {
	Strategy            string
	Check               string
	CheckInterval       int
	CheckTimeout        int
	CheckTolerance      int
	CheckLatencySamples int
	CheckDisabledOnly   bool
	MaxFailures         int
	DialTimeout         int
	RelayTimeout        int
	IntFace             string
}
