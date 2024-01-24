package http

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/url"
)

func New() Http {
	return &implement{}
}

type implement struct {
	headers            map[string]string
	body               []byte
	proxyUrl           string
	insecureSkipVerify bool
	basicAuth          struct {
		enabled  bool
		username string
		password string
	}
}

func (h *implement) clone() *implement {
	headers := make(map[string]string)
	for k, v := range h.headers {
		headers[k] = v
	}

	var body []byte
	if h.body != nil {
		body = make([]byte, len(h.body))
		copy(body, h.body)
	} else {
		body = nil
	}

	return &implement{
		headers:            headers,
		body:               body,
		proxyUrl:           h.proxyUrl,
		insecureSkipVerify: h.insecureSkipVerify,
		basicAuth:          h.basicAuth,
	}
}

func (h *implement) WithHeaders(headers map[string]string) Http {
	httpClone := h.clone()
	httpClone.headers = headers
	return httpClone
}

func (h *implement) WithProxy(url string) Http {
	httpClone := h.clone()
	httpClone.proxyUrl = url
	return httpClone
}

func (h *implement) SetBasicAuth(username string, password string) Http {
	httpClone := h.clone()
	httpClone.basicAuth.enabled = true
	httpClone.basicAuth.username = username
	httpClone.basicAuth.password = password
	return httpClone
}

func (h *implement) WithBody(body []byte) Http {
	httpClone := h.clone()
	httpClone.body = body
	return httpClone
}

func (h *implement) InsecureSkipVerify() Http {
	httpClone := h.clone()
	httpClone.insecureSkipVerify = true
	return httpClone
}

func (h *implement) Connect(url string) (*http.Response, error) {
	return h.do(http.MethodConnect, url)
}

func (h *implement) Delete(url string) (*http.Response, error) {
	return h.do(http.MethodDelete, url)
}

func (h *implement) Get(url string) (*http.Response, error) {
	return h.do(http.MethodGet, url)
}

func (h *implement) Head(url string) (*http.Response, error) {
	return h.do(http.MethodHead, url)
}

func (h *implement) Options(url string) (*http.Response, error) {
	return h.do(http.MethodOptions, url)
}

func (h *implement) Patch(url string) (*http.Response, error) {
	return h.do(http.MethodPatch, url)
}

func (h *implement) Post(url string) (*http.Response, error) {
	return h.do(http.MethodPost, url)
}

func (h *implement) Put(url string) (*http.Response, error) {
	return h.do(http.MethodPut, url)
}

func (h *implement) Trace(url string) (*http.Response, error) {
	return h.do(http.MethodTrace, url)
}

func (h *implement) do(method string, u string) (*http.Response, error) {
	req, err := http.NewRequest(method, u, bytes.NewBuffer(h.body))
	if err != nil {
		return nil, err
	}

	var proxyFunc func(*http.Request) (*url.URL, error) = nil
	if h.proxyUrl != "" {
		proxy, err := url.Parse(h.proxyUrl)
		if err != nil {
			return nil, err
		}
		proxyFunc = http.ProxyURL(proxy)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: h.insecureSkipVerify,
			},
			Proxy: proxyFunc,
		},
	}

	for k, v := range h.headers {
		req.Header.Add(k, v)
	}

	// close connection immediately
	req.Close = true

	if h.basicAuth.enabled {
		req.SetBasicAuth(h.basicAuth.username, h.basicAuth.password)
	}

	var rsp *http.Response
	rsp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
