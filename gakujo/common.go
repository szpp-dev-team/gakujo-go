package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/szpp-dev-team/gakujo-go/scrape"
)

const (
	HostName          = "https://gakujo.shizuoka.ac.jp"
	IdpHostName       = "https://idp.shizuoka.ac.jp"
	GeneralPurposeUrl = "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
)

type Client struct {
	client *http.Client
	jar    *cookiejar.Jar
	token  string // org.apache.struts.taglib.html.TOKEN
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     jar,
		Timeout: 5 * time.Minute,
	}
	return &Client{
		client: &httpClient,
		jar:    jar,
	}
}

// search a cookie "JSESSIONID" from c.jar
// if not found, return ""
func (c *Client) SessionID() string {
	u, _ := url.Parse(HostName)
	for _, cookie := range c.jar.Cookies(u) {
		if cookie.Name == "JSESSIONID" {
			return cookie.Value
		}
	}

	return ""
}

// 任意のページを取得し、apache Token を取得する
func (c *Client) GetPage(url string, data url.Values) ([]byte, error) {
	data.Set("org.apache.struts.taglib.html.TOKEN", c.token)
	resp, err := c.postForm(url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status was %d(expext %d)", resp.StatusCode, http.StatusOK)
	}
	page := &bytes.Buffer{}
	r := io.TeeReader(resp.Body, page)
	token, err := scrape.ApacheToken(r)
	if err != nil {
		return nil, err // getPage では必ず apache Token が含まれるページを取得するはず
	}
	c.token = token // トークンを更新
	return page.Bytes(), nil
}

// http.Get wrapper
func (c *Client) get(url string, param ...url.Values) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if len(param) > 0 {
		req.URL.RawQuery = param[0].Encode()
	}
	return c.client.Do(req)
}

// http.PostForm wrapper
func (c *Client) postForm(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.client.Do(req)
}

type UnexpectedStatusCodeError struct {
	expected int
	actual   int
}

func (e *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("unexpected status code: expected %d, actual %d", e.expected, e.actual)
}

func unexpectedStatusCodeError(expected, actual int) error {
	return &UnexpectedStatusCodeError{expected, actual}
}
