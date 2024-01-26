package rule

import (
	"CheckProxy/pkg/httpRequest"
	"fmt"
	"hash/fnv"
	"net"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"CheckProxy/glider/pkg/log"
	"CheckProxy/glider/proxy"
)

// forwarder slice orderd by priority.
type priSlice []*Forwarder

func (p priSlice) Len() int           { return len(p) }
func (p priSlice) Less(i, j int) bool { return p[i].Priority() > p[j].Priority() }
func (p priSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// FwdrGroup is a forwarder group.
type FwdrGroup struct {
	name     string
	config   *Strategy
	fwdrs    priSlice
	avail    []*Forwarder // available forwarders
	mu       sync.RWMutex
	index    uint32
	priority uint32
	next     func(addr string) *Forwarder
	localip  string
	//fwdok      []string
	dbFilePath string
}

// NewFwdrGroup returns a new forward group.
func NewFwdrGroup(rulePath string, s []string, c *Strategy) *FwdrGroup {
	var fwdrs []*Forwarder
	for _, chain := range s {
		fwdr, err := ForwarderFromURL(chain, c.IntFace,
			time.Duration(c.DialTimeout)*time.Second, time.Duration(c.RelayTimeout)*time.Second)
		if err != nil {
			fmt.Printf("[-]代理解析错误!(%s)\n", err)
			continue
		}
		fwdr.SetMaxFailures(uint32(c.MaxFailures))
		fwdrs = append(fwdrs, fwdr)
	}

	if len(fwdrs) == 0 {
		log.Fatal("[-]代理格式均无效!")
	}

	name := strings.TrimSuffix(filepath.Base(rulePath), filepath.Ext(rulePath))
	return newFwdrGroup(name, fwdrs, c)
}

// newFwdrGroup returns a new FwdrGroup.
func newFwdrGroup(name string, fwdrs []*Forwarder, c *Strategy) *FwdrGroup {
	p := &FwdrGroup{name: name, fwdrs: fwdrs, config: c}
	sort.Sort(p.fwdrs)

	p.init()

	// default scheduler
	p.next = p.scheduleRR

	// if there're more than 1 forwarders, we care about the strategy.
	if count := len(fwdrs); count > 1 {
		switch c.Strategy {
		case "rr":
			p.next = p.scheduleRR
			log.F("[strategy] %s: %d forwarders forward in round robin mode.", name, count)
		case "ha":
			p.next = p.scheduleHA
			log.F("[strategy] %s: %d forwarders forward in high availability mode.", name, count)
		case "lha":
			p.next = p.scheduleLHA
			log.F("[strategy] %s: %d forwarders forward in latency based high availability mode.", name, count)
		case "dh":
			p.next = p.scheduleDH
			log.F("[strategy] %s: %d forwarders forward in destination hashing mode.", name, count)
		default:
			p.next = p.scheduleRR
			log.F("[strategy] %s: not supported forward mode '%s', use round robin mode for %d forwarders.", name, c.Strategy, count)
		}
	}

	for _, f := range fwdrs {
		f.AddHandler(p.onStatusChanged)
	}

	return p
}

// Dial connects to the address addr on the network net.
func (p *FwdrGroup) Dial(network, addr string) (net.Conn, proxy.Dialer, error) {
	nd := p.NextDialer(addr)
	c, err := nd.Dial(network, addr)
	return c, nd, err
}

// DialUDP connects to the given address.
func (p *FwdrGroup) DialUDP(network, addr string) (pc net.PacketConn, dialer proxy.UDPDialer, err error) {
	nd := p.NextDialer(addr)
	pc, err = nd.DialUDP(network, addr)
	return pc, nd, err
}

// NextDialer returns the next dialer.
func (p *FwdrGroup) NextDialer(dstAddr string) proxy.Dialer {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.avail) == 0 {
		return p.fwdrs[atomic.AddUint32(&p.index, 1)%uint32(len(p.fwdrs))]
	}

	return p.next(dstAddr)
}

// Priority returns the active priority of dialer.
func (p *FwdrGroup) Priority() uint32 { return atomic.LoadUint32(&p.priority) }

// SetPriority sets the active priority of daler.
func (p *FwdrGroup) SetPriority(pri uint32) { atomic.StoreUint32(&p.priority, pri) }

// init traverse d.fwdrs and init the available forwarder slice.
func (p *FwdrGroup) init() {
	for _, f := range p.fwdrs {
		if f.Enabled() {
			p.SetPriority(f.Priority())
			break
		}
	}

	p.avail = nil
	for _, f := range p.fwdrs {
		if f.Enabled() && f.Priority() >= p.Priority() {
			p.avail = append(p.avail, f)
		}
	}

	if len(p.avail) == 0 {
		// no available forwarders, set priority to 0 to check all forwarders in check func
		p.SetPriority(0)
		// log.F("[group] no available forwarders, please check your config file or network settings")
	}
}

