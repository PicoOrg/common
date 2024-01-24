package http

import (
	"net/http"
)

type Http interface {
	InsecureSkipVerify() (http Http)
	WithBody(body []byte) (http Http)
	WithHeaders(headers map[string]string) (http Http)
	WithProxy(url string) (http Http)
	SetBasicAuth(username string, password string) (http Http)

	Connect(url string) (rsp *http.Response, err error)
	Delete(url string) (rsp *http.Response, err error)
	Get(url string) (rsp *http.Response, err error)
	Head(url string) (rsp *http.Response, err error)
	Options(url string) (rsp *http.Response, err error)
	Patch(url string) (rsp *http.Response, err error)
	Post(url string) (rsp *http.Response, err error)
	Put(url string) (rsp *http.Response, err error)
	Trace(url string) (rsp *http.Response, err error)
}
