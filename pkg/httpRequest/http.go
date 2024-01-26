package httpRequest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var GetHeaderAll = map[string]string{
	"Accept":          "*/*",
	"Accept-Language": "en-US;q=0.9,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.102 Safari/537.36",
	"Connection":      "close",
	"Cache-Control":   "max-age=0",
}

func Headreq(targe string, proxy string) (statusCode int, err error) {
	req, err := http.NewRequest("HEAD", targe, nil)
	if err != nil {
		return 0, fmt.Errorf("Headreq创建请求出现错误:%v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")

	// 创建自定义的Transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: func(req *http.Request) (*url.URL, error) {
			if proxy != "" {
				return url.Parse(proxy)
			}
			return nil, nil
		},
	}
	client := &http.Client{
		Timeout:   5 * time.Second, //设置超时时间为5秒
		Transport: tr,              //跳过https证书认证
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func Postreq(target string, headers map[string]string, body string, proxy string, timeOutMap ...int) (statusCode int, respHeader http.Header, respBody []byte, err error) {
	timeOut := 5 //设置超时时间默认5s
	if len(timeOutMap) != 0 {
		for _, item := range timeOutMap {
			timeOut = item
		}
	}

	//创建请求
	req, err := http.NewRequest("POST", target, bytes.NewBufferString(body))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("[-]Postreq创建请求出现错误:%v", err)
	}

	//设置请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 创建自定义的Transport
	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if proxy != "" {
				return url.Parse(proxy)
			}
			return nil, nil
		},
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true, //禁用压缩
		DisableKeepAlives:  true, //关闭长连接
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁用重定向
			return http.ErrUseLastResponse
		},
		Timeout:   time.Duration(timeOut) * time.Second, //设置超时时间
		Transport: tr,                                   //跳过https证书认证
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	Body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("[-]getResq读取响应时发送错误：%v\n", err)
	}
	return resp.StatusCode, resp.Header, Body, nil
}

func GetResq(target string, headers map[string]string, proxy string) (statusCode int, respHeader http.Header, respBody []byte, err error) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("[-]getResq请求创建出错:%v", err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// 创建自定义的Transport
	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if proxy != "" {
				return url.Parse(proxy)
			}
			return nil, nil
		},
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true, //禁用压缩
		DisableKeepAlives:  true, //关闭长连接
	}
	client := &http.Client{
		Timeout:   5 * time.Second, //设置超时时间为5秒
		Transport: tr,              //跳过https证书认证
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	Body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("[-]getResq读取响应时发送错误：%v\n", err)
	}
	return resp.StatusCode, resp.Header, Body, nil
}
