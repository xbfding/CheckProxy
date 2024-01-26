package rule

// Proxy implements the proxy.Proxy interface with rule support.
type Proxy struct {
	main *FwdrGroup
	//all       []*FwdrGroup
	//domainMap sync.Map
	//ipMap     sync.Map
	//cidrMap   sync.Map
}

// NewProxy returns a new rule proxy.
func NewProxy(mainForwarders []string, mainStrategy *Strategy, dbFilePath string) *Proxy {
	rd := &Proxy{main: NewFwdrGroup("main", mainForwarders, mainStrategy)}
	rd.main.dbFilePath = dbFilePath

	return rd
}

// Check checks availability of forwarders inside proxy.
func (p *Proxy) Check() (fwdok []string) {
	return p.main.Check()
}
