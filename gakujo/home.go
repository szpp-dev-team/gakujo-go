package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) Home() (model.HomeInfo, error) {
	body, err := c.fetchHomeHtml()
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
