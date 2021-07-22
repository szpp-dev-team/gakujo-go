package api

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strings"
	"time"
)

const (
	HostName    = "https://gakujo.shizuoka.ac.jp"
	IdpHostName = "https://idp.shizuoka.ac.jp"
)

type Client struct {
	client *http.Client
	jar    *cookiejar.Jar
}

func NewClient() *Client {
	jar, _ := cookiejar.New(
		nil,
	)
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     jar,
		Timeout: 30 * time.Second,
	}
	return &Client{
		client: &httpClient,
		jar:    jar,
	}
}

// save cookie "Set-Cookies" into client.cookie
func (c *Client) request(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	b, _ := httputil.DumpResponse(resp, true)
	if strings.Contains(string(b), "不正な操作") {
		return nil, errors.New("不正な操作が行われました")
	}

	return resp, nil
}

// http.Get wrapper
func (c *Client) get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return c.request(req)
}

// http.Get wrapper
func (c *Client) getWithReferer(url, referer string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Referer", referer)
	return c.request(req)
}
