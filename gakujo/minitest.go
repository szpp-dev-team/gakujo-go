package gakujo

import (
	"bytes"
	"io"
	"net/url"
	"strconv"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) MinitestRows(option *model.MinitestSearchOption) ([]model.MinitestRow, error) {
	page, err := c.fetchMinitestRowsPage(option)
	if err != nil {
		return nil, err
	}
	return scrape.MinitestRows(io.NopCloser(bytes.NewReader(page)))
}

func (c *Client) fetchMinitestRowsPage(option *model.MinitestSearchOption) ([]byte, error) {
	if _, err := c.fetchGeneralPurposeMinitestRowsPage(); err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("schoolYear", strconv.Itoa(option.SchoolYear))
	data.Set("semesterCode", strconv.Itoa(option.SemesterCode.Int()))
	body, err := c.getPage("https://gakujo.shizuoka.ac.jp/portal/test/student/searchList/search", data)
	if err != nil {
		return nil, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	return io.ReadAll(body)
}

func (c *Client) fetchGeneralPurposeMinitestRowsPage() ([]byte, error) {
	data := url.Values{}
	data.Set("headTitle", "小テスト一覧")
	data.Set("menuCode", "A03")
	data.Set("nextPath", "/test/student/searchList/initialize")
	page, err := c.getPage("https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/", data)
	if err != nil {
		return nil, err
	}
	defer func() {
		page.Close()
		_, _ = io.Copy(io.Discard, page)
	}()
	return io.ReadAll(page)
}
