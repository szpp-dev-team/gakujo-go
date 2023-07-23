package gakujo

import (
	"bytes"
	"net/url"

	"github.com/szpp-dev-team/gakujo-go/model"
	"github.com/szpp-dev-team/gakujo-go/scrape"
)

func (c *Client) Home() (*model.HomeInfo, error) {
	b, err := c.fetchHomeHtml()
	if err != nil {
		return nil, err
	}
	taskRows, err := scrape.TaskRows(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &model.HomeInfo{
		TaskRows: taskRows,
	}, nil
}

func (c *Client) fetchHomeHtml() ([]byte, error) {
	datas := make(url.Values)
	datas.Set("headTitle", "ホーム")
	datas.Set("menuCode", "Z07") // TODO: menucode を定数化(まとめてやる)
	datas.Set("nextPath", "/home/home/initialize")
	return c.GetPage(GeneralPurposeUrl, datas)
}
