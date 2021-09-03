package gakujo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

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

func (kc *KyoumuClient) ChusenRegistrationRows() ([]*model.ChusenRegistrationRow, error) {
	b, err := kc.c.fetchChusenRegistrationHtml()
	if err != nil {
		return nil, err
	}
	return scrape.ChusenRegistrationRows(bytes.NewReader(b))
}

func (kc *KyoumuClient) PostChusenRegistration(chusenRows []*model.ChusenRegistrationRow) error {
	data := url.Values{}
	data.Set("x", strconv.Itoa(rand.Int()%100))
	data.Set("y", strconv.Itoa(rand.Int()%100))
	data.Set("RisshuForm.jikanwariVector", "AA")

	for i, chusenRow := range chusenRows {
		if chusenRow == nil {
			return fmt.Errorf("chusenRow(index: %d) is nil", i)
		}
		if chusenRow.ChoiceRank < 0 || 3 < chusenRow.ChoiceRank {
			return errors.New("choice rank is out of range")
		}
		data.Set(chusenRow.AttrName, strconv.Itoa(chusenRow.ChoiceRank))
	}

	resp, err := kc.c.postForm("https://gakujo.shizuoka.ac.jp/kyoumu/chuusenRishuuRegist.do", data)
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

func (c *Client) fetchChusenRegistrationHtml() ([]byte, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/chuusenRishuuInit.do;"

	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	q := req.URL.Query()
	q.Set("jsessionid", c.SessionID())
	q.Set("mainMenuCode", "019")
	q.Set("parentMenuCode", "001")
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
