package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Login(username, password string) error {
	if err := c.fetchGakujoJSESSIONID(); err != nil {
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

func (c *Client) fetchGakujoJSESSIONID() error {
	resp, err := c.get("https://gakujo.shizuoka.ac.jp/portal/")
	if err != nil {
		return err
	}
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
	fmt.Println("Redilect to", url)
	resp, err := c.get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}
	return resp.Header.Get("Location"), nil
}

func (c *Client) login(url, username, password string) error {
	htmlReadCloser, err := c.postSSOexecution(url, username, password)
	if err != nil {
		return err
	}
	relayState, samlResponse, err := scrapRelayStateAndSAMLResponse(htmlReadCloser)
	if err != nil {
		return err
	}

	location, err := c.fetchSSOinitLoginLocation(relayState, samlResponse)
	if err != nil {
		return err
	}
	fmt.Println("Redilect to:", location)
	resp, err := c.get(location)
	if err != nil {
		return err
	}
	fmt.Println(resp.Request)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}
	b, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(b), "不正な操作") {
		return errors.New("は？")
	}

	return nil
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
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s\nResponse status was %d(expect %d)", reqUrl, resp.StatusCode, http.StatusFound)
	}

	return resp.Header.Get("Location"), nil
}

func scrapRelayStateAndSAMLResponse(htmlReader io.ReadCloser) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", "", err
	}
	selection := doc.Find("html > body > form > div > input")
	html, _ := selection.Html()
	fmt.Println(html)
	relayState, ok := selection.Attr("value")
	if !ok {
		return "", "", errors.New("RelayState was not found")
	}
	selection = selection.Next()
	samlResponse, ok := selection.Attr("value")
	if !ok {
		return "", "", errors.New("SAMLResponse was not found")
	}

	return relayState, samlResponse, nil
}
