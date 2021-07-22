package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HomeInfo struct {
	TaskRows []TaskRow // 未提出課題一覧
}

type TaskRow struct {
	Type     TaskType
	Deadline time.Time
	Name     string
	Index    int
}

func (c *Client) Home() (HomeInfo, error) {
	body, err := c.fetchHomeHtml()
	if err != nil {
		return HomeInfo{}, err
	}
	defer body.Close()
	taskRows, err := scrapTasks(body)
	if err != nil {
		return HomeInfo{}, err
	}

	return HomeInfo{
		TaskRows: taskRows,
	}, nil
}

func scrapTasks(r io.Reader) ([]TaskRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	taskRows := make([]TaskRow, 0)
	doc.Find("#tbl_submission > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		taskType, inerr := ToTasktype(selection.Find("td.arart > span > span").Text())
		if inerr != nil {
			err = inerr
			return
		}
		deadlineText := selection.Find("td.daytime").Text()
		var deadline time.Time
		if deadlineText != "" {
			deadline, inerr = time.Parse("2006/01/02 15:04", deadlineText)
			if inerr != nil {
				err = inerr
				return
			}
		}
		taskRow := TaskRow{
			Type:     taskType,
			Deadline: deadline,
			Name:     selection.Find("td:nth-child(3) > a").Text(),
			Index:    i,
		}
		taskRows = append(taskRows, taskRow)
	})
	if err != nil {
		return nil, err
	}
	return taskRows, nil
}

func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/home/home/initialize"

	params := make(url.Values)
	params.Set("EXCLUDE_SET", "")

	resp, err := c.client.PostForm(reqURL, params)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, c.client.Jar.Cookies(resp.Request.URL))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expect %d)", resp.StatusCode, http.StatusOK)
	}

	return resp.Body, nil
}
