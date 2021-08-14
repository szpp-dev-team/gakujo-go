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

func (c *Client) NoticeDetail(index string) (*model.NoticeDetail, error) {
	body, err := c.fetchNoiceDetailhtml(index)
	if err != nil {
		return nil, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	noticeDetail, err := scrape.NoticeDetail(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return nil, err
	}
	return noticeDetail, nil
}

func (c *Client) fetchNoiceDetailhtml(index string) (io.ReadCloser, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/classcontact/classContactList/goDetail/" + index
	data := make(url.Values)
	data.Set("headTitle", "授業連絡一覧")
	data.Set("menuCode", "A01")
	data.Set("nextPath", "/classcontact/classContactList/initialize")
	_, err := c.getPage(GeneralPurposeUrl, data)
	if err != nil {
		return nil, err
	}
	data.Set("tbl_A01_01_length", "10")
	data.Set("schoolYear", "2021")
	data.Set("semesterCode", "1")
	data.Set("reportDateStart", "2021/02/01")
	return c.getPage(reqUrl, data)
}
