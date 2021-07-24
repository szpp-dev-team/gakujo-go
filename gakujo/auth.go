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
	loginAPIurl, err := c.fetchLoginAPIurl()
	if err != nil {
		return err
	}

	if err := c.login(IdpHostName+loginAPIurl, username, password); err != nil {
		return err
	}

	return nil
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
	params := url.Values{}
	params.Set("mistakeChecker", "0")

	resp, err := c.client.PostForm("https://gakujo.shizuoka.ac.jp/portal/login/preLogin/preLogin", params)
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

func (c *Client) fetchLoginAPIurl() (string, error) {
	url, err := c.fetchSSOSAMLRequestLocation()
	if err != nil {
		return "", err
	}
	resp, err := c.get(url)
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

	return c.initialize()
}

func (c *Client) postSSOexecution(reqUrl, username, password string) (io.ReadCloser, error) {
	params := make(url.Values)
	params.Set("j_username", username)
	params.Set("j_password", password)
	params.Set("_eventId_proceed", "")

	resp, err := c.client.PostForm(reqUrl, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return resp.Body, nil
}

func (c *Client) fetchSSOSAMLRequestLocation() (string, error) {
	url := HostName + "/portal/shibbolethlogin/shibbolethLogin/initLogin/sso"
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	resp, err := c.request(req)
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

func (c *Client) fetchSSOinitLoginLocation(relayState, samlResponse string) (string, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/Shibboleth.sso/SAML2/POST"

	params := make(url.Values)
	params.Set("RelayState", relayState)
	params.Set("SAMLResponse", samlResponse)

	resp, err := c.client.PostForm(reqUrl, params)
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

	params := make(url.Values)
	params.Set("EXCLUDE_SET", "")

	resp, err := c.postForm(reqURL, params)
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
