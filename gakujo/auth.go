package gakujo

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/szpp-dev-team/gakujo-go/scrape"
)

func (c *Client) Login(username, password string) error {
	if err := c.fetchGakujoPortalJSESSIONID(); err != nil {
		return err
	}

	if err := c.getTopPageJsp(); err != nil {
		return err
	}
	if err := c.postPreLogin(); err != nil {
		return err
	}
	loc, err := c.shibbolethlogin()
	if err != nil {
		return err
	}

	// セッションがないとき
	if loc != "" {
		loginAPIurl, err := c.fetchLoginAPIurl(loc)
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
		return unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}
	return nil
}

func (c *Client) getTopPageJsp() error {
	unixmilli := time.Now().UnixNano() / int64(time.Millisecond)
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/UI/jsp/topPage/topPage.jsp?_=" + strconv.FormatInt(unixmilli, 10))
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}
	return nil
}

func (c *Client) postPreLogin() error {
	data := url.Values{}
	data.Set("mistakeChecker", "0")
	resp, err := c.client.PostForm("https://gakujo.shizuoka.ac.jp/portal/login/preLogin/preLogin", data)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
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
		return "", unexpectedStatusCodeError(http.StatusFound, resp.StatusCode)
	}
	return resp.Header.Get("Location"), nil
}

func (c *Client) login(reqUrl, username, password string) error {
	rc, err := c.postSSOexecution(reqUrl, username, password)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, rc)
		rc.Close()
	}()

	relayState, samlResponse, err := scrape.RelayStateAndSAMLResponse(rc)
	if err != nil {
		return err
	}
	location, err := c.postShibbolethSAML2POST(relayState, samlResponse)
	if err != nil {
		return err
	}
	resp, err := c.get(location, http.Header{"Referer": {"https://idp.shizuoka.ac.jp/"}})
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}

	if err := c.postAccessEnvironmentRegist(); err != nil {
		return err
	}

	c.initialize()

	return nil
}

func (c *Client) postSSOexecution(uri, username, password string) (io.ReadCloser, error) {
	data := make(url.Values)
	data.Set("j_username", username)
	data.Set("j_password", password)
	data.Set("_eventId_proceed", "")
	resp, err := c.client.PostForm(uri, data)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}
	return resp.Body, nil
}

func (c *Client) shibbolethlogin() (string, error) {
	req, _ := http.NewRequest(http.MethodPost, HostName+"/portal/shibbolethlogin/shibbolethLogin/initLogin/sso", nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return "", unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}
	return resp.Header.Get("Location"), nil
}

func (c *Client) postShibbolethSAML2POST(relayState, samlResponse string) (string, error) {
	data := make(url.Values)
	data.Set("RelayState", relayState)
	data.Set("SAMLResponse", samlResponse)
	resp, err := c.client.PostForm("https://gakujo.shizuoka.ac.jp/Shibboleth.sso/SAML2/POST", data)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", unexpectedStatusCodeError(http.StatusFound, resp.StatusCode)
	}
	return resp.Header.Get("Location"), nil
}

// アクセス環境登録
func (c *Client) postAccessEnvironmentRegist() error {
	data := make(url.Values)
	data.Set("accessEnvName", uuid.NewString())
	data.Set("checkBox", "on")
	data.Set("newAccessKey", "")
	resp, err := c.postForm("https://gakujo.shizuoka.ac.jp/portal/common/accessEnvironmentRegist/goHome/", data)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusCodeError(http.StatusOK, resp.StatusCode)
	}
	return nil
}

func (c *Client) initialize() error {
	data := make(url.Values)
	data.Set("EXCLUDE_SET", "")
	if _, err := c.GetPage("https://gakujo.shizuoka.ac.jp/portal/home/home/initialize", data); err != nil {
		return err
	}
	return nil
}
