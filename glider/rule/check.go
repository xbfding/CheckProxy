package rule

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"CheckProxy/glider/pkg/pool"
)

// Checker is a forwarder health checker.
type Checker interface {
	Check(fwdr *Forwarder) (proxy string, currentIP string, err error)
}

type httpChecker struct {
	addr    string
	uri     string
	expect  string
	timeout time.Duration

	tlsConfig  *tls.Config
	serverName string

	regex *regexp.Regexp
}

func newHttpChecker(addr, uri, expect string, timeout time.Duration, withTLS bool) *httpChecker {
	c := &httpChecker{
		addr:    addr,
		uri:     uri,
		expect:  expect,
		timeout: timeout,
		regex:   regexp.MustCompile(expect),
	}

	if _, p, _ := net.SplitHostPort(addr); p == "" {
		if withTLS {
			c.addr = net.JoinHostPort(addr, "443")
		} else {
			c.addr = net.JoinHostPort(addr, "80")
		}
	}
	c.serverName = c.addr[:strings.LastIndex(c.addr, ":")]
	if withTLS {
		c.tlsConfig = &tls.Config{ServerName: c.serverName}
	}
	return c
}

// Check implements the Checker interface.
func (c *httpChecker) Check(fwdr *Forwarder) (proxy string, currentIP string, err error) {
	startTime := time.Now()
	rc, err := fwdr.Dial("tcp", c.addr)
	if err != nil {
		return fwdr.url, "", err
	}

	if c.tlsConfig != nil {
		tlsConn := tls.Client(rc, c.tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			tlsConn.Close()
			return fwdr.url, "", err
		}
		rc = tlsConn
	}
	defer rc.Close()

	if c.timeout > 0 {
		rc.SetDeadline(time.Now().Add(c.timeout))
	}
	req := fmt.Sprintf("GET %s HTTP/1.1\nHost: %s\nConnection: keep-alive\nCache-Control: max-age=0\nsec-ch-ua: \"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"\nsec-ch-ua-mobile: ?0\nsec-ch-ua-platform: \"Windows\"\nUpgrade-Insecure-Requests: 1\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\nSec-Fetch-Site: none\nSec-Fetch-Mode: navigate\nSec-Fetch-User: ?1\nSec-Fetch-Dest: document\nAccept-Encoding: gzip, deflate, br, zstd\nAccept-Language: zh-CN,zh;q=0.9\n\n", c.uri, c.serverName)
	if _, err = io.WriteString(rc, req); err != nil {
		return fwdr.url, "", err
	}

	r := pool.GetBufReader(rc)
	defer pool.PutBufReader(r)
	elapsed := time.Since(startTime)
	line, _ := r.ReadString(0)
	if err != nil {
		return fwdr.url, "", err
	}
	if elapsed > c.timeout {
		return "", "", errors.New("timeout")
	}
	if line == "" {
		return fwdr.url, "", errors.New("nil")
	}
	currentIP, _ = ExtractMyipInfo(line)

	if CheckIpCN(line) {
		return fwdr.url, currentIP, nil
	}
	return fwdr.url, currentIP, errors.New("NO CN")
}
