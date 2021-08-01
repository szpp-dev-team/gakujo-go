package gakujo

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
		Timeout: 5 * time.Minute,
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

	// validation
	if strings.Contains(string(b), "不正な操作") {
		return nil, errors.New("不正な操作が行われました")
	}

	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}

// get page which needs org.apache.struts.taglib.html.TOKEN and save its token
func (c *Client) getPage(url string, datas url.Values) (io.ReadCloser, error) {
	datas.Set("org.apache.struts.taglib.html.TOKEN", c.token)

	resp, err := c.postForm(url, datas)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expext %d)", resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token, err := scrape.ApacheToken(io.NopCloser(bytes.NewReader(b)))
	if err != nil {
		// getPage では必ず apache Token が含まれるページを取得するはず
		return nil, err
	}
	c.token = token

	return io.NopCloser(bytes.NewReader(b)), nil
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
