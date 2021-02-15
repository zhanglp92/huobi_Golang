package internal

import (
	"crypto/tls"
	"github.com/huobirdcenter/huobi_golang/logging/perflogger"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HttpGet(url string) (string, error) {
	logger := perflogger.GetInstance()
	logger.Start()

	defer func() { logger.StopAndLog("GET", url) }()
	return httpRequestHandle.get(url)
}

func HttpPost(url string, body string) (string, error) {
	logger := perflogger.GetInstance()
	logger.Start()

	defer func() { logger.StopAndLog("POST", url) }()
	return httpRequestHandle.post(url, body)
}

var httpRequestHandle = &httpRequest{}

type httpRequest struct {
}

func (m *httpRequest) post(url, body string) (string, error) {
	if cli, err := m.client(); err != nil {
		return "", err
	} else if resp, err := cli.Post(url, "application/json", strings.NewReader(body)); err != nil {
		return "", err
	} else {
		defer func() { _ = resp.Body.Close() }()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return "", err
		} else {
			return string(body), nil
		}
	}
}

func (m *httpRequest) get(url string) (string, error) {
	if cli, err := m.client(); err != nil {
		return "", err
	} else if resp, err := cli.Get(url); err != nil {
		return "", err
	} else {
		defer func() { _ = resp.Body.Close() }()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return "", err
		} else {
			return string(body), nil
		}
	}
}

func (m *httpRequest) client() (*http.Client, error) {
	if proxy, err := url.Parse(m.getProxy()); err != nil {
		return nil, err
	} else {
		return &http.Client{
			Transport: &http.Transport{
				Proxy:           http.ProxyURL(proxy),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}, nil
	}
}

func (m *httpRequest) getProxy() string {
	return os.Getenv("http_proxy")
}
