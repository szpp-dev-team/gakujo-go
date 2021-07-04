package api

import (
	"net/http"
	"net/http/cookiejar"
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
	jar, _ := cookiejar.New(nil)
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
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
	c.jar.SetCookies(resp.Request.URL, resp.Cookies())

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
