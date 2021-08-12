package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
	"github.com/szpp-dev-team/gakujo-api/util"
)

type KyoumuClient struct {
	c *Client
}

func (c *Client) NewKyoumuClient() (*KyoumuClient, error) {
	if err := c.initializeShibboleth("kyoumu"); err != nil {
		return nil, err
	}
	if err := c.kyoumuLoginStudent(); err != nil {
		return nil, err
	}
	return &KyoumuClient{
		c: c,
	}, nil
}

func (kc *KyoumuClient) SeisekiRows() ([]*model.SeisekiRow, error) {
	b, err := kc.c.fetchSeisekiHtml()
	if err != nil {
		return nil, err
	}
	return scrape.SeisekiRows(bytes.NewReader(b))
}

func (kc *KyoumuClient) DepartmentGpa() (*model.DepartmentGpa, error) {
	if _, err := kc.c.fetchSeisekiHtml(); err != nil {
		return nil, err
	}
	b, err := kc.c.fetchDepartmentGpaHtml()
	if err != nil {
		return nil, err
	}
	return scrape.DepartmentGpa(bytes.NewBuffer(b))
}

func (c *Client) initializeShibboleth(renkeiType string) error {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/home/systemCooperationLink/initializeShibboleth?renkeiType=" + renkeiType
	resp, err := c.postForm(
		reqUrl,
		url.Values{"org.apache.struts.taglib.html.TOKEN": []string{c.token}},
	)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return util.RespStatusError(resp.StatusCode, http.StatusOK)
	}
	return nil
}

func (c *Client) kyoumuLoginStudent() error {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/sso/loginStudent.do"
	resp, err := c.postForm(reqURL, url.Values{"loginID": {""}})
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return util.RespStatusError(resp.StatusCode, http.StatusOK)
	}
	return nil
}

func (c *Client) fetchSeisekiHtml() ([]byte, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/seisekiSearchStudentInit.do;"

	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	q := req.URL.Query()
	q.Set("jsessionid", c.SessionID())
	q.Set("mainMenuCode", "008")
	q.Set("parentMenuCode", "007")
	req.URL.RawQuery = q.Encode()

	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expected %d)", resp.StatusCode, http.StatusOK)
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	return io.ReadAll(resp.Body)
}

func (c *Client) fetchDepartmentGpaHtml() ([]byte, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/departmentGpa.do;"

	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	q := req.URL.Query()
	q.Set("jsessionid", c.SessionID())
	req.URL.RawQuery = q.Encode()

	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expected %d)", resp.StatusCode, http.StatusOK)
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	return io.ReadAll(resp.Body)
}
