package gakujo

import (
	"fmt"
	"io"
	"net/url"

	"github.com/szpp-dev-team/gakujo-go/model"
	"github.com/szpp-dev-team/gakujo-go/scrape"
)

func (c *Client) ClassNoticeRows(opt *model.ClassNoticeSearchOption) ([]model.ClassNoticeRow, error) {
	body, err := c.fetchClassNoticeSearchPage(opt)
	if err != nil {
		return nil, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	return scrape.ClassNoticeRows(body)
}

func (c *Client) ClassNoticeDetail(row *model.ClassNoticeRow, opt *model.ClassNoticeSearchOption) (*model.ClassNoticeDetail, error) {
	body, err := c.fetchClassNoticeDetailPage(row.Index, opt)
	if err != nil {
		return nil, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()

	classNoticeDetail, err := scrape.ClassNoticeDetail(body)
	if err != nil {
		return nil, err
	}

	return classNoticeDetail, nil
}

func (c *Client) fetchClassNoticeSearchPage(opt *model.ClassNoticeSearchOption) (io.ReadCloser, error) {
	body, err := c.fetchGeneralPurposeClassHomePage()
	if err != nil {
		return nil, err
	}
	body.Close()
	_, _ = io.Copy(io.Discard, body)

	body, err = c.fetchGeneralPurposeClassNoticePage()
	if err != nil {
		return nil, err
	}
	body.Close()
	_, _ = io.Copy(io.Discard, body)

	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/classcontact/classContactList/selectClassContactList"
	data := opt.Formdata()
	return c.getPage(reqUrl, *data)
}

func (c *Client) fetchClassNoticeDetailPage(index int, opt *model.ClassNoticeSearchOption) (io.ReadCloser, error) {
	reqUrl := fmt.Sprintf("https://gakujo.shizuoka.ac.jp/portal/classcontact/classContactList/goDetail/%d", index)
	data := opt.Formdata()
	return c.getPage(reqUrl, *data)
}

func (c *Client) fetchGeneralPurposeClassNoticePage() (io.ReadCloser, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
	data := url.Values{}
	data.Set("headTitle", "授業サポート")
	data.Set("menuCode", "A01")
	data.Set("nextPath", "/classcontact/classContactList/initialize")
	return c.getPage(reqUrl, data)
}

func (c *Client) fetchGeneralPurposeClassHomePage() (io.ReadCloser, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
	data := url.Values{}
	data.Set("headTitle", "ホーム")
	data.Set("menuCode", "A00")
	data.Set("nextPath", "/classsupporttop/classSupportTop/initialize")
	return c.getPage(reqUrl, data)
}
