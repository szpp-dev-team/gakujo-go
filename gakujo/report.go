package gakujo

import (
	"bytes"
	"io"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) ReportList() (model.ReportListInfo, error) {
	body, err := c.fetchReportListHtml()
	if err != nil {
		return model.ReportListInfo{}, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()

	b, _ := io.ReadAll(body)
	if err != nil {
		return model.ReportListInfo{}, err
	}

	reportRows, err := scrape.ReportRows(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return model.ReportListInfo{}, err
	}
	return model.ReportListInfo{
		ReportRows: reportRows,
	}, nil
}

func (c *Client) fetchReportListHtml() (io.ReadCloser, error) {
	datas := make(url.Values)
	datas.Set("headTitle", "授業サポート")
	datas.Set("menuCode", "A02") // TODO: menucode を定数化(まとめてやる)
	datas.Set("nextPath", "/report/student/searchList/initialize")

	return c.getPage(GeneralPurposeUrl, datas)
}
