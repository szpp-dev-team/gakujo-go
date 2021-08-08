package gakujo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) Login(username, password string) error {
	if err := c.fetchGakujoPortalJSESSIONID(); err != nil {
		return err
	}

	if err := c.fetchGakujoRootJSESSIONID(); err != nil {
		return err
	}

	if err := c.preLogin(); err != nil {
		return err
	}
	resp, err := c.shibbolethlogin()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return fmt.Errorf("Response status was %d(expect %d or %d)", resp.StatusCode, http.StatusOK, http.StatusFound)
	}

	// セッションがないとき
	if resp.StatusCode == http.StatusFound {
		loginAPIurl, err := c.fetchLoginAPIurl(resp.Header.Get("Location"))
		if err != nil {
			return err
		}
		if err := c.login(IdpHostName+loginAPIurl, username, password); err != nil {
			return err
		}
	}

	return c.initialize()
}

func (c *Client) fetchGakujoPortalJSESSIONID() error {
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/portal/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) fetchGakujoRootJSESSIONID() error {
	unixmilli := time.Now().UnixNano() / 1000000
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/UI/jsp/topPage/topPage.jsp?_=" + strconv.FormatInt(unixmilli, 10))
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) preLogin() error {
	datas := url.Values{}
	datas.Set("mistakeChecker", "0")

	resp, err := c.client.PostForm("https://gakujo.shizuoka.ac.jp/portal/login/preLogin/preLogin", datas)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) fetchLoginAPIurl(SSOSAMLRequestURL string) (string, error) {
	resp, err := c.get(SSOSAMLRequestURL)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}
	return resp.Header.Get("Location"), nil
}

func (c *Client) login(reqUrl, username, password string) error {
	htmlReadCloser, err := c.postSSOexecution(reqUrl, username, password)
	if err != nil {
		return err
	}
	relayState, samlResponse, err := scrape.RelayStateAndSAMLResponse(htmlReadCloser)
	if err != nil {
		return err
	}
	htmlReadCloser.Close()
	_, _ = io.Copy(io.Discard, htmlReadCloser)

	location, err := c.fetchSSOinitLoginLocation(relayState, samlResponse)
	if err != nil {
		return err
	}

	resp, err := c.getWithReferer(location, "https://idp.shizuoka.ac.jp/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c *Client) postSSOexecution(reqUrl, username, password string) (io.ReadCloser, error) {
	datas := make(url.Values)
	datas.Set("j_username", username)
	datas.Set("j_password", password)
	datas.Set("_eventId_proceed", "")

	resp, err := c.client.PostForm(reqUrl, datas)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return resp.Body, nil
}

func (c *Client) shibbolethlogin() (*http.Response, error) {
	url := HostName + "/portal/shibbolethlogin/shibbolethLogin/initLogin/sso"
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	resp, err := c.request(req)
	resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp, err
}

func (c *Client) fetchSSOinitLoginLocation(relayState, samlResponse string) (string, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/Shibboleth.sso/SAML2/POST"

	datas := make(url.Values)
	datas.Set("RelayState", relayState)
	datas.Set("SAMLResponse", samlResponse)

	resp, err := c.client.PostForm(reqUrl, datas)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s\nResponse status was %d(expect %d)", reqUrl, resp.StatusCode, http.StatusFound)
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) initialize() error {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/home/home/initialize"

	datas := make(url.Values)
	datas.Set("EXCLUDE_SET", "")

	rc, err := c.getPage(reqURL, datas)
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
		_, _ = io.Copy(io.Discard, rc)
	}()

	return nil
}
