package gakujo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"
)

type CookieJarJSON struct {
	Url     UrlJSON      `json:"url"`
	Cookies []CookieJSON `json:"cookies"`
}

type UrlJSON struct {
	Scheme     string `json:"scheme"`
	Host       string `json:"host"`
	Path       string `json:"path"`
	ForceQuery bool   `json:"force_query"`
}

type CookieJSON struct {
	Name    string    `json:"name"`
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
}

// return url.URL
func (c *CookieJarJSON) urlURL() *url.URL {
	return &url.URL{
		Scheme:     c.Url.Scheme,
		Host:       c.Url.Host,
		Path:       c.Url.Path,
		ForceQuery: c.Url.ForceQuery,
	}
}

// return http.Cookie
func (c *CookieJarJSON) httpCookies() []*http.Cookie {
	cookies := []*http.Cookie{}
	for _, cookie := range c.Cookies {
		cookies = append(
			cookies,
			&http.Cookie{
				Name:    cookie.Name,
				Value:   cookie.Value,
				Expires: cookie.Expires,
			},
		)
	}
	return cookies
}

// dump cookies to cookies.json
func (c *Client) DumpCookies() error {
	file, err := os.OpenFile("cookies.json", os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.Marshal(c.cookieJars)
	if err != nil {
		return err
	}
	if _, err := file.Write(b); err != nil {
		return err
	}

	return nil
}

func (c *Client) saveCookies(url *url.URL, cookies []*http.Cookie) error {
	urlJson := UrlJSON{
		Scheme:     url.Scheme,
		Host:       url.Host,
		Path:       url.Path,
		ForceQuery: url.ForceQuery,
	}
	cookieJsons := []CookieJSON{}
	for _, cookie := range cookies {
		cookieJson := CookieJSON{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Expires: cookie.Expires,
		}
		cookieJsons = append(cookieJsons, cookieJson)
	}
	c.cookieJars = append(
		c.cookieJars,
		CookieJarJSON{
			Url:     urlJson,
			Cookies: cookieJsons,
		},
	)
	return nil
}

// load cookies from "cookies.json" and set cookies to c.client.Jar
// please call this first if you use
func (c *Client) LoadCookies() error {
	b, err := os.ReadFile("cookies.json") // ファイルの場所は要検討
	if err != nil {
		return err
	}
	cookieJarJsons := []CookieJarJSON{}
	if err := json.Unmarshal(b, &cookieJarJsons); err != nil {
		return err
	}

	for _, cookieJarJson := range cookieJarJsons {
		c.client.Jar.SetCookies(
			cookieJarJson.urlURL(),
			cookieJarJson.httpCookies(),
		)
	}

	return nil
}
