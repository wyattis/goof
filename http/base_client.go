package http

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewBaseClient(base string, client *http.Client) *BaseClient {
	b := &BaseClient{
		Base:   base,
		Client: client,
	}
	if client == nil {
		b.Client = http.DefaultClient
	}
	return b
}

type BaseClient struct {
	Base   string
	Header http.Header
	*http.Client
}

func (c *BaseClient) fixPath(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return c.Base + path
}

func (c *BaseClient) Get(path string) (res *http.Response, err error) {
	return c.Client.Get(c.fixPath(path))
}

func (c *BaseClient) Head(path string) (res *http.Response, err error) {
	return c.Client.Head(c.fixPath(path))
}

func (c *BaseClient) Post(path, contentType string, body io.Reader) (res *http.Response, err error) {
	return c.Client.Post(c.fixPath(path), contentType, body)
}

func (c *BaseClient) PostForm(path string, data url.Values) (res *http.Response, err error) {
	return c.Client.PostForm(c.fixPath(path), data)
}

func (c *BaseClient) Do(req *http.Request) (res *http.Response, err error) {
	uri := req.URL.String()
	req.URL, err = url.Parse(c.fixPath(uri))
	if err != nil {
		return
	}
	return c.Client.Do(req)
}