// onStatusChanged will be called when fwdr's status changed.
func (p *FwdrGroup) onStatusChanged(fwdr *Forwarder) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if fwdr.Enabled() {
		if fwdr.Priority() == p.Priority() {
			p.avail = append(p.avail, fwdr)
		} else if fwdr.Priority() > p.Priority() {
			p.init()
		}
		log.F("[group] %s: %s(%d) changed status from DISABLED to ENABLED (%d of %d currently enabled)",
			p.name, fwdr.Addr(), fwdr.Priority(), len(p.avail), len(p.fwdrs))
	} else {
		for i, f := range p.avail {
			if f == fwdr {
				p.avail[i], p.avail = p.avail[len(p.avail)-1], p.avail[:len(p.avail)-1]
				break
			}
		}
		log.F("[group] %s: %s(%d) changed status from ENABLED to DISABLED (%d of %d currently enabled)",
			p.name, fwdr.Addr(), fwdr.Priority(), len(p.avail), len(p.fwdrs))
	}

	if len(p.avail) == 0 {
		p.init()
	}
}

// Check runs the forwarder checks.
func (p *FwdrGroup) Check() (fwdok []string) {

	if !strings.Contains(p.config.Check, "://") {
		p.config.Check += "://"
	}

	u, err := url.Parse(p.config.Check)
	if err != nil {
		log.F("[group] %s: parse check config error: %s, disable health checking", p.name, err)
		return
	}

	addr := u.Host
	timeout := time.Duration(p.config.CheckTimeout) * time.Second

	var checker Checker
	expect := "HTTP" // default: check the first 4 chars in response
	params, _ := url.ParseQuery(u.Fragment)
	if ex := params.Get("expect"); ex != "" {
		expect = ex
	}
	checker = newHttpChecker(addr, u.RequestURI(), expect, timeout, u.Scheme == "https")

	log.F("[+] %s: using check config: %s", p.name, p.config.Check)
	p.myIp() // get location ip
	fwdoks := p.runCheck(checker)
	return fwdoks
}

func (p *FwdrGroup) runCheck(checker Checker) (fwdoks []string) {
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 30)

	for i := 0; i < len(p.fwdrs); i++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(idx int, c Checker) {
			defer wg.Done()

			//p.mu.Lock() //加锁保护共享数据
			fwdok := p.check(p.fwdrs[idx], c)
			if fwdok != "" {
				fwdoks = append(fwdoks, fwdok)
			}
			//p.mu.Unlock() //解锁
			<-ch
		}(i, checker)
	}
	wg.Wait()
	return fwdoks
}

func (p *FwdrGroup) myIp() {
	Headers := map[string]string{
		"User-Agent": "curl/8.4.0",
		"Accept":     "*/*",
	}
	for i := 0; i < 3; i++ {
		statusCode, _, respBody, err := httpRequest.GetResq("https://www.taobao.com/help/getip.php", Headers, "")
		if statusCode != 200 {
			fmt.Println("[-](https://www.taobao.com/help/getip.php)触发风控,正在重试!\n[-]响应内容:", string(respBody))
		} else if err != nil {
			log.Fatal("[-]真实IP获取失败:", err)
		} else {
			p.localip, _ = ExtractMyipInfo(string(respBody))
			return
		}
		time.Sleep(time.Duration(3) * time.Second)
	}

}

func (p *FwdrGroup) check(fwdr *Forwarder, checker Checker) (fwdok string) {
	for i := 0; i < 3; i++ {
		fwdrUrl, currentIP, err := checker.Check(fwdr)
		if err != nil {
			if err.Error() == "timeout" {
				fmt.Printf("[-]请求超时,正在重试!(%s)\n", fwdr.url)
				continue
			}
			if err.Error() == "nil" {
				fmt.Printf("[-]验证失败:(%s)\n", fwdrUrl)
				continue
			}
			if err.Error() == "NO CN" {
				fmt.Printf("[-]非大陆,不记录:(%s)\n", fwdrUrl)
				break
			}
			fmt.Printf("[-]验证发生错误:(%s)(%s)\n", fwdrUrl, err.Error())
			continue
		}
		if currentIP != p.localip {
			fmt.Printf("[+]大陆存活,记录:(%s)\n", fwdrUrl)
			return fwdrUrl
		}
		intval := time.Duration(3) * time.Second
		time.Sleep(intval)
	}
	return ""
}

// Round Robin.
func (p *FwdrGroup) scheduleRR(dstAddr string) *Forwarder {
	return p.avail[atomic.AddUint32(&p.index, 1)%uint32(len(p.avail))]
}

// High Availability.
func (p *FwdrGroup) scheduleHA(dstAddr string) *Forwarder {
	return p.avail[0]
}

// Latency based High Availability.
func (p *FwdrGroup) scheduleLHA(dstAddr string) *Forwarder {
	oldfwdr, newfwdr := p.avail[0], p.avail[0]
	lowest := oldfwdr.Latency()
	for _, f := range p.avail {
		if f.Latency() < lowest {
			lowest = f.Latency()
			newfwdr = f
		}
	}
	tolerance := int64(p.config.CheckTolerance) * int64(time.Millisecond)
	if newfwdr.Latency() < (oldfwdr.Latency() - tolerance) {
		return newfwdr
	}
	return oldfwdr
}

// Destination Hashing.
func (p *FwdrGroup) scheduleDH(dstAddr string) *Forwarder {
	fnv1a := fnv.New32a()
	fnv1a.Write([]byte(dstAddr))
	return p.avail[fnv1a.Sum32()%uint32(len(p.avail))]
}
