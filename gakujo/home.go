package gakujo

import (
	"bytes"
	"io"
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
	datas := make(url.Values)
	datas.Set("headTitle", "ホーム")
	datas.Set("menuCode", "Z07") // TODO: menucode を定数化(まとめてやる)
	datas.Set("nextPath", "/home/home/initialize")

	return c.getPage(GeneralPurposeUrl, datas)
}

func (c *Client) NoticeDetail() (*model.NoticeDetail, error) {
	body, err := c.fetchNoiceDetailhtml()
	if err != nil {
		return &model.NoticeDetail{}, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	noticeDetail, err := scrape.NoticeDetail(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return &model.NoticeDetail{}, err
	}
	return noticeDetail, nil
}

func (c *Client) fetchNoiceDetailhtml() (io.ReadCloser, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/portaltopcommon/newsForTop/deadLineForTop"

	params := make(url.Values)
	params.Set("newsTargetIndexNo", "20")

	return c.getPage(reqURL, params)
}

func (c *Client) ClassNotice() (string /**model.ClassNoticeRow*/, error) {
	body, err := c.fetchClassNoticeRow()
	if err != nil {
		return " ", err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	//classNoticeRow,err := scrape.ClassNoticeRow(io.NopCloser(bytes.NewBuffer(b)))
	return string(b), err
}

func (c *Client) fetchClassNoticeRow() (io.ReadCloser, error) {
	params := make(url.Values)
	params.Set("headTitle:", "授業連絡一覧")
	params.Set("menuCode:", "A01")
	params.Set("nextPath:", "/classcontact/classContactList/initialize")

	return c.getPage(GeneralPurposeUrl, params)
}
