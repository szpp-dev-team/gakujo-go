package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HomeInfo struct {
	TaskRows   []TaskRow   // 未提出課題一覧
	NoticeRows []NoticeRow // お知らせ
}

type TaskRow struct {
	Type     TaskType
	Deadline time.Time
	Name     string
	Index    int
}

type NoticeRow struct {
	Type        NoticeType
	SubType     SubNoticeType
	important   bool
	Date        time.Time
	Title       string
	Affiliation string
	Index       int
}

func (c *Client) Home() (HomeInfo, error) {
	body, err := c.fetchHomeHtml()
	if err != nil {
		return HomeInfo{}, err
	}
	defer body.Close()
	b, _ := io.ReadAll(body)
	taskRows, err := scrapTasks(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return HomeInfo{}, err
	}
	noticeRows, err := scrapNotices(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return HomeInfo{}, err
	}
	return HomeInfo{
		TaskRows:   taskRows,
		NoticeRows: noticeRows,
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

func scrapNotices(r io.Reader) ([]NoticeRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	noticeRows := make([]NoticeRow, 0)
	doc.Find("#tbl_news > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		noticeType, inerr := ToNoticetype(selection.Find("td.arart > span > span > a").Text())
		if inerr != nil {
			err = inerr
			return
		}
		titleLine := selection.Find("td.title > a").Text()
		snt, important, title, inerr := scrapTitleLine(titleLine)
		if inerr != nil {
			err = inerr
			return
		}
		dateText := selection.Find("td.day").Text()
		date, inerr := time.Parse("2006/01/02", dateText)
		if inerr != nil {
			err = inerr
			return
		}
		noticeRow := NoticeRow{
			Type:        noticeType,
			SubType:     snt,
			important:   important,
			Title:       title,
			Date:        date,
			Affiliation: selection.Find("td.syozoku").Text(),
			Index:       i,
		}
		noticeRows = append(noticeRows, noticeRow)
	})
	if err != nil {
		return nil, err
	}
	return noticeRows, nil
}

// return (SubNoticeType, isImportant, title)
func scrapTitleLine(s string) (SubNoticeType, bool, string, error) {
	big := false
	squ := false
	bigText := ""
	squText := ""
	important := false
	title := ""
	for _, c := range s {
		if c == '【' {
			big = true
			continue
		}
		if c == '】' {
			big = false
			continue
		}
		if c == '[' {
			squ = true
			continue
		}
		if c == ']' {
			squ = false
			continue
		}
		if big {
			bigText += string(c)
		} else if squ {
			squText += string(c)
		} else {
			title += string(c)
		}
	}
	if bigText == "重要" {
		important = true
	}
	snt := SNTNone
	if squText != "" {
		var err error
		snt, err = ToSubNoticetype(squText)
		if err != nil {
			return 0, false, "", err
		}
	}
	return snt, important, title, nil
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
