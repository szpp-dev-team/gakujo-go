package gakujo

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/szpp-dev-team/gakujo-go/model"
	"github.com/szpp-dev-team/gakujo-go/scrape"
)

func (c *Client) ClassNoticeRows(opt *model.ClassNoticeSearchOption) ([]model.ClassNoticeRow, error) {
	b, err := c.fetchClassNoticeSearchPage(opt)
	if err != nil {
		return nil, err
	}
	return scrape.ClassNoticeRows(bytes.NewReader(b))
}

func (c *Client) ClassNoticeDetail(row *model.ClassNoticeRow, opt *model.ClassNoticeSearchOption) (*model.ClassNoticeDetail, error) {
	b, err := c.fetchClassNoticeDetailPage(row.Index, opt)
	if err != nil {
		return nil, err
	}
	classNoticeDetail, err := scrape.ClassNoticeDetail(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return classNoticeDetail, nil
}

func (c *Client) fetchClassNoticeSearchPage(opt *model.ClassNoticeSearchOption) ([]byte, error) {
	if _, err := c.fetchGeneralPurposeClassHomePage(); err != nil {
		return nil, err
	}
	if _, err := c.fetchGeneralPurposeClassNoticePage(); err != nil {
		return nil, err
	}
	return c.GetPage("https://gakujo.shizuoka.ac.jp/portal/classcontact/classContactList/selectClassContactList", opt.Formdata())
}

func (c *Client) fetchClassNoticeDetailPage(index int, opt *model.ClassNoticeSearchOption) ([]byte, error) {
	return c.GetPage(
		fmt.Sprintf("https://gakujo.shizuoka.ac.jp/portal/classcontact/classContactList/goDetail/%d", index),
		opt.Formdata(),
	)
}

func (c *Client) fetchGeneralPurposeClassNoticePage() ([]byte, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
	data := url.Values{}
	data.Set("headTitle", "授業サポート")
	data.Set("menuCode", "A01")
	data.Set("nextPath", "/classcontact/classContactList/initialize")
	return c.GetPage(reqUrl, data)
}

func (c *Client) fetchGeneralPurposeClassHomePage() ([]byte, error) {
	reqUrl := "https://gakujo.shizuoka.ac.jp/portal/common/generalPurpose/"
	data := url.Values{}
	data.Set("headTitle", "ホーム")
	data.Set("menuCode", "A00")
	data.Set("nextPath", "/classsupporttop/classSupportTop/initialize")
	return c.GetPage(reqUrl, data)
}
