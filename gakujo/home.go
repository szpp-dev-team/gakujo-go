package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) Home() (model.HomeInfo, error) {
	body, err := c.fetchNoiceDetailhtml()
	if err != nil {
		return model.HomeInfo{}, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	taskRows, err := scrape.TaskRows(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return model.HomeInfo{}, err
	}
	noticeRows, err := scrape.NoticeRows(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return model.HomeInfo{}, err
	}
	return model.HomeInfo{
		TaskRows:   taskRows,
		NoticeRows: noticeRows,
	}, nil
}

func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {
	reqURL := GeneralPurposeUrl

	params := make(url.Values)
	params.Set("org.apache.struts.taglib.html.TOKEN", c.token)
	params.Set("headTitle", "ホーム")
	params.Set("menuCode", "Z07") // TODO: 定数化(まとめてやる)
	params.Set("nextPath", "/home/home/initialize")

	resp, err := c.postForm(reqURL, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return resp.Body, nil
}

func (c *Client) fetchNoiceDetailhtml() (io.ReadCloser, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/portaltopcommon/newsForTop/deadLineForTop"

	params := make(url.Values)
	params.Set("org.apache.struts.taglib.html.TOKEN", c.token)
	params.Set("newsTargetIndexNo", "0")

	resp, err := c.postForm(reqURL, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expext %d)", resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))
	return resp.Body, nil
}
