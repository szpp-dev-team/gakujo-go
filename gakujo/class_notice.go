package gakujo

import (
	"bytes"
	"io"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) ClassNotice() ( /*string*/ []model.ClassNoticeRow, error) {
	body, err := c.fetchClassNoticeRow()
	if err != nil {
		return /*" "*/ nil, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	classNoticeRow, err := scrape.ClassNoticeRow(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return /*" "*/ nil, err
	}
	return /*classNoticeRow*/ classNoticeRow, err
}

func (c *Client) fetchClassNoticeRow() (io.ReadCloser, error) {
	params := make(url.Values)
	params.Set("headTitle", "授業連絡一覧")
	params.Set("menuCode", "A01")
	params.Set("nextPath", "/classcontact/classContactList/initialize")

	return c.getPage(GeneralPurposeUrl, params)
}
