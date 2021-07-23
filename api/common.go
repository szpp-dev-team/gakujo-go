package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/szpp-dev-team/gakujo-api/scrape"
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
	b, _ := io.ReadAll(resp.Body)

	// これクソコード
	token, err := scrape.ApacheToken(io.NopCloser(bytes.NewReader(b)))
	if err != nil {
		switch err.(type) {
		case *scrape.ErrorNotFound:
			fmt.Fprintln(os.Stderr, err.Error())
		default:
			return nil, err
		}
	}
	c.token = token

	if strings.Contains(string(b), "不正な操作") {
		return nil, errors.New("不正な操作が行われました")
	}

	resp.Body = io.NopCloser(bytes.NewReader(b))
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

// http.PostForm wrapper
func (c *Client) postForm(url string, datas url.Values) (*http.Response, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(datas.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func saveCookies(url *url.URL, cookies *[]http.Cookie) error {
	file, err := os.Create("cookies")
	if err != nil {
		return err
	}
	for _, cookie := range *cookies {
		_, _ = file.WriteString(fmt.Sprintf("%v=%v\n", cookie.Name, cookie.Value))
	}
	return nil
}
