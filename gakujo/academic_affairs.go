package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

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

type OverCapasityError struct {
	message string
}

func (e OverCapasityError) Error() string {
	return e.message
}

func (kc *KyoumuClient) PostRishuRegistration(formData *model.PostKamokuFormData) error {
	// よくわからんけど登録 api までに経由するであろうページにアクセスしないと500 になる
	if _, err := kc.c.fetchRishuInitHtml(); err != nil {
		return err
	}
	if _, err := kc.c.fetchSearchKamokuInitHtml(formData.Youbi, formData.Jigen); err != nil {
		return err
	}
	reqUrl := "https://gakujo.shizuoka.ac.jp/kyoumu/searchKamoku.do"
	data := formData.FormData()
	data.Set("button_kind.registKamoku.x", "16")
	data.Set("button_kind.registKamoku.y", "10")
	resp, err := kc.c.postForm(reqUrl, data)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return util.RespStatusError(resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return errors.WithStack(err)
	}
	errMsg := doc.Find("font.txt12:nth-child(8) > ul:nth-child(1) > li:nth-child(1)").Text()
	if errMsg == "定員数を超えているため、登録できません。" {
		return errors.WithStack(OverCapasityError{errMsg})
	}
	if errMsg != "" {
		return errors.WithStack(fmt.Errorf("unexpected error: %s", errMsg))
	}
	return nil
}

func (kc *KyoumuClient) GetRishuuInit() error {
	b, err := kc.c.fetchRishuInitHtml()
	if err != nil {
		return err
	}
	fmt.Println(string(b))
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

func (c *Client) fetchRishuInitHtml() ([]byte, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/rishuuInit.do"

	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	q := req.URL.Query()
	q.Set("jsessionid", c.SessionID())
	q.Set("mainMenuCode", "002")
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

func (c *Client) fetchSearchKamokuInitHtml(youbi, jigen int) ([]byte, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/kyoumu/searchKamokuInit.do"

	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	q := req.URL.Query()
	q.Set("jsessionid", c.SessionID())
	q.Set("youbi", strconv.Itoa(youbi))
	q.Set("jigen", strconv.Itoa(jigen))
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
